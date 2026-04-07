package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// 1. Discovery immediato dell'ambiente (Sensi)
	myDistro := NewDistro()

	if len(os.Args) < 2 {
		printUsage()
		return
	}

	// 2. Gestione dei sotto-comandi (Logica)
	switch os.Args[1] {
	case "produce":
		handleProduce(os.Args[2:], myDistro)
	case "kill":
		handleKill()
	case "detect":
		handleDetect(myDistro)
	case "version":
		fmt.Printf("coa v0.1.0 - The Mind of remaster\n")
	default:
		fmt.Printf("\033[1;31mError:\033[0m Unknown command '%s'\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

// getOaPath cerca il braccio operativo (oa) nel sistema o nel percorso relativo
func getOaPath() string {
	path, err := exec.LookPath("oa")
	if err == nil {
		return path
	}
	if _, err := os.Stat("../oa/oa"); err == nil {
		return "../oa/oa"
	}
	if _, err := os.Stat("./oa/oa"); err == nil {
		return "./oa/oa"
	}
	return "oa"
}

// bridgeConfigs sovrascrive la configurazione mkinitcpio nella liveroot per Arch
func bridgeConfigs(d *Distro, workPath string) error {
	if d.FamilyID != "archlinux" {
		return nil
	}

	presetName := "live-arch.conf"
	if d.DistroID == "manjaro" || d.DistroID == "biglinux" {
		presetName = fmt.Sprintf("live-%s.conf", d.DistroID)
	}

	src := fmt.Sprintf("/tmp/coa/configs/mkinitcpio/%s", presetName)
	// Destinazione: Sovrascriviamo il file standard per evitare parametri extra nel chroot
	dst := filepath.Join(workPath, "liveroot", "etc", "mkinitcpio.conf")

	fmt.Printf("\033[1;34m[coa]\033[0m Overwriting liveroot /etc/mkinitcpio.conf with %s...\n", presetName)

	if err := copyFile(src, dst); err != nil {
		return err
	}

	return os.Chmod(dst, 0644)
}

// copyFile è una utility per la copia fisica tra percorsi host
func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

// handleProduce gestisce la creazione della ISO
func handleProduce(args []string, d *Distro) {
	produceCmd := flag.NewFlagSet("produce", flag.ExitOnError)
	mode := produceCmd.String("mode", "standard", "standard, clone, or crypted")
	workPath := produceCmd.String("path", "/home/eggs", "working directory")
	produceCmd.Parse(args)

	// 1. Estrazione Assets Config sull'host
	tempConfigPath := "/tmp/coa/configs"
	fmt.Printf("\033[1;32m[coa]\033[0m Extracting internal configurations to %s...\n", tempConfigPath)
	if err := ExtractConfigs(tempConfigPath); err != nil {
		log.Fatalf("\033[1;31m[coa]\033[0m Asset extraction failed: %v", err)
	}

	// 2. Download/Verifica Bootloaders (Indispensabile per action_iso)
	fmt.Printf("\033[1;32m[coa]\033[0m Ensuring bootloaders are present in %s...\n", BootloaderRoot)
	if _, err := EnsureBootloaders(); err != nil {
		log.Fatalf("\033[1;31m[coa]\033[0m Bootloader retrieval failed: %v", err)
	}

	// 3. Fase di Preparazione (Crea la liveroot tramite OverlayFS)
	fmt.Printf("\033[1;32m[coa]\033[0m Preparing environment...\n")
	prePlan := FlightPlan{
		PathLiveFs: *workPath,
		Mode:       *mode,
		Plan:       []Action{{Command: "action_prepare"}},
	}
	executePlan(prePlan)

	// 4. Ponte delle configurazioni (Host -> Liveroot)
	if err := bridgeConfigs(d, *workPath); err != nil {
		log.Fatalf("\033[1;31m[coa]\033[0m Config bridging failed: %v", err)
	}

	// 5. Esecuzione del piano di produzione completo
	fmt.Printf("\033[1;32m[coa]\033[0m Starting production flight...\n")
	flightPlan := GeneratePlan(d, *mode, *workPath)
	
	// Rimuoviamo action_prepare dal piano finale per evitare ri-montaggi
	if len(flightPlan.Plan) > 0 && flightPlan.Plan[0].Command == "action_prepare" {
		flightPlan.Plan = flightPlan.Plan[1:]
	}

	executePlan(flightPlan)
}

// handleKill gestisce la pulizia (ex eggs kill)
func handleKill() {
	fmt.Println("\033[1;33m[coa]\033[0m Freeing the nest...")
	oaPath := getOaPath()
	cmd := exec.Command("sudo", oaPath, "cleanup") 
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("\033[1;31m[coa]\033[0m Cleanup failed: %v\n", err)
	}
}

// handleDetect mostra le info rilevate dal discovery
func handleDetect(d *Distro) {
	fmt.Println("\033[1;34m--- coa distro detect ---\033[0m")
	fmt.Printf("Host Distro:     %s\n", d.DistroID)
	fmt.Printf("Family:          %s\n", d.FamilyID)
	fmt.Printf("DistroLike:      %s\n", d.DistroLike)
	fmt.Printf("Codename:        %s\n", d.CodenameID)
	fmt.Printf("Release:         %s\n", d.ReleaseID)
	fmt.Printf("DistroUniqueID:  %s\n", d.DistroUniqueID)
}

// executePlan trasforma il piano in JSON e lo dà in pasto a oa
func executePlan(plan FlightPlan) {
	jsonData, err := json.MarshalIndent(plan, "", "  ")
	if err != nil {
		log.Fatalf("\033[1;31m[coa]\033[0m JSON Error: %v", err)
	}

	tmpJsonPath := "/tmp/plan_coa_tmp.json"
	err = os.WriteFile(tmpJsonPath, jsonData, 0644)
	if err != nil {
		log.Fatalf("\033[1;31m[coa]\033[0m Temp file error: %v", err)
	}
	defer os.Remove(tmpJsonPath)

	oaPath := getOaPath()
	cmd := exec.Command("sudo", oaPath, tmpJsonPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("\n\033[1;31m[coa]\033[0m Engine error: %v", err)
	}
}

func printUsage() {
	fmt.Println("\033[1;32mcoa (Cova) - The Artisan Orchestrator\033[0m")
	fmt.Println("\nUsage:")
	fmt.Println("  coa produce [--mode standard|clone|crypted] [--path /path]")
	fmt.Println("  coa kill")
	fmt.Println("  coa detect")
	fmt.Println("  coa version")
}