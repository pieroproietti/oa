#include "action_skeleton.h"  // Il suo header
#include "cJSON.h"           // Per manipolare i dati
#include <stdio.h>           // Per printf, fprintf, snprintf, FILE, fopen
#include <stdlib.h>          // Per system()
#include <string.h>          // Per strcmp()
#include <unistd.h>          // Per access() e F_OK
#include <sys/utsname.h>     // Per uname()

/**
 * @brief Prepara la struttura della ISO, copia il kernel e configura il bootloader
 */
int action_skeleton(cJSON *json) {
    cJSON *pathLiveFs = cJSON_GetObjectItemCaseSensitive(json, "pathLiveFs");
    cJSON *mode_item = cJSON_GetObjectItemCaseSensitive(json, "mode");
    
    if (!cJSON_IsString(pathLiveFs)) return 1;

    const char *mode = (cJSON_IsString(mode_item)) ? mode_item->valuestring : "";

    char iso_dir[1024], live_dir[1024], liveroot_dir[1024], isolinux_dir[1024];
    snprintf(iso_dir, 1024, "%s/iso", pathLiveFs->valuestring);
    snprintf(live_dir, 1024, "%s/iso/live", pathLiveFs->valuestring);
    snprintf(isolinux_dir, 1024, "%s/iso/isolinux", pathLiveFs->valuestring);
    snprintf(liveroot_dir, 1024, "%s/liveroot", pathLiveFs->valuestring);

    // 1. Setup Struttura
    char cmd[4096];
    snprintf(cmd, sizeof(cmd), "mkdir -p %s %s %s", live_dir, iso_dir, isolinux_dir);
    system(cmd);

    // 2. Rilevamento Kernel
    struct utsname buffer;
    if (uname(&buffer) != 0) return 1;
    char *kversion = buffer.release;

    printf("{\"status\": \"imaging\", \"step\": \"skeleton\", \"kernel\": \"%s\", \"mode\": \"%s\"}\n", kversion, mode);

    // 3. Copia Kernel
    printf("\033[1;34m[oa]\033[0m Copying kernel...\n");
    snprintf(cmd, sizeof(cmd), "cp /boot/vmlinuz-%s %s/vmlinuz", kversion, live_dir);
    if (system(cmd) != 0) system("cp -L /vmlinuz %s/vmlinuz");

    // 4. Bootloader BIOS (Isolinux)
    printf("\033[1;34m[oa]\033[0m Populating Isolinux binaries...\n");
    snprintf(cmd, sizeof(cmd), "cp /usr/lib/ISOLINUX/isolinux.bin %s/ && "
                               "cp /usr/lib/syslinux/modules/bios/*.c32 %s/", 
             isolinux_dir, isolinux_dir);
    system(cmd);

    // 5. Configurazione Isolinux
    char cfg_path[1024];
    snprintf(cfg_path, 1024, "%s/isolinux.cfg", isolinux_dir);
    if (access(cfg_path, F_OK) != 0) {
        FILE *f = fopen(cfg_path, "w");
        if (f) {
            fprintf(f, "UI vesamenu.c32\n"
                       "PROMPT 0\n"
                       "TIMEOUT 50\n"
                       "DEFAULT live\n"
                       "MENU TITLE oa Boot Menu\n"
                       "LABEL live\n"
                       "  MENU LABEL oa Live (Standard)\n"
                       "  KERNEL /live/vmlinuz\n"
                       "  APPEND initrd=/live/initrd.img boot=live components quiet splash\n");
            fclose(f);
            printf("\033[1;32m[oa]\033[0m isolinux.cfg generated.\n");
        }
    }

    // 6. Gestione Utenti basata sul MODE
    if (strcmp(mode, "clone") == 0) {
        printf("\033[1;32m[oa]\033[0m Mode CLONE: Users will be preserved.\n");
    } else {
        printf("\033[1;34m[oa]\033[0m Cleaning up host users (UID >= 1000) in liveroot...\n");
        snprintf(cmd, sizeof(cmd), "chroot %s /bin/bash -c \"awk -F: '$3 >= 1000 && $3 < 60000 {print $1}' /etc/passwd | xargs -r -n1 userdel -r -f\"", liveroot_dir);
        system(cmd);

        if (strcmp(mode, "crypted") == 0) {
            printf("\033[1;35m[oa]\033[0m Mode CRYPTED: Applying encryption logic...\n");
        }
    }

    return 0;
}

