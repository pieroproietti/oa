#include "cJSON.h"
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/utsname.h>
#include "action_iso.h"

/**
 * @brief Finalizza la ISO avviabile
 */
int action_iso(cJSON *json) {
    cJSON *pathLiveFs = cJSON_GetObjectItemCaseSensitive(json, "pathLiveFs");
    cJSON *volid = cJSON_GetObjectItemCaseSensitive(json, "volume_id");
    cJSON *iso_name = cJSON_GetObjectItemCaseSensitive(json, "filename");

    if (!cJSON_IsString(pathLiveFs)) return 1;

    char iso_root[1024], output_iso[1024];
    snprintf(iso_root, 1024, "%s/iso", pathLiveFs->valuestring);
    snprintf(output_iso, 1024, "%s/%s", pathLiveFs->valuestring, 
             cJSON_IsString(iso_name) ? iso_name->valuestring : "live-system.iso");

    char cmd[8192];
    snprintf(cmd, sizeof(cmd),
             "xorriso -as mkisofs -J -joliet-long -r -l -iso-level 3 "
             "-isohybrid-mbr /usr/lib/ISOLINUX/isohdpfx.bin "
             "-partition_offset 16 -V '%s' "
             "-b isolinux/isolinux.bin -c isolinux/boot.cat "
             "-no-emul-boot -boot-load-size 4 -boot-info-table "
             "-o %s %s/",
             cJSON_IsString(volid) ? volid->valuestring : "OA_LIVE",
             output_iso, iso_root);

    printf("\n\033[1;34m[oa ISO Mode]\033[0m Finalizing ISO...\n");
    return system(cmd);
}
