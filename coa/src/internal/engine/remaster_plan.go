package engine

import (
	"coa/src/internal/distro"
	"coa/src/internal/pilot"
	"fmt"
	"path/filepath"
)

func generatePlan(d *distro.Distro, mode string, workPath string) FlightPlan {
	// 0. Inizializzazione della base del Piano di Volo
	plan := FlightPlan{
		PathLiveFs:      workPath,
		Mode:            mode,
		Family:          d.FamilyID,
		BootloadersPath: BootloaderRoot,
	}

	// =================================================================
	// FASE 1: THE BRAIN (Operazioni Interne - Chroot)
	// Modifichiamo il cuore del sistema (Utenti, Initramfs, Configurazioni)
	// =================================================================
	profile := pilot.GetBrainProfile(d, mode, workPath)

	for _, t := range profile.Tasks {
		for _, cmd := range t.Commands {
			plan.Plan = append(plan.Plan, Action{
				Command:    "oa_shell",
				Info:       t.Description,
				RunCommand: cmd,
				Chroot:     t.Chroot,
			})
		}
	}

	// =================================================================
	// FASE 2: LIVE STRUCT (Operazioni Esterne - Host)
	// Copiamo il kernel e l'initrd (appena rigenerato!) nella ISO
	// =================================================================

	// NOTA: I percorsi di vmlinuz e initrd variano da distro a distro.
	// Di solito in Debian/Ubuntu puntano ai symlink in root. Essendo eseguiti sull'host,
	// dobbiamo puntare all'interno della liveroot.
	kernelSrc := filepath.Join(workPath, "liveroot", "vmlinuz")
	initrdSrc := filepath.Join(workPath, "liveroot", "initrd.img")

	lsActions, err := generateLiveStruct(workPath, kernelSrc, initrdSrc)
	if err != nil {
		fmt.Printf("\033[1;33m[WARNING]\033[0m Impossibile generare LiveStruct: %v\n", err)
	} else {
		plan.Plan = append(plan.Plan, lsActions...)
	}

	// =================================================================
	// FASE 3: BOOTLOADERS (Operazioni Esterne - Host)
	// Iniettiamo GRUB, ISOLINUX e generiamo l'immagine EFI FAT32
	// =================================================================
	blActions, err := LiveBootloader(workPath)
	if err != nil {
		fmt.Printf("\033[1;33m[WARNING]\033[0m Impossibile configurare Bootloader: %v\n", err)
	} else {
		plan.Plan = append(plan.Plan, blActions...)
	}

	// =================================================================
	// FASE 4: SQUASHFS (Il Congelamento)
	// Chiudiamo il filesystem modificato in un blocco compresso
	// =================================================================
	excludeFilePath := generateExcludeList(mode)
	sqActions, err := generateSquashfs(workPath, "xz", excludeFilePath)
	if err != nil {
		fmt.Printf("\033[1;31m[ERRORE]\033[0m Generazione SquashFS fallita: %v\n", err)
	} else {
		plan.Plan = append(plan.Plan, sqActions...)
	}

	// =================================================================
	// FASE 5: XORRISO (Il Confezionamento)
	// Creiamo la ISO Ibrida avviabile da USB e CD
	// =================================================================
	isoName := getIsoName(d)
	isoActions, err := GenerateIso(workPath, isoName, "OA_LIVE")
	if err != nil {
		fmt.Printf("\033[1;31m[ERRORE]\033[0m Generazione ISO fallita: %v\n", err)
	} else {
		plan.Plan = append(plan.Plan, isoActions...)
	}

	// =================================================================
	// FASE 6: CLEANUP (Il Demolitore)
	// =================================================================
	plan.Plan = append(plan.Plan, Action{
		Command: "oa_remaster_cleanup",
		Info:    "Smontaggio filesystem virtuali e pulizia finale",
	})

	return plan
}
