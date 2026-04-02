package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
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
	// 1. Cerca nel PATH di sistema (es. /usr/local/bin/oa)
	path, err := exec.LookPath("oa")
	if err == nil {
		return path
	}

	// 2. Fallback per lo sviluppo (se siamo nella cartella coa/)
	if _, err := os.Stat("../oa/oa"); err == nil {
		return "../oa/oa"
	}

	// 3. Fallback estremo (se siamo nella root artisan/)
	if _, err := os.Stat("./oa/oa"); err == nil {
		return "./oa/oa"
	}

	return "oa" // Ultima speranza: che sudo lo trovi da solo
}

// handleProduce gestisce la creazione della ISO (ex eggs produce)
func handleProduce(args []string, d *Distro) {
	produceCmd := flag.NewFlagSet("produce", flag.ExitOnError)
	mode := produceCmd.String("mode", "standard", "standard, clone, or crypted")
	workPath := produceCmd.String("path", "/home/eggs", "working directory")
	produceCmd.Parse(args)

	fmt.Printf("\033[1;32m[coa]\033[0m Starting production in \033[1m%s\033[0m mode...\n", *mode)
	
	// Generazione del piano dinamico in memoria (da plan.go)
	flightPlan := GeneratePlan(d, *mode, *workPath)
	
	// Passaggio al motore C
	executePlan(flightPlan)
}

// handleKill gestisce la pulizia (ex eggs kill)
func handleKill() {
	fmt.Println("\033[1;33m[coa]\033[0m Freeing the nest... (Cleaning mounts and temp files)")
	
	oaPath := getOaPath()
	// Invochiamo oa con un argomento diretto di cleanup (oa deve saperlo gestire)
	// o passiamo un mini-piano JSON di sola pulizia.
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

	// File temporaneo in /tmp per evitare problemi di permessi di scrittura locale
	tmpJsonPath := "/tmp/plan_coa_tmp.json"
	err = os.WriteFile(tmpJsonPath, jsonData, 0644)
	if err != nil {
		log.Fatalf("\033[1;31m[coa]\033[0m Temp file error: %v", err)
	}
	defer os.Remove(tmpJsonPath) // Pulizia automatica del JSON alla fine

	oaPath := getOaPath()
	fmt.Printf("\033[1;33m[coa]\033[0m Invoking engine: \033[1m%s\033[0m\n", oaPath)

	// Chiamata definitiva al motore C tramite sudo
	cmd := exec.Command("sudo", oaPath, tmpJsonPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("\n\033[1;31m[coa]\033[0m Engine error: %v", err)
	}
	fmt.Println("\n\033[1;32m[coa]\033[0m Hatching completed successfully.")
}

func printUsage() {
	fmt.Println("\033[1;32mcoa (Cova) - The Artisan Orchestrator\033[0m")
	fmt.Println("\nUsage:")
	fmt.Println("  coa produce [--mode standard|clone|crypted] [--path /path]")
	fmt.Println("  coa kill")
	fmt.Println("  coa detect")
	fmt.Println("  coa version")
}