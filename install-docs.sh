#!/bin/bash

# coa - Installazione manuali e autocompletamento
# Eseguire con sudo!

if [ "$EUID" -ne 0 ]; then
  echo -e "\033[1;31m[Errore]\033[0m Per favore esegui questo script come root (sudo ./install-docs.sh)"
  exit 1
fi

echo -e "\033[1;34m[coa docs]\033[0m Inizio installazione documentazione e completamenti..."

# --- 1. MAN PAGES ---
MAN_DIR="/usr/share/man/man1"

if [ -d "./docs/man" ]; then
    echo -e "  -> Installazione Man Pages in \033[1m$MAN_DIR\033[0m..."
    mkdir -p "$MAN_DIR"
    cp ./docs/man/*.1 "$MAN_DIR/"
    
    # Aggiorna l'indice dei manuali di sistema
    echo "  -> Aggiornamento database mandb..."
    mandb -q
    echo -e "\033[1;32m[+]\033[0m Man pages installate con successo! Prova a digitare: \033[1mman coa\033[0m"
else
    echo -e "\033[1;33m[!]\033[0m Cartella ./docs/man non trovata. Hai prima eseguito 'coa docs'?"
fi

echo ""

# --- 2. BASH COMPLETION ---
# Il percorso moderno standard per i completamenti bash è /usr/share/bash-completion/completions/
BASH_COMP_DIR="/usr/share/bash-completion/completions"

if [ -f "./docs/completion/coa.bash" ]; then
    echo -e "  -> Installazione Bash Completion in \033[1m$BASH_COMP_DIR\033[0m..."
    mkdir -p "$BASH_COMP_DIR"
    
    # Copiamo il file rimuovendo l'estensione, come vuole lo standard bash-completion
    cp ./docs/completion/coa.bash "$BASH_COMP_DIR/coa"
    
    echo -e "\033[1;32m[+]\033[0m Bash completion installato!"
    echo -e "    \033[1;33mNota:\033[0m Per attivarlo subito nel terminale corrente, esegui: \033[1msource $BASH_COMP_DIR/coa\033[0m"
else
    echo -e "\033[1;33m[!]\033[0m File ./docs/completion/coa.bash non trovato. Hai prima eseguito 'coa docs'?"
fi

# (Opzionale) ZSH e FISH
echo ""
echo -e "\033[1;34m[coa docs]\033[0m Installazione completata."
echo "Se usi ZSH, copia ./docs/completion/coa.zsh in un percorso coperto da \$fpath (es. /usr/local/share/zsh/site-functions/_coa)"
echo "Se usi Fish, copia ./docs/completion/coa.fish in ~/.config/fish/completions/"

