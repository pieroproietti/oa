// Copyright 2026 Piero Proietti <piero.proietti@gmail.com>.
// All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Variabili globali per i flag di Cobra
var (
	modeFlag  string
	pathFlag  string
	cleanFlag bool
)

// checkSudoRequirements verifica se il comando è lanciato con i privilegi corretti.
func checkSudoRequirements(cmdName string, shouldBeRoot bool) {
	isRoot := os.Geteuid() == 0

	if shouldBeRoot && !isRoot {
		fmt.Printf("\n\033[1;31m[ERROR]\033[0m The command 'coa %s' requires root privileges.\n", cmdName)
		fmt.Printf("Please use: \033[1msudo coa %s\033[0m\n\n", cmdName)
		os.Exit(1)
	}

	if !shouldBeRoot && isRoot {
		fmt.Printf("\n\033[1;31m[ERROR]\033[0m Do not run 'coa %s' with sudo.\n", cmdName)
		fmt.Printf("Run it as a normal user: \033[1mcoa %s\033[0m\n\n", cmdName)
		os.Exit(1)
	}
}

func main() {
	// Discovery immediato dell'ambiente (Sensi)
	myDistro := NewDistro()

	// --- ROOT COMMAND ---
	var rootCmd = &cobra.Command{
		Use:   "coa",
		Short: "coa (Cova) - The Artisan Orchestrator",
		Long:  "coa is the intelligent orchestrator written in Go, designed to be the \"Mind\" behind the GNU/Linux system remastering process.",
	}

	// --- PRODUCE COMMAND ---
	var produceCmd = &cobra.Command{
		Use:   "produce",
		Short: "Start a system remastering production flight",
		Run: func(cmd *cobra.Command, args []string) {
			checkSudoRequirements(cmd.Name(), true)
			handleProduce(modeFlag, pathFlag, myDistro)
		},
	}
	produceCmd.Flags().StringVar(&modeFlag, "mode", "standard", "standard, clone, or crypted")
	produceCmd.Flags().StringVar(&pathFlag, "path", "/home/eggs", "working directory")

	// --- EXPORT COMMAND ---
	var exportCmd = &cobra.Command{
		Use:   "export",
		Short: "Export artifacts (iso, pkg) to a remote Proxmox storage",
		Run: func(cmd *cobra.Command, args []string) {
			checkSudoRequirements(cmd.Name(), false)
			handleExport(cleanFlag)
		},
	}
	exportCmd.Flags().BoolVar(&cleanFlag, "clean", false, "remove previous versions")

	// --- KILL COMMAND ---
	var killCmd = &cobra.Command{
		Use:   "kill",
		Short: "Free the nest and unmount filesystems",
		Run: func(cmd *cobra.Command, args []string) {
			checkSudoRequirements(cmd.Name(), true)
			handleKill()
		},
	}

	// --- DETECT COMMAND ---
	var detectCmd = &cobra.Command{
		Use:   "detect",
		Short: "Show host distribution discovery info",
		Run: func(cmd *cobra.Command, args []string) {
			checkSudoRequirements(cmd.Name(), false)
			handleDetect(myDistro)
		},
	}

	// --- KRILL COMMAND ---
	var krillCmd = &cobra.Command{
		Use:   "krill",
		Short: "Start the system installation (The Hatching)",
		Run: func(cmd *cobra.Command, args []string) {
			checkSudoRequirements(cmd.Name(), true)
			handleKrill()
		},
	}

	// --- ADAPT COMMAND ---
	var adaptCmd = &cobra.Command{
		Use:   "adapt",
		Short: "Adapt monitor resolution for VMs",
		Run: func(cmd *cobra.Command, args []string) {
			checkSudoRequirements(cmd.Name(), false)
			handleAdapt()
		},
	}

	// --- VERSION COMMAND ---
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of coa",
		Run: func(cmd *cobra.Command, args []string) {
			checkSudoRequirements(cmd.Name(), false)
			fmt.Printf("coa v%s - The Mind of remaster\n", AppVersion)
		},
	}

	// --- DOCS COMMAND ---
	var docsCmd = &cobra.Command{
		Use:    "docs",
		Short:  "Generate man pages, markdown wiki, and completion scripts",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			checkSudoRequirements(cmd.Name(), false)
			generateDocs(rootCmd)
		},
	}

	// --- BUILD COMMAND ---
	var buildCmd = &cobra.Command{
		Use:   "build",
		Short: "Compile binaries and generate native distribution packages (.deb, PKGBUILD)",
		Run: func(cmd *cobra.Command, args []string) {
			checkSudoRequirements(cmd.Name(), false)
			handleBuild(myDistro)
		},
	}

	// Registrazione comandi
	rootCmd.AddCommand(
		adaptCmd,
		buildCmd,
		detectCmd,
		docsCmd,
		exportCmd,
		killCmd,
		krillCmd,
		produceCmd,
		versionCmd,
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
