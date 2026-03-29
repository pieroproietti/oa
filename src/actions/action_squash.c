#include "cJSON.h"
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/utsname.h>

#include "action_squash.h"
#include "helpers.h" // Per append_eggs_exclusion
#include "cJSON.h"

/**
 * @brief Crea il filesystem compresso SquashFS
 */
int action_squash(cJSON *json) {
    cJSON *pathLiveFs = cJSON_GetObjectItemCaseSensitive(json, "pathLiveFs");
    cJSON *comp = cJSON_GetObjectItemCaseSensitive(json, "compression");
    cJSON *comp_lvl = cJSON_GetObjectItemCaseSensitive(json, "compression_level");
    cJSON *exclude_file = cJSON_GetObjectItemCaseSensitive(json, "exclude_list");
    cJSON *mode_item = cJSON_GetObjectItemCaseSensitive(json, "mode");

    if (!cJSON_IsString(pathLiveFs)) return 1;

    const char *mode = (cJSON_IsString(mode_item)) ? mode_item->valuestring : "";

    char final_exclude_path[1024] = "";
    if (cJSON_IsString(exclude_file) && access(exclude_file->valuestring, F_OK) == 0) {
        strncpy(final_exclude_path, exclude_file->valuestring, 1024);
    } else if (access("/usr/share/oa/exclusion.list", F_OK) == 0) {
        strncpy(final_exclude_path, "/usr/share/oa/exclusion.list", 1024);
    }

    long nprocs = sysconf(_SC_NPROCESSORS_ONLN);
    int level = cJSON_IsNumber(comp_lvl) ? comp_lvl->valueint : 3;
    const char *comp_str = cJSON_IsString(comp) ? comp->valuestring : "zstd";

    char liveroot[1024], squash_out[1024];
    snprintf(liveroot, 1024, "%s/liveroot", pathLiveFs->valuestring);
    snprintf(squash_out, 1024, "%s/iso/live/filesystem.squashfs", pathLiveFs->valuestring);

    char session_excludes[4096] = "";
    const char *fexcludes[] = {
        "boot/efi/EFI", "boot/loader/entries/", "etc/fstab", "var/lib/docker/",
        "proc/*", "sys/*", "dev/*", "run/*", "tmp/*"
    };
    for (size_t i = 0; i < 9; i++) append_eggs_exclusion(session_excludes, 4096, fexcludes[i]);

    if (strcmp(mode, "clone") != 0) {
        append_eggs_exclusion(session_excludes, 4096, "home/*");
        append_eggs_exclusion(session_excludes, 4096, "root/*");
    }

    char cmd[8192], comp_opts[256] = "";
    if (strcmp(comp_str, "zstd") == 0) snprintf(comp_opts, 256, "-Xcompression-level %d", level);

    snprintf(cmd, sizeof(cmd), "mksquashfs %s %s -comp %s %s -processors %ld -b 1M -noappend -wildcards", 
             liveroot, squash_out, comp_str, comp_opts, nprocs);

    if (strlen(final_exclude_path) > 0) {
        snprintf(cmd + strlen(cmd), 8192 - strlen(cmd), " -ef %s", final_exclude_path);
    }
    if (strlen(session_excludes) > 0) {
        snprintf(cmd + strlen(cmd), 8192 - strlen(cmd), " -e%s", session_excludes);
    }

    printf("\n\033[1;34m[oa Turbo Squash]\033[0m Cores: %ld | Lvl: %d | Mode: %s\n", nprocs, level, mode);
    return system(cmd);
}

