package cmd

import (
	"coa/src/internal/pilot"

	"github.com/spf13/cobra"
)

// brainCmd rappresenta il comando principale per la gestione del cervello
var brainCmd = &cobra.Command{
	Use:   "brain",
	Short: "Manage the modular configuration (brain.d)",
	Long:  `The 'brain' command provides tools to validate and manage the distribution-specific YAML configurations located in brain.d.`,
}

// lintCmd esegue la pulizia e la validazione dei file YAML
var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Validate brain structure and add missing headers",
	Long:  `Scan the brain.d directory, verify YAML syntax, and ensure every file has the correct functional header.`,
	Example: `  # Validate and annotate all brain files
  coa brain lint`,
	Run: func(cmd *cobra.Command, args []string) {
		// Non serve sudo per il lint di solito, ma se vuoi essere sicuro:
		CheckSudoRequirements(cmd.Name(), false)
		pilot.RunBrainLint()
	},
}

func init() {
	// Aggiungiamo lint come sotto-comando di brain
	brainCmd.AddCommand(lintCmd)

	// Aggiungiamo brain al comando principale rootCmd
	rootCmd.AddCommand(brainCmd)
}
