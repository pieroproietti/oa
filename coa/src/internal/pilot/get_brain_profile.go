package pilot

import (
	"coa/src/internal/distro"
	"fmt"
	"strings"
)

func GetBrainProfile(d *distro.Distro, mode string, workPath string) *BrainProfile {
	profile := &BrainProfile{Tasks: []Task{}}

	// 1. AREA IDENTITY (Chiamata singola e pulita)
	if idTask, err := createIdentityTask(d, mode); err == nil && idTask != nil {
		profile.Tasks = append(profile.Tasks, *idTask)
	}

	// 2. AREA INITRD
	if iniTask, err := createInitrdTask(d); err == nil && iniTask != nil {
		profile.Tasks = append(profile.Tasks, *iniTask)
	}

	// 3. AREA BOOT (Solo per i parametri, niente comandi shell qui)
	var bootCfg BootConfig
	if err := readAreaConfig(d.FamilyID, "boot", &bootCfg); err == nil {
		profile.Tasks = append(profile.Tasks, Task{
			Name:        "boot",
			Description: "Kernel boot parameters",
			// Lasciamo Commands vuoto per evitare l'errore 'run_command mancante'
		})
	}

	// 4, In GetBrainProfile
	if layTask, err := createLayoutTask(d, workPath); err == nil {
		profile.Tasks = append(profile.Tasks, *layTask)
	}

	return profile
}

// createIdentityTask
func createIdentityTask(d *distro.Distro, mode string) (*Task, error) {
	var cfg IdentityConfig
	if err := readAreaConfig(d.FamilyID, "identity", &cfg); err != nil {
		return nil, err
	}

	task := &Task{
		Name:        "identity",
		Description: "User & Group Identity artisan",
		Chroot:      true,
		Commands:    []string{},
	}

	// 1. CASO CLONE: Non facciamo nulla, preserviamo gli utenti dell'host
	if mode == "clone" {
		task.Description += " (Preserving host identities - Clone Mode)"
		return task, nil
	}

	// 2. CASO STANDARD o CRYPTED: Modifichiamo utenti e gruppi
	if mode == "standard" || mode == "crypted" {
		login := "live" // Nome dell'utente live di default

		// A. PULIZIA (Purge)
		// Rimuoviamo gli utenti umani esistenti (UID 1000-59999)
		purgeCmd := "awk -F: '$3 < 1000 || $3 > 59999' /etc/passwd > /tmp/p && mv /tmp/p /etc/passwd && " +
			"awk -F: '$3 < 1000 || $3 > 59999' /etc/group > /tmp/g && mv /tmp/g /etc/group"
		task.Commands = append(task.Commands, purgeCmd)

		// B. INIEZIONE NUOVA IDENTITÀ (Passwd & Shadow)
		// Usiamo i tuoi comandi echo o useradd
		task.Commands = append(task.Commands, fmt.Sprintf("echo '%s:x:1000:1000:%s,,,:/home/%s:/bin/bash' >> /etc/passwd", login, login, login))
		// Inseriamo l'hash della password (recuperato dal brain o default)
		task.Commands = append(task.Commands, fmt.Sprintf("echo '%s:$6$hash_di_esempio...:19000:0:99999:7:::' >> /etc/shadow", login))

		// C. MAPPATURA GRUPPI (La logica "Artisan" universale che abbiamo discusso)
		groupMapping := getGroupMappingCommand(login, cfg.AdminGroup, cfg.UserGroups)
		task.Commands = append(task.Commands, groupMapping)

		// D. ABILITAZIONE SUDOERS
		sudoCmd := fmt.Sprintf("sed -i 's/^# %%%s ALL=(ALL:ALL) ALL/%%%s ALL=(ALL:ALL) ALL/' /etc/sudoers",
			cfg.AdminGroup, cfg.AdminGroup)
		task.Commands = append(task.Commands, sudoCmd)

		// E. SETUP HOME (Skel & Permissions)
		homeCmd := fmt.Sprintf("mkdir -p /home/%s && cp -a /etc/skel/. /home/%s/ 2>/dev/null || true && chown -R 1000:1000 /home/%s",
			login, login, login)
		task.Commands = append(task.Commands, homeCmd)
	}

	return task, nil
}

// createInitrdTask è una funzione interna (minuscola) che costruisce il Task specifico.
func createInitrdTask(d *distro.Distro) (*Task, error) {
	var cfg InitrdConfig

	// Carica il file (es. debian/initrd.yaml)
	if err := readAreaConfig(d.FamilyID, "initrd", &cfg); err != nil {
		return nil, err
	}

	task := &Task{
		Name:        "initrd",
		Description: fmt.Sprintf("Regenerating initramfs for %s", d.FamilyID),
		Files:       make(map[string]string),
		Commands:    []string{},
		Chroot:      true,
	}

	// Ora carichiamo i dati direttamente (niente più .Initrd.Live o wrapper)
	for path, content := range cfg.Files {
		task.Files[path] = content
	}

	if cfg.Command != "" {
		task.Commands = append(task.Commands, cfg.Command)
	}

	return task, nil
}

func createLayoutTask(d *distro.Distro, workPath string) (*Task, error) {
	var cfg LayoutConfig
	if err := readAreaConfig(d.FamilyID, "layout", &cfg); err != nil {
		return nil, err
	}

	task := &Task{
		Name:        "layout",
		Description: "Creating ISO layout symlinks",
		Commands:    []string{},
		Chroot:      false, // Lavoriamo sui file della ISO
	}

	for dst, src := range cfg.Links {
		// Ricostruiamo il comando che avevi nell'engine
		linkCmd := fmt.Sprintf("mkdir -p $(dirname %s/iso/%s) && ln -sf %s %s/iso/%s",
			workPath, dst, src, workPath, dst)
		task.Commands = append(task.Commands, linkCmd)
	}

	return task, nil
}

// getGroupMappingCommand chiamata da createIdentityTask sopra
func getGroupMappingCommand(login string, adminGroup string, userGroups []string) string {
	// Uniamo il gruppo admin ai gruppi utente
	allGroups := append([]string{adminGroup}, userGroups...)
	var commands []string

	for _, g := range allGroups {
		if g == "" {
			continue
		}
		// Questo sed cerca la riga del gruppo e appende ",utente" alla fine
		// Funziona anche se la lista utenti è vuota (es. wheel:x:10:)
		cmd := fmt.Sprintf("sed -i -E '/^%s:/ s/([^:]*)$/\\1,%s/' /etc/group", g, login)
		commands = append(commands, cmd)
	}

	// Pulizia: trasforma ":,utente" in ":utente" per i gruppi che erano inizialmente vuoti
	commands = append(commands, "sed -i 's/:,/:/g' /etc/group")

	return strings.Join(commands, " && ")
}
