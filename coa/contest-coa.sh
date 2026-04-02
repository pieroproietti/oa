#!/bin/bash

# Abilita nullglob: se un pattern (es. *.c) non trova nulla, restituisce un array vuoto
# invece di restituire la stringa letterale "*.c"
shopt -s nullglob

# Genera un numero casuale tra 000 e 999
RAND_SUFFIX=$(printf "%03d" $((RANDOM % 1000)))
CONTEXT="CONTEXT_COA_${RAND_SUFFIX}.txt"

# Elenco dei file dinamico: bash espanderà in automatico i percorsi
FILES=(
  docs/*.md
  src/*
  go.mod
  go.sum
  coa-contest.sh
  Makefile
  README.md
)

(
  echo '````'
  for f in "${FILES[@]}"; do
    if [ -f "$f" ]; then
      echo "### 📄 FILE: $f"
      
      filename=$(basename "$f")
      ext="${filename##*.}"

      case "$ext" in
        c|h)      lang="c" ;;
        sh)       lang="bash" ;;
        json)     lang="json" ;;
        md)       lang="markdown" ;;
        *)        lang="text" ;;
      esac

      if [ "$filename" == "Makefile" ]; then lang="make"; fi
      if [ "$filename" == "context.sh" ]; then lang="bash"; fi

      echo '```'"$lang"
      cat "$f"
      echo '```'
      echo ""
    else
      echo -e "\033[1;31m⚠️ ERRORE: $f non trovato\033[0m" >&2
    fi
  done
  echo '````'
) > "$CONTEXT"

# Disabilita nullglob per evitare comportamenti inattesi nel resto del terminale
shopt -u nullglob

echo -e "\033[1;32m[oa]\033[0m File \033[1m$CONTEXT\033[0m generato con successo!"
scp "$CONTEXT" artisan@192.168.1.2:/home/artisan
rm "$CONTEXT"