# Genera un numero casuale tra 000 e 999
RAND_SUFFIX=$(printf "%03d" $((RANDOM % 1000)))
CONTEST="CONTEST_${RAND_SUFFIX}.txt"

(
  echo '````'
  for f in CHANGELOG.md \
            README.md \
            include/oa.h \
            src/main.c \
            src/actions/action_initrd.c \
            src/actions/action_iso.c \
            src/actions/action_prepare.c \
            src/actions/action_remaster.c \
            src/actions/action_run.c \
            src/actions/action_scan.c \
            src/actions/action_squash.c \
            src/actions/action_users.c; 
    do
    if [ -f "$f" ]; then
      echo "### 📄 FILE: $f"
      # Determina l'estensione per l'evidenziazione del codice
      ext="${f##*.}"
      if [ "$ext" == "c" ] || [ "$ext" == "h" ]; then lang="c"; else lang="markdown"; fi
      echo '```'$lang
      cat "$f"
      echo '```'
      echo ""
    else
      echo "⚠️ ERRORE: $f non trovato"
    fi
  done
  echo '````'
) > $CONTEST

echo -e "\033[1;32m[oa]\033[0m File \033[1m$CONTEST\033[0m generato con successo!"
scp $CONTEST artisan@192.168.1.2:/home/artisan
rm $CONTEST