package engine

import (
	"fmt"
	"path/filepath"
)

// generateSquashfs crea l'azione per comprimere la liveroot nel filesystem.squashfs
func generateSquashfs(workPath string, compType string, excludeFile string) ([]Action, error) {
	var actions []Action

	// 1. Definiamo i percorsi esatti usando workPath
	liveroot := filepath.Join(workPath, "liveroot")
	squashDest := filepath.Join(workPath, "isodir", "live", "filesystem.squashfs")

	// 2. Fallback per la compressione
	if compType == "" {
		compType = "xz"
	}

	// 3. Creazione del comando per mksquashfs
	// -ef usa il file delle esclusioni che abbiamo generato in remaster_plan.go
	cmd := fmt.Sprintf(
		"mksquashfs %s %s -comp %s -b 1M -noappend -wildcards -ef %s",
		liveroot, squashDest, compType, excludeFile,
	)

	// 4. Aggiungiamo l'azione da passare al motore C
	actions = append(actions, Action{
		Command:    "oa_shell",
		Info:       fmt.Sprintf("Compressione SquashFS (%s) - Questa operazione richiederà tempo...", compType),
		RunCommand: cmd,
	})

	return actions, nil
}
