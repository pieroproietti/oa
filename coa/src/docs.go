package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// generateDocs si occupa di creare in automatico la wiki, le man pages e l'autocompletamento
func generateDocs(rootCmd *cobra.Command) {
	fmt.Println("\033[1;34m[coa docs]\033[0m Generating project documentation...")

	// 1. Generazione Markdown per la Wiki
	wikiPath := "./docs/wiki"
	os.MkdirAll(wikiPath, 0755)
	if err := doc.GenMarkdownTree(rootCmd, wikiPath); err != nil {
		fmt.Printf("\033[1;31mError generating markdown:\033[0m %v\n", err)
	} else {
		fmt.Printf("\033[1;32m[+] Markdown Wiki generated in %s\033[0m\n", wikiPath)
	}

	// 2. Generazione Man Pages
	manPath := "./docs/man"
	os.MkdirAll(manPath, 0755)
	header := &doc.GenManHeader{
		Title:   "COA",
		Section: "1",
		Source:  "penguins-eggs",
		Manual:  "coa - The Artisan Orchestrator",
	}
	if err := doc.GenManTree(rootCmd, header, manPath); err != nil {
		fmt.Printf("\033[1;31mError generating man pages:\033[0m %v\n", err)
	} else {
		fmt.Printf("\033[1;32m[+] Man pages generated in %s\033[0m\n", manPath)
	}

	// 3. Generazione Autocompletion Scripts (Bash, Zsh, Fish)
	compPath := "./docs/completion"
	os.MkdirAll(compPath, 0755)

	rootCmd.GenBashCompletionFile(filepath.Join(compPath, "coa.bash"))
	rootCmd.GenZshCompletionFile(filepath.Join(compPath, "coa.zsh"))
	rootCmd.GenFishCompletionFile(filepath.Join(compPath, "coa.fish"), true)

	fmt.Printf("\033[1;32m[+] Autocompletion scripts generated in %s\033[0m\n", compPath)
}
