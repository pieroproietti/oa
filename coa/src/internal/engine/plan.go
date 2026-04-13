// Copyright 2026 Piero Proietti <piero.proietti@gmail.com>.
// All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package engine

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"coa/src/internal/distro"
)

// BootloaderRoot definisce dove vengono estratti i bootloader.
// L'abbiamo spostata qui da utils.go per centralizzare le costanti del motore.
const BootloaderRoot = "/tmp/coa/bootloaders"

// Action rappresenta un singolo blocco "command" nell'array "plan"
type Action struct {
	Command         string   `json:"command"`
	VolID           string   `json:"volid,omitempty"`
	OutputISO       string   `json:"output_iso,omitempty"`
	CryptedPassword string   `json:"crypted_password,omitempty"`
	RunCommand      string   `json:"run_command,omitempty"`
	ExcludeList     string   `json:"exclude_list,omitempty"`
	BootParams      string   `json:"boot_params,omitempty"` // Parametri dinamici per il bootloader
	Args            []string `json:"args,omitempty"`
}

// UserConfig definisce la struttura per la creazione nativa dell'utente live
type UserConfig struct {
	Login    string   `json:"login"`
	Password string   `json:"password"`
	Gecos    string   `json:"gecos"`
	Home     string   `json:"home"`
	Shell    string   `json:"shell"`
	Groups   []string `json:"groups"`
}

// FlightPlan è l'oggetto JSON principale inviato al motore oa
type FlightPlan struct {
	PathLiveFs      string       `json:"pathLiveFs"`
	Mode            string       `json:"mode"`
	Family          string       `json:"family"` // Passiamo la famiglia per logiche specifiche in C
	InitrdCmd       string       `json:"initrd_cmd"`
	BootloadersPath string       `json:"bootloaders_path"`
	Users           []UserConfig `json:"users"`
	Plan            []Action     `json:"plan"`
}

// generateExcludeList crea il file .list dinamico per mksquashfs
func generateExcludeList(mode string) string {
	outPath := "/tmp/coa/excludes.list"
	var excludes []string

	// 1. Esclusioni Base (Pulizia della ISO)
	excludes = append(excludes,
		"boot/efi/EFI",
		"boot/loader/entries/",
		"etc/fstab",
		"var/lib/docker/",
	)

	// 2. Esclusioni specifiche per modalità
	if mode != "clone" && mode != "crypted" {
		excludes = append(excludes, "root/*")
	}

	// 3. Esclusioni Utente
	userList := "/etc/coa/exclusion.list"
	if _, err := os.Stat(userList); os.IsNotExist(err) {
		userList = "conf/exclusion.list"
	}

	if data, err := os.ReadFile(userList); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" && !strings.HasPrefix(line, "#") {
				excludes = append(excludes, line)
			}
		}
	}

	os.MkdirAll("/tmp/coa", 0755)
	os.WriteFile(outPath, []byte(strings.Join(excludes, "\n")+"\n"), 0644)

	return outPath
}

