// src/main.c

#include "cJSON.h"
#include "exec_logic.h"
#include "image_logic.h"
#include "mount_logic.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// Helper per leggere il file JSON
char *read_file(const char *filename) {
  FILE *f = fopen(filename, "rb");
  if (!f)
    return NULL;
  fseek(f, 0, SEEK_END);
  long len = ftell(f);
  fseek(f, 0, SEEK_SET);
  char *data = malloc(len + 1);
  if (data) {
    fread(data, 1, len, f);
    data[len] = '\0';
  }
  fclose(f);
  return data;
}

// Il "Vigile Urbano": smista i verbi ai vari moduli
int execute_verb(cJSON *task) {
  cJSON *command = cJSON_GetObjectItemCaseSensitive(task, "command");
  if (!cJSON_IsString(command) || (command->valuestring == NULL))
    return 1;

  printf("Vitellus: Executing action '%s'...\n", command->valuestring);

  if (strcmp(command->valuestring, "action_prepare") == 0)
    return action_prepare(task);
  if (strcmp(command->valuestring, "action_skeleton") == 0) 
    return action_skeleton(task); // <-- AGGIUNTO!    
  if (strcmp(command->valuestring, "action_cleanup") == 0)
    return action_cleanup(task);
  if (strcmp(command->valuestring, "action_run") == 0)
    return action_run(task);
  if (strcmp(command->valuestring, "action_squash") == 0)
    return action_squash(task);
  if (strcmp(command->valuestring, "action_iso") == 0)
    return action_iso(task);

  fprintf(stderr, "{\"error\": \"Unknown command '%s'\"}\n",
          command->valuestring);
  return 1;
}

int main(int argc, char *argv[]) {
  if (argc < 2) {
    printf("Vitellus Engine v0.2\nUsage: %s <plan.json>\n", argv[0]);
    return 1;
  }

  char *json_data = read_file(argv[1]);
  if (!json_data) {
    fprintf(stderr, "Error: Could not read file %s\n", argv[1]);
    return 1;
  }

  cJSON *json = cJSON_Parse(json_data);
  if (!json) {
    fprintf(stderr, "Error: Invalid JSON format\n");
    free(json_data);
    return 1;
  }

  cJSON *plan = cJSON_GetObjectItemCaseSensitive(json, "plan");
  cJSON *global_path = cJSON_GetObjectItemCaseSensitive(json, "pathLiveFs");
  int final_status = 0;

  // --- LOGICA DEL PIANO DI VOLO ---
  if (cJSON_IsArray(plan)) {
    cJSON *task;
    cJSON_ArrayForEach(task, plan) {
      // Ereditarietà automatica: se il task non ha pathLiveFs, usa quella del
      // nido (globale)
      if (global_path &&
          !cJSON_GetObjectItemCaseSensitive(task, "pathLiveFs")) {
        cJSON_AddItemToObject(task, "pathLiveFs",
                              cJSON_Duplicate(global_path, 1));
      }

      if (execute_verb(task) != 0) {
        fprintf(
            stderr,
            "{\"status\": \"halted\", \"error\": \"Plan failed at task\"}\n");
        final_status = 1;
        break; // Ci fermiamo se un passo del piano fallisce
      }
    }
  } else {
    // Supporto per il vecchio modo: comando singolo nel root del JSON
    final_status = execute_verb(json);
  }

  cJSON_Delete(json);
  free(json_data);
  return final_status;
}
