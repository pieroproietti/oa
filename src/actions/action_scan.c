#define _XOPEN_SOURCE 500 // Necessario per nftw
#include "cJSON.h"
#include <ftw.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "action_scan.h"


// Usiamo una variabile globale (purtroppo richiesta dall'interfaccia di nftw)
static uint64_t current_total_bytes = 0;
static uint64_t current_file_count = 0;

// Callback: viene eseguita per ogni file/directory trovata
static int scan_callback(const char *fpath, const struct stat *sb, int tflag,
                         struct FTW *ftwbuf) {
  // Escludiamo i link simbolici dal conteggio della dimensione per non contare
  // due volte
  if (tflag == FTW_F) {
    current_total_bytes += sb->st_size;
    current_file_count++;
  }

  // Qui in futuro aggiungeremo: if (is_excluded(fpath)) return
  // FTW_SKIP_SUBTREE;

  return 0; // 0 significa "continua la scansione"
}

/**
 * 
 */
int action_scan(cJSON *json) {
  cJSON *path_obj = cJSON_GetObjectItemCaseSensitive(json, "path");

  if (!cJSON_IsString(path_obj) || (path_obj->valuestring == NULL)) {
    fprintf(stderr, "{\"error\": \"Path non specificato per lo scan\"}\n");
    return 1;
  }

  current_total_bytes = 0;
  current_file_count = 0;

  // nftw(path, callback, nopenfd, flags)
  // FTW_PHYS: non segue i link simbolici
  // FTW_MOUNT: non esce dal filesystem corrente (molto utile per noi!)
  if (nftw(path_obj->valuestring, scan_callback, 64, FTW_PHYS | FTW_MOUNT) ==
      -1) {
    fprintf(stderr, "{\"error\": \"Errore durante la scansione di %s\"}\n",
            path_obj->valuestring);
    return 1;
  }

  // Risposta in JSON per Eggs/Adrian
  cJSON *response = cJSON_CreateObject();
  cJSON_AddStringToObject(response, "status", "ok");
  cJSON_AddStringToObject(response, "path", path_obj->valuestring);
  cJSON_AddNumberToObject(response, "total_bytes", (double)current_total_bytes);
  cJSON_AddNumberToObject(response, "total_mb",
                          (double)current_total_bytes / 1024 / 1024);
  cJSON_AddNumberToObject(response, "file_count", (double)current_file_count);

  char *out = cJSON_PrintUnformatted(response);
  printf("%s\n", out);

  free(out);
  cJSON_Delete(response);
  return 0;
}
