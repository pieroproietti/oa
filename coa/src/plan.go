// Copyright 2026 Piero Proietti <piero.proietti@gmail.com>.
// All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

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
func GeneratePlan(d *Distro, mode string, workPath string) FlightPlan {
	plan := FlightPlan{
		PathLiveFs: workPath,
		Mode:       mode,
		Family:     d.FamilyID, // Comunichiamo la famiglia al braccio C
	}

	// 1. Configurazione Parametri di Boot e Initramfs
	bootParams := "boot=live components quiet splash"
	if d.FamilyID == "archlinux" {
		bootParams = "archisobasedir=arch archisolabel=OA_LIVE quiet splash"
	}

	switch d.FamilyID {
	case "debian":
		plan.InitrdCmd = "mkinitramfs -o {{out}} {{ver}}"
		plan.BootloadersPath = ""
	case "archlinux":
		// Il trucco di Arch: base dir e label per archiso
		// bootParams = "archisobasedir=live archisolabel=OA_LIVE quiet splash"
		plan.InitrdCmd = "mkinitcpio -g {{out}} -k {{ver}}"
		plan.BootloadersPath = BootloaderRoot
	case "fedora", "opensuse":
		plan.InitrdCmd = "dracut --nomadas --force {{out}} {{ver}}"
		plan.BootloadersPath = BootloaderRoot
	default:
		plan.InitrdCmd = "mkinitramfs -o {{out}} {{ver}}"
		plan.BootloadersPath = ""
	}

	// 2. Configurazione Utenti (Globale)
	if mode == "standard" {
		adminGroup := "sudo"
		if d.FamilyID == "archlinux" || d.FamilyID == "fedora" {
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

	// Task specifici per Fedora
	if d.FamilyID == "fedora" {
		plan.Plan = append(plan.Plan, Action{
			Command:    "sys_run",
			RunCommand: "cp",
			Args:       []string{"/tmp/coa/configs/dracut/fedora.conf", "/etc/dracut.conf.d/coa.conf"},
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
			CryptedPassword: "evolution",
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
