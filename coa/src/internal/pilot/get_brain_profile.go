package pilot

import (
	"coa/src/internal/distro"
	"fmt"
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

func createIdentityTask(d *distro.Distro, mode string) (*Task, error) {
	var cfg IdentityConfig
	if err := readAreaConfig(d.FamilyID, "identity", &cfg); err != nil {
		return nil, err
	}

	task := &Task{
		Name:        "identity",
		Description: "Artisan identity injection (passwd/shadow)",
		Files:       make(map[string]string),
		Commands:    []string{},
		Chroot:      true,
	}

	// 1. PULIZIA (Purge): Se non siamo in clone, puliamo gli utenti dell'host
	if mode != "clone" && mode != "crypted" {
		purgeCmd := "awk -F: '$3 < 1000 || $3 > 59999' /etc/passwd > /tmp/p && mv /tmp/p /etc/passwd && " +
			"awk -F: '$3 < 1000 || $3 > 59999' /etc/group > /tmp/g && mv /tmp/g /etc/group"
		task.Commands = append(task.Commands, purgeCmd)
	}

	// 2. INIEZIONE: Creiamo l'utente live (esempio semplificato per un utente)
	// In futuro potrai iterare su una lista di utenti dal profilo
	login := "live"
	uid, gid := 1000, 1000
	gecos := "live,,,"
	home := "/home/live"
	shell := "/bin/bash"
	passwordHash := "$6$wM.wY0QtatvbQMHZ$QtIKXSpIsp2Sk57.Ny.JHk7hWDu.lxPtUYaTOiBnP4WBG5KS6JpUlpXj2kcSaaMje7fr01uiGmxZhE8kfZRqv."

	passwdLine := fmt.Sprintf("%s:x:%d:%d:%s:%s:%s", login, uid, gid, gecos, home, shell)
	shadowLine := fmt.Sprintf("%s:%s:19000:0:99999:7:::", login, passwordHash)

	injectCmd := fmt.Sprintf("echo '%s' >> /etc/passwd && echo '%s' >> /etc/shadow", passwdLine, shadowLine)
	task.Commands = append(task.Commands, injectCmd)

	// 3. HOME & SKEL: La tua logica C tradotta in shell
	homeCmd := fmt.Sprintf("mkdir -p %s && cp -a /etc/skel/. %s/ 2>/dev/null || true && chown -R %d:%d %s",
		home, home, uid, gid, home)
	task.Commands = append(task.Commands, homeCmd)

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
