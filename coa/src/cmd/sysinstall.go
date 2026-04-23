package cmd

import (
	"coa/src/internal/calamares" // <-- Importa il nuovo package
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var sysinstallCmd = &cobra.Command{
	Use:   "sysinstall",
	Short: "Launch the graphical system installer (Calamares + OA)",
	Long: `The 'sysinstall' command prepares the environment and launches 
the Calamares graphical installer. Once the GUI finishes partitioning 
and unpacking, the OA engine will take over to finalize the bootloader.`,
	Example: `  # Launch the system installer
  sudo coa sysinstall`,
	Run: func(cmd *cobra.Command, args []string) {
		CheckSudoRequirements(cmd.Name(), true)

		// Affidiamo tutto al nuovo configuratore
		err := calamares.SetupAndLaunch()
		if err != nil {
			fmt.Printf("\n\033[1;31m[ERRORE CRITICO]\033[0m L'installazione si è interrotta: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("\n\033[1;32m[SUCCESSO]\033[0m Installazione completata!")
	},
}

func init() {
	rootCmd.AddCommand(sysinstallCmd)
}
