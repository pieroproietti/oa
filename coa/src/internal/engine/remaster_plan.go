package engine

import (
	"coa/src/internal/distro"
	"coa/src/internal/pilot"
	"fmt"
)

func generatePlan(d *distro.Distro, mode string, workPath string) FlightPlan {
	// 1. Il Cervello fornisce il nuovo Profilo Modulare
	profile := pilot.GetBrainProfile(d, mode, workPath)
	// Debug rapido per contare i pezzi
	for _, t := range profile.Tasks {
		fmt.Printf("DEBUG: Task '%s' caricato con %d comandi\n", t.Name, len(t.Commands))
	}

	// Fallback temporanei per le parti non ancora migrate (Identity/Boot)
	bootParams := "boot=live components"

	plan := FlightPlan{
		PathLiveFs:      workPath,
		Mode:            mode,
		Family:          d.FamilyID,
		BootloadersPath: BootloaderRoot,
	}

	for _, t := range profile.Tasks {
		for _, cmd := range t.Commands {
			plan.Plan = append(plan.Plan, Action{
				Command:    "oa_sys_shell",
				Info:       t.Description,
				RunCommand: cmd,
				Chroot:     t.Chroot,
			})
		}
	}

	// --- 4. ISO & BOOTLOADERS (Ancora legati alle vecchie azioni C) ---
	excludeFilePath := generateExcludeList(mode)

	plan.Plan = append(plan.Plan,
		Action{Command: "oa_remaster_livestruct"},
		Action{Command: "oa_remaster_isolinux", BootParams: bootParams},
		Action{Command: "oa_remaster_uefi", BootParams: bootParams},
	)

	// fix da rimuovere sovrascrive uefi menu
	err := pilot.GenerateBootConfig(d.FamilyID, profile)
	if err != nil {
		// Se fallisce, stampiamo l'errore ma proviamo a procedere
		// (o blocchiamo qui se preferisci un approccio più rigido)
		fmt.Printf("[ERRORE] Il Pilota non ha generato il boot config: %v\n", err)
	}
	plan.Plan = append(plan.Plan, Action{
		Command:    "oa_sys_shell",
		RunCommand: "cp /tmp/coa/grub.cfg.final /home/eggs/iso/boot/grub/grub.cfg",
		Chroot:     false,
		Info:       "Injecting generated GRUB configuration",
	})

	// --- 5. CHIUSURA (Fasi pesanti in C) ---
	plan.Plan = append(plan.Plan, Action{
		Command:     "oa_remaster_squash",
		ExcludeList: excludeFilePath,
	})

	isoName := getIsoName(d)
	plan.Plan = append(plan.Plan, Action{
		Command:   "oa_remaster_iso",
		VolID:     "OA_LIVE",
		OutputISO: isoName,
	})

	plan.Plan = append(plan.Plan, Action{Command: "oa_remaster_cleanup"})

	return plan
}
