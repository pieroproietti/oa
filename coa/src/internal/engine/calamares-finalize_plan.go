package engine

import (
	"bytes"
	"encoding/json"
	"os"
)

// IsUEFI controlla se il sistema è avviato in modalità UEFI
func IsUEFI() bool {
	_, err := os.Stat("/sys/firmware/efi")
	return !os.IsNotExist(err)
}

// GenerateFinalizePlan crea il JSON per l'ultimo step di Calamares
func GenerateFinalizePlan() error {
	isUEFI := IsUEFI()
	var grubCmd string

	if isUEFI {
		grubCmd = "grub-install --target=x86_64-efi --efi-directory=/boot/efi --bootloader-id=OA > /var/log/grub-debug.log 2>&1 && grub-mkconfig -o /boot/grub/grub.cfg >> /var/log/grub-debug.log 2>&1"
	} else {
		grubCmd = "grub-install /dev/sda && grub-mkconfig -o /boot/grub/grub.cfg >> /var/log/grub-debug.log 2>&1"
	}

	// TODO: Qui un giorno aggiungeremo lo switch per Arch/Debian/Fedora per l'initramfs
	// initramfsCmd := "update-initramfs -u -k all"

	plan := FlightPlan{
		Mode:       "install",
		PathLiveFs: "/tmp/coa/calamares-root",
		Plan: []Action{
			{
				Command:    "oa_shell",
				Info:       "Installazione bootloader (GRUB)",
				RunCommand: grubCmd,
				Chroot:     true,
			},
		},
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	err := enc.Encode(plan)
	if err != nil {
		return err
	}

	os.MkdirAll("/tmp/coa", 0755)
	return os.WriteFile("/tmp/coa/finalize-plan.json", buf.Bytes(), 0644)
}
