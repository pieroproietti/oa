package engine

import (
	"fmt"
	"os"
	"path/filepath"
)

// LiveBootloader utilizza la variabile globale BootloaderRoot
func LiveBootloader(workDir string) ([]Action, error) {
	var actions []Action
	isodir := filepath.Join(workDir, "isodir")

	// 1. Definiamo i percorsi di destinazione interni alla ISO
	grubEfiDest := filepath.Join(isodir, "boot/grub/x86_64-efi")
	isolinuxDest := filepath.Join(isodir, "isolinux")
	efiBootDest := filepath.Join(isodir, "EFI/BOOT")

	// 2. Creazione Directory tramite Go
	dirs := []string{grubEfiDest, isolinuxDest, efiBootDest}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return nil, err
		}
	}

	// 3. Il Trampolino GRUB nativo in Go
	trampolino := "search --file --set=root /live/vmlinuz\n" +
		"set prefix=($root)/boot/grub\n" +
		"configfile $prefix/grub.cfg\n"

	trampolinoPath := filepath.Join(efiBootDest, "grub.cfg")
	if err := os.WriteFile(trampolinoPath, []byte(trampolino), 0644); err != nil {
		return nil, err
	}

	// 4. AZIONI PER IL MOTORE OA: Copia usando BootloaderRoot

	// A. Copia dei moduli GRUB (Usiamo cp -a con il suffisso /.)
	srcGrub := filepath.Join(BootloaderRoot, "grub/x86_64-efi")
	actions = append(actions, Action{
		Command:    "oa_shell",
		Info:       "Copia moduli GRUB da penguins-bootloaders",
		RunCommand: fmt.Sprintf("cp -a %s/. %s/", srcGrub, grubEfiDest),
	})

	// B. Copia dei file ISOLINUX (MAIUSCOLO e cp -a con il suffisso /.)
	srcIsolinux := filepath.Join(BootloaderRoot, "ISOLINUX")
	actions = append(actions, Action{
		Command:    "oa_shell",
		Info:       "Copia file ISOLINUX da penguins-bootloaders",
		RunCommand: fmt.Sprintf("cp -a %s/. %s/", srcIsolinux, isolinuxDest),
	})

	// C. Costruzione dell'immagine EFI (efi.img)
	efiImgPath := filepath.Join(isodir, "boot/grub/efi.img")
	srcEfiMonolithic := filepath.Join(BootloaderRoot, "grub/monolithic/grubx64.efi")

	cmdEfiImg := fmt.Sprintf(
		"dd if=/dev/zero of=%s bs=1k count=2048 && "+
			"mkfs.vfat %s && "+
			"mmd -i %s ::/EFI && mmd -i %s ::/EFI/BOOT && "+
			"mcopy -i %s %s ::/EFI/BOOT/BOOTX64.EFI && "+
			"mcopy -i %s %s ::/EFI/BOOT/grub.cfg",
		efiImgPath, efiImgPath, efiImgPath, efiImgPath, efiImgPath, srcEfiMonolithic, efiImgPath, trampolinoPath,
	)

	actions = append(actions, Action{
		Command:    "oa_shell",
		Info:       "Creazione immagine EFI (efi.img)",
		RunCommand: cmdEfiImg,
	})

	return actions, nil
}
