package engine

import (
	"fmt"
	"path/filepath"
)

// GenerateLiveStruct ricrea fedelmente la logica del vecchio remaster_livestruct.c
// Si aspetta i percorsi del kernel e dell'initrd già pronti (costruiti prima da coa)
func generateLiveStruct(workDir string, sourceKernel string, sourceInitrd string) ([]Action, error) {
	var actions []Action

	// In penguins-eggs la cartella è SEMPRE "live", per tutte le distro!
	liveDirPath := filepath.Join(workDir, "isodir", "live")

	// 1. Azione: Creazione della struttura (mkdir -p /isodir/live)
	actions = append(actions, Action{
		Command:    "oa_shell",
		Info:       "Creazione della struttura isodir/live",
		RunCommand: fmt.Sprintf("mkdir -p %s", liveDirPath),
	})

	// 2. Azione: Copia del Kernel (vmlinuz)
	actions = append(actions, Action{
		Command:    "oa_shell",
		Info:       "Copia del Kernel standardizzato in /live/vmlinuz",
		RunCommand: fmt.Sprintf("cp %s %s/vmlinuz", sourceKernel, liveDirPath),
	})

	// 3. Azione: Copia dell'Initrd (initrd.img)
	actions = append(actions, Action{
		Command:    "oa_shell",
		Info:       "Copia dell'Initrd standardizzato in /live/initrd.img",
		RunCommand: fmt.Sprintf("cp %s %s/initrd.img", sourceInitrd, liveDirPath),
	})

	return actions, nil
}
