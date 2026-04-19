echo "--- START BRAIN DUMP ---"
for distro_dir in *; do
    if [ -d "$distro_dir" ]; then
        distro_name=$(basename "$distro_dir")
        echo "DISTRO: $distro_name"
        for yaml_file in "$distro_dir"/*.yaml; do
            echo "  FILE: $(basename "$yaml_file")"
            echo "  CONTENT:"
            # Aggiunge uno spazio di indentazione al contenuto per leggibilità
            sed 's/^/    /' "$yaml_file"
            echo "  ---"
        done
    fi
done
echo "--- END BRAIN DUMP ---"