// GeneratePlan costruisce il piano di volo dinamico
func GeneratePlan(d *distro.Distro, mode string, workPath string) FlightPlan {
	plan := FlightPlan{
		PathLiveFs:      workPath,
		Mode:            mode,
		Family:          d.FamilyID,
		BootloadersPath: BootloaderRoot,
	}

	// 1. Configurazione Parametri di Boot
	bootParams := "boot=live components"
	switch d.FamilyID {
	case "archlinux":
		bootParams = "archisobasedir=arch archisolabel=OA_LIVE"
	case "fedora", "rhel", "centos", "rocky", "almalinux", "opensuse":
		bootParams = "root=live:CDLABEL=OA_LIVE rd.live.image rd.live.dir=live rd.live.squashimg=filesystem.squashfs selinux=0"
	}

	// SWITCH PER L'INITRAMFS
	switch d.FamilyID {
	case "debian":
		plan.InitrdCmd = "mkinitramfs -o {{out}} {{ver}}"
	case "archlinux":
		// Interno a chroot
		plan.InitrdCmd = fmt.Sprintf("mkinitcpio -c %s/liveroot/etc/mkinitcpio.conf -g {{out}} -k {{ver}}", workPath)
	case "fedora", "rhel", "centos", "rocky", "almalinux":
		// Interno a chroot
		plan.InitrdCmd = "dracut --no-hostonly --nomdadmconf --nolvmconf --xz --add dmsquash-live --add rootfs-block --add bash --force {{out}} {{ver}}"
	default:
		plan.InitrdCmd = "mkinitramfs -o {{out}} {{ver}}"
	}

	// 2. Configurazione Utenti (Globale)
	if mode == "standard" {
		adminGroup := "sudo"
		if d.FamilyID == "archlinux" || d.FamilyID == "fedora" || d.FamilyID == "rhel" || d.FamilyID == "centos" || d.FamilyID == "rocky" || d.FamilyID == "almalinux" {
			adminGroup = "wheel"
		}

		plan.Users = []UserConfig{
			{
				Login:    "live",
				Password: "$6$wM.wY0QtatvbQMHZ$QtIKXSpIsp2Sk57.Ny.JHk7hWDu.lxPtUYaTOiBnP4WBG5KS6JpUlpXj2kcSaaMje7fr01uiGmxZhE8kfZRqv.",
				Gecos:    "live,,,",
				Home:     "/home/live",
				Shell:    "/bin/bash",
				Groups:   []string{"cdrom", "audio", "video", "plugdev", "netdev", "autologin", adminGroup},
			},
		}
	} else {
		plan.Users = []UserConfig{}
	}

	// 3. Assemblaggio dinamico del piano
	plan.Plan = []Action{
		{Command: "lay_users"},
	}

	// Task specifici per Fedora/RHEL
	if d.FamilyID == "fedora" || d.FamilyID == "rhel" || d.FamilyID == "centos" || d.FamilyID == "rocky" || d.FamilyID == "almalinux" {
		
		targetConfDir := fmt.Sprintf("%s/liveroot/etc/dracut.conf.d", workPath)
		targetConfPath := fmt.Sprintf("%s/coa.conf", targetConfDir) 

		// Prepariamo il testo del file formattato per il comando echo -e
		dracutConfig := `hostonly="no"\nadd_dracutmodules+=" dmsquash-live rootfs-block bash "\ncompress="xz"`
		
		// Costruiamo il comando shell che scrive direttamente il file
		writeCmd := fmt.Sprintf(`echo -e '%s' > %s`, dracutConfig, targetConfPath)

		// 1. Assicuriamoci che la cartella esista
		plan.Plan = append(plan.Plan, Action{
			Command:    "sys_run",
			RunCommand: "mkdir",
			Args:       []string{"-p", targetConfDir},
		})

		// 2. Scriviamo il file a destinazione usando sh -c (nessun file temporaneo perso nei mount!)
		plan.Plan = append(plan.Plan, Action{
			Command:    "sys_run",
			RunCommand: "sh",
			Args:       []string{"-c", writeCmd},
		})
	}

	excludeFilePath := generateExcludeList(mode)

	// Aggiungiamo le azioni principali passando i bootParams
	plan.Plan = append(plan.Plan,
		Action{Command: "lay_initrd"},
		Action{Command: "lay_livestruct"},
		Action{Command: "lay_isolinux", BootParams: bootParams},
		Action{Command: "lay_uefi", BootParams: bootParams},
		Action{
			Command:     "lay_squash",
			ExcludeList: excludeFilePath,
		},
	)

	if mode == "crypted" {
		plan.Plan = append(plan.Plan, Action{
			Command:         "lay_crypted",
			CryptedPassword: "evolution", // Qui potremmo in futuro parametrizzarlo da CLI
		})
	}

	// Definizione nome ISO
	hostname, _ := os.Hostname()
	timestamp := time.Now().Format("2006-01-02_1504")
	arch := runtime.GOARCH

	var nameParts []string
	nameParts = append(nameParts, d.DistroID)
	if d.CodenameID != "" {
		nameParts = append(nameParts, d.CodenameID)
	} else if d.ReleaseID != "" {
		nameParts = append(nameParts, d.ReleaseID)
	}
	if hostname != "" {
		nameParts = append(nameParts, hostname)
	}

	distroTag := strings.Join(nameParts, "-")
	isoName := fmt.Sprintf("egg-of_%s_%s_%s.iso", distroTag, arch, timestamp)

	plan.Plan = append(plan.Plan, Action{
		Command:   "lay_iso",
		VolID:     "OA_LIVE",
		OutputISO: isoName,
	})

	plan.Plan = append(plan.Plan, Action{Command: "lay_cleanup"})

	return plan
}
