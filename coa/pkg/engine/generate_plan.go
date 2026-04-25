package engine

import (
	"coa/pkg/pilot" // Importiamo i tipi definiti nel pilota
	"encoding/json"
	"os"
	"path/filepath"
)

func GeneratePlan(yamlSteps []pilot.YamlStep, familyID string, isRemaster bool, workPath string) (string, error) {
	var plan OAPlan

	// Definiamo l'utente classico "live/evolution"
	// In futuro questo verrà da pilot.LoadConfig()
	defaultUser := pilot.User{
		Login:    "live",
		Password: "$6$oa-tools$uTKAYeAVn.Y.Dy2To6HXsHt1Gt4HpMghmOV93a46jFY7hkAQ3tk7eRTKjcvSYDf5sOf3qnKzyyPYXurKp9ST3.", // evolution
		Home:     "/home/live",
		Shell:    "/bin/bash",
		Groups:   []string{"sudo", "audio", "video", "cdrom", "plugdev", "netdev"},
		UID:      1000,
		GID:      1000,
	}

	for _, step := range yamlSteps {
		switch step.Command {

		case "oa_mount_logic":
			// Esplode la vecchia logica del C in tanti task JSON
			plan.Plan = append(plan.Plan, expandMountLogic(workPath)...)

		case "oa_users":
			plan.Plan = append(plan.Plan, OATask{
				Command:    "oa_shell",
				Info:       "Creazione home directory da /etc/skel",
				RunCommand: "mkdir -p " + workPath + "/liveroot/home/live && cp -a " + workPath + "/liveroot/etc/skel/. " + workPath + "/liveroot/home/live/",
			})

			plan.Plan = append(plan.Plan, OATask{
				Command:    "oa_users",
				Info:       "Iniezione identità live/live",
				PathLiveFs: workPath,
				Users:      []pilot.User{defaultUser},
			})

		case "oa_umount":
			plan.Plan = append(plan.Plan, OATask{
				Command:    "oa_umount",
				Info:       "Pulizia finale dei mount",
				PathLiveFs: workPath,
			})

		default:
			// Passaggio diretto dallo YAML al JSON (permette a debian.yaml di chiamare qualsiasi verbo C)
			plan.Plan = append(plan.Plan, OATask{
				Command:    step.Command, // <--- Usa il comando originale dello YAML!
				Info:       step.Description,
				RunCommand: step.RunCommand,
				Chroot:     step.Chroot,
				PathLiveFs: workPath,
				Path:       step.Path, // <--- Passa il parametro
				Src:        step.Src,  // <--- Passa il parametro
				Dst:        step.Dst,  // <--- Passa il parametro
			})
		}
	}

	// Scrittura del file JSON finale
	return savePlan(plan)
}

// Aggiunto (string, error) alla firma
func savePlan(plan OAPlan) (string, error) {
	// Definiamo chiaramente dove andrà a finire
	targetDir := "/tmp/coa"
	targetFile := "oa-plan.json"
	fullPath := filepath.Join(targetDir, targetFile)

	// 1. Creiamo la directory e gestiamo l'errore
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", err
	}

	// 2. Marshalling del JSON
	file, err := json.MarshalIndent(plan, "", "  ")
	if err != nil {
		return "", err // Gestiamo l'errore se il JSON è malformato
	}

	// 3. Scrittura del file nel percorso ASSOLUTO
	if err := os.WriteFile(fullPath, file, 0644); err != nil {
		return "", err
	}

	// 4. Se tutto è andato bene, restituiamo il percorso e nessun errore
	return fullPath, nil
}
