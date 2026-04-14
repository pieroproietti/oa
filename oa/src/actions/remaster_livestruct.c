#include "oa.h"
#include <sys/utsname.h>

int remaster_livestruct(OA_Context *ctx) {
    // 1. Recupero percorsi
    cJSON *path_obj = cJSON_GetObjectItemCaseSensitive(ctx->task, "pathLiveFs");
    if (!path_obj) path_obj = cJSON_GetObjectItemCaseSensitive(ctx->root, "pathLiveFs");
    if (!cJSON_IsString(path_obj)) return 1;

    const char *work_base = path_obj->valuestring;
    char liveroot_boot[PATH_SAFE];
    char iso_live_dir[PATH_SAFE];

    // Definiamo i percorsi basandoci sulla nostra area di lavoro
    snprintf(liveroot_boot, sizeof(liveroot_boot), "%s/liveroot/boot", work_base);
    snprintf(iso_live_dir, sizeof(iso_live_dir), "%s/iso/live", work_base);

    // Prepariamo la directory di destinazione
    char cmd[CMD_MAX];
    snprintf(cmd, sizeof(cmd), "mkdir -p %s", iso_live_dir);
    system(cmd);

    // 2. Identificazione kernel corrente
    struct utsname buffer;
    uname(&buffer);
    char *kver = buffer.release;

    printf("\033[1;34m[oa LIVESTRUCT]\033[0m Transporting boot files to ISO structure...\n");

    // -------------------------------------------------------------------------
    // 3. COPIA DELL'INITRD (Fondamentale!)
    // -------------------------------------------------------------------------
    // Cerchiamo l'initrd.img che abbiamo generato con oa_sys_shell
    snprintf(cmd, sizeof(cmd), "cp %s/initrd.img %s/initrd.img && chmod 644 %s/initrd.img", 
             liveroot_boot, iso_live_dir, iso_live_dir);
    
    if (system(cmd) != 0) {
        LOG_WARN("initrd.img not found, trying fallback to host-named initramfs...");
        snprintf(cmd, sizeof(cmd), "cp %s/initramfs-%s.img %s/initrd.img 2>/dev/null", 
                 liveroot_boot, kver, iso_live_dir);
        system(cmd);
    }

    // -------------------------------------------------------------------------
    // 4. COPIA DEL KERNEL (VMLINUZ)
    // -------------------------------------------------------------------------
    
    // TENTATIVO A: Arch Linux LTS (quello che stai usando tu)
    snprintf(cmd, sizeof(cmd), "cp %s/vmlinuz-linux-lts %s/vmlinuz 2>/dev/null", liveroot_boot, iso_live_dir);
    if (system(cmd) == 0) {
        LOG_INFO("Kernel extracted: vmlinuz-linux-lts");
    } 
    // TENTATIVO B: Arch Linux Standard
    else {
        snprintf(cmd, sizeof(cmd), "cp %s/vmlinuz-linux %s/vmlinuz 2>/dev/null", liveroot_boot, iso_live_dir);
        if (system(cmd) == 0) {
            LOG_INFO("Kernel extracted: vmlinuz-linux");
        }
        // TENTATIVO C: Debian/Ubuntu Style (vmlinuz-VERSION)
        else {
            snprintf(cmd, sizeof(cmd), "cp %s/vmlinuz-%s %s/vmlinuz 2>/dev/null", liveroot_boot, kver, iso_live_dir);
            if (system(cmd) != 0) {
                LOG_ERR("Could not find any kernel in %s", liveroot_boot);
                return 1;
            }
        }
    }

    // Forza permessi corretti per il bootloader
    snprintf(cmd, sizeof(cmd), "chmod 644 %s/vmlinuz", iso_live_dir);
    system(cmd);

    printf("\033[1;32m[oa LIVESTRUCT]\033[0m Kernel and initrd are now in %s\n", iso_live_dir);
    return 0;
}
