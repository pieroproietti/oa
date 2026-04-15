// Copyright 2026 Piero Proietti <piero.proietti@gmail.com>.
// All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package engine

import (
	"encoding/json"
	"fmt"
	"os"

	"coa/src/internal/distro"
)

func generateInstallPlan(ans *KrillAnswers, disk string, d *distro.Distro) {
	fmt.Println("\n\033[1;34m[krill]\033[0m Compiling JSON flight plan for the C ..")

	squashPath := getSquashfsPath()
	if squashPath == "" {
		fmt.Printf("\033[1;31m[ERROR]\033[0m Could not locate pristine filesystem.squashfs!\n")
		return
	}

	hashedPass := generateHashedPassword(ans.Password)

	// 1. Inizializzazione del Piano di Volo
	plan := FlightPlan{
		PathLiveFs: "/mnt/krill-target",
		Mode:       "install",
		Plan:       []Action{},
	}

	// 2. PREPARAZIONE DISCO E FILESYSTEM
	plan.Plan = append(plan.Plan, Action{Command: "oa_install_partition", RunCommand: disk})

	// FIX PER FEDORA 43: Formattazione EXT4 con flag di compatibilità per GRUB
	if d.FamilyID == "fedora" || d.FamilyID == "rhel" {
		plan.Plan = append(plan.Plan, Action{
			Command:    "oa_sys_shell",
			RunCommand: fmt.Sprintf("mkfs.fat -F32 %s2 && mkfs.ext4 -O ^metadata_csum_seed,^orphan_file %s3", disk, disk),
		})
	} else {
		plan.Plan = append(plan.Plan, Action{Command: "oa_install_format", RunCommand: disk})
	}

	plan.Plan = append(plan.Plan,
		Action{Command: "oa_install_unpack", RunCommand: disk, Args: []string{squashPath}},
		Action{Command: "oa_install_prepare", RunCommand: disk},
	)

	// 3. LOGICA DISTRO-SPECIFICA (The Mind)
	var shellInitrdCmd string
	var shellGrubCmd string

	switch d.FamilyID {
	case "archlinux":
		shellInitrdCmd = "mkinitcpio -P"
		if isUEFI() {
			shellGrubCmd = "grub-install --target=x86_64-efi --efi-directory=/boot/efi --bootloader-id=coa --recheck && grub-mkconfig -o /boot/grub/grub.cfg"
		} else {
			shellGrubCmd = "grub-install --target=i386-pc --recheck " + disk + " && grub-mkconfig -o /boot/grub/grub.cfg"
		}
	case "fedora", "rhel", "centos", "rocky", "almalinux":
		shellInitrdCmd = "dracut --force --regenerate-all"
		// FIX PER FEDORA: Disabilitazione BLS e installazione fisica GRUB
		if isUEFI() {
			shellGrubCmd = "echo 'GRUB_ENABLE_BLSCFG=false' >> /etc/default/grub && " +
				"grub2-install --target=x86_64-efi --efi-directory=/boot/efi --bootloader-id=fedora --recheck && " +
				"grub2-mkconfig -o /boot/grub2/grub.cfg && " +
				"cp /boot/grub2/grub.cfg /boot/efi/EFI/fedora/grub.cfg"
		} else {
			shellGrubCmd = "echo 'GRUB_ENABLE_BLSCFG=false' >> /etc/default/grub && " +
				"grub2-install --target=i386-pc --recheck " + disk + " && " +
				"grub2-mkconfig -o /boot/grub2/grub.cfg"
		}
	default: // Debian/Ubuntu
		shellInitrdCmd = "update-initramfs -u -k all"
		if isUEFI() {
			shellGrubCmd = "grub-install --target=x86_64-efi --efi-directory=/boot/efi --bootloader-id=coa --recheck && update-grub"
		} else {
			shellGrubCmd = "grub-install --target=i386-pc --recheck " + disk + " && update-grub"
		}
	}

	// 4. ESECUZIONE AZIONI UNIVERSALI E CHROOT
	plan.Plan = append(plan.Plan,
		Action{Command: "oa_install_fstab", RunCommand: disk},
		Action{
			Command:    "oa_sys_shell",
			RunCommand: fmt.Sprintf("echo %s > /etc/hostname && systemd-machine-id-setup", ans.Hostname),
			Chroot:     true,
		},
		Action{
			Command:    "oa_sys_shell",
			RunCommand: shellInitrdCmd,
			Chroot:     true,
		},
		Action{
			Command:    "oa_sys_shell",
			RunCommand: shellGrubCmd,
			Chroot:     true,
		},
		Action{Command: "oa_install_users"},
	)

	// 5. PULIZIA SELETTIVA E SYNC (FIX PER HANG SU FEDORA)
	if d.FamilyID == "debian" {
		plan.Plan = append(plan.Plan, Action{
			Command:    "oa_sys_shell",
			RunCommand: "rm -rf /var/log/installer /var/lib/live/config /etc/sudoers.d/live-user 2>/dev/null",
			Chroot:     true,
		})
	}

	plan.Plan = append(plan.Plan,
		Action{Command: "oa_sys_shell", RunCommand: "sync"},
		Action{Command: "oa_install_cleanup"},
	)

	// 6. Configurazione Utente Primario
	plan.Users = []UserConfig{
		{
			Login:    ans.Username,
			Password: hashedPass,
			Gecos:    ans.Username + ",,,",
			Home:     "/home/" + ans.Username,
			Shell:    "/bin/bash",
			Groups:   []string{"sudo", "wheel", "video", "audio", "plugdev", "netdev"},
		},
	}

	// 7. Salvataggio ed Esecuzione
	jsonData, _ := json.MarshalIndent(plan, "", "  ")
	outPath := "/tmp/oa-plan.json"
	os.WriteFile(outPath, jsonData, 0644)

	fmt.Printf("\033[1;32m[SUCCESS]\033[0m Flight plan ready at %s\n", outPath)
	ExecutePlan(plan)
}
