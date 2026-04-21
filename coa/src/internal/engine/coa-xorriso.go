package engine

import (
	"fmt"
	"path/filepath"
	// "coa/src/utils" // Adatta l'import se BootloaderRoot è altrove
)

// GenerateIso crea l'azione per impacchettare isodir nella ISO finale
func GenerateIso(workDir string, isoName string, volId string) ([]Action, error) {
	var actions []Action

	// 1. Definiamo i percorsi chiave
	isodir := filepath.Join(workDir, "isodir")
	isoDest := filepath.Join(workDir, isoName) // Es: /home/eggs/custom.iso

	// Se il Volume ID è vuoto, diamo un default sensato (max 32 caratteri, no spazi strani)
	if volId == "" {
		volId = "PENGUINS_EGGS_LIVE"
	}

	// 2. Recuperiamo il file MBR ibrido dai bootloader scaricati.
	// Questo file è vitale per far avviare la chiavetta USB sui vecchi BIOS.
	mbrPath := filepath.Join(BootloaderRoot, "isolinux/isohdpfx.bin")

	// 3. Costruzione del "Mega Comando" xorriso
	// -as mkisofs: usa la sintassi standard e compatibile
	// -r -J -joliet-long: abilita i nomi di file lunghi (Windows/Linux compatibilità)
	// -b: punta al bootloader BIOS
	// -e: punta all'immagine UEFI FAT che abbiamo creato al passo precedente
	// -isohybrid-*: rende la ISO scrivibile direttamente su USB con 'dd'
	cmd := fmt.Sprintf(
		"xorriso -as mkisofs -r -J -joliet-long -l -cache-inodes "+
			"-V '%s' "+
			"-isohybrid-mbr %s -partition_offset 16 "+
			"-A '%s' "+
			"-b isolinux/isolinux.bin -c isolinux/boot.cat "+
			"-no-emul-boot -boot-load-size 4 -boot-info-table "+
			"-eltorito-alt-boot "+
			"-e boot/grub/efi.img -no-emul-boot -isohybrid-gpt-basdat "+
			"-o %s %s",
		volId, mbrPath, volId, isoDest, isodir,
	)

	// 4. Deleghiamo l'esecuzione della ISO cruda al motore C
	actions = append(actions, Action{
		Command:    "oa_shell", // Eseguito sull'host, non in chroot
		Info:       fmt.Sprintf("Creazione della ISO finale: %s", isoName),
		RunCommand: cmd,
	})

	return actions, nil
}
