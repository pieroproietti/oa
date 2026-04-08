/*
 * oa: eggs in my dialect🥚🥚
 * remastering core
 *
 * Author: Piero Proietti <piero.proietti@gmail.com>
 * License: GPL-3.0-or-later
 */
#include "oa.h"

// Helper per leggere il file JSON
char *read_file(const char *filename) {
    FILE *f = fopen(filename, "rb");
    if (!f) return NULL;
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

// Il "Vigile Urbano": smista i verbi ai vari moduli tramite OA_Context
int execute_verb(cJSON *root, cJSON *task) {
    cJSON *command = cJSON_GetObjectItemCaseSensitive(task, "command");
    if (!cJSON_IsString(command) || (command->valuestring == NULL)) {
        LOG_ERR("Task without a valid 'command' field found.");
        return 1;
    }

    const char *cmd_name = command->valuestring;
    OA_Context ctx = { .root = root, .task = task };

    LOG_INFO(">>> dispatching to: %s", cmd_name);
    printf("\033[1;34m[oa]\033[0m Executing action '%s'...\n", cmd_name);
    int status = 1; // Default a errore

    // --- FASE 1: LAY (Creazione Uovo / Remastering) ---
    if (strcmp(cmd_name, "lay_prepare") == 0)          status = lay_prepare(&ctx);
    else if (strcmp(cmd_name, "lay_cleanup") == 0)     status = lay_cleanup(&ctx);
    else if (strcmp(cmd_name, "lay_crypted") == 0)     status = lay_crypted(&ctx);
    else if (strcmp(cmd_name, "lay_initrd") == 0)      status = lay_initrd(&ctx);
    else if (strcmp(cmd_name, "lay_iso") == 0)         status = lay_iso(&ctx);
    else if (strcmp(cmd_name, "lay_isolinux") == 0)    status = lay_isolinux(&ctx);
    else if (strcmp(cmd_name, "lay_livestruct") == 0)  status = lay_livestruct(&ctx);
    else if (strcmp(cmd_name, "lay_squash") == 0)      status = lay_squash(&ctx);
    else if (strcmp(cmd_name, "lay_uefi") == 0)        status = lay_uefi(&ctx);
    else if (strcmp(cmd_name, "lay_users") == 0)       status = lay_users(&ctx);

    // --- FASE 2: HATCH (Schiusa / Installazione Fisica) ---
    else if (strcmp(cmd_name, "hatch_partition") == 0) status = hatch_partition(&ctx);
    else if (strcmp(cmd_name, "hatch_format") == 0)    status = hatch_format(&ctx);
    else if (strcmp(cmd_name, "hatch_unpack") == 0)    status = hatch_unpack(&ctx);
    else if (strcmp(cmd_name, "hatch_fstab") == 0)     status = hatch_fstab(&ctx);
    else if (strcmp(cmd_name, "hatch_users") == 0)     status = hatch_users(&ctx);
    else if (strcmp(cmd_name, "hatch_uefi") == 0)      status = hatch_uefi(&ctx);    

    // --- FASE 3: SYS (Utility Generiche) ---
    else if (strcmp(cmd_name, "sys_run") == 0)         status = sys_run(&ctx);
    else if (strcmp(cmd_name, "sys_scan") == 0)        status = sys_scan(&ctx);
    else if (strcmp(cmd_name, "sys_suspend") == 0)     status = sys_suspend(&ctx);
    
    // --- ERRORE ---
    else {
        LOG_ERR("Unknown command requested: %s", cmd_name);
        fprintf(stderr, "Error: Unknown command '%s'\n", cmd_name);
        return 1;
    }

    if (status == 0) {
        LOG_INFO("<<< action %s finished successfully", cmd_name);
    } else {
        LOG_ERR("action %s returned exit status %d", cmd_name, status);
    }

    return status;
}

int main(int argc, char *argv[]) {
    if (argc < 2) {
        printf("oa engine v%s\nUsage: %s <plan.json>\n",OA_VERSION, argv[0]);
        return 1;
    }

    // Inizializziamo il logger subito (es. oa.log per chiarezza)
    oa_init_log("oa.log");
    LOG_INFO("=== STARTING OA ENGINE ===");
    LOG_INFO("Input plan: %s", argv[1]);

    char *json_data = read_file(argv[1]);
    if (!json_data) {
        LOG_ERR("Could not read file: %s", argv[1]);
        fprintf(stderr, "Error: Could not read file %s\n", argv[1]);
        oa_close_log();
        return 1;
    }

    cJSON *json = cJSON_Parse(json_data);
    if (!json) {
        LOG_ERR("JSON parsing failed for %s", argv[1]);
        fprintf(stderr, "Error: Invalid JSON format\n");
        free(json_data);
        oa_close_log();
        return 1;
    }

    cJSON *plan = cJSON_GetObjectItemCaseSensitive(json, "plan");
    int final_status = 0;

    // --- LOGICA DEL PIANO DI VOLO ---
    if (cJSON_IsArray(plan)) {
        LOG_INFO("Plan detected: processing %d tasks", cJSON_GetArraySize(plan));
        cJSON *task;
        int step = 0;
        cJSON_ArrayForEach(task, plan) {
            step++;
            LOG_INFO("--- Task %d ---", step);
            if (execute_verb(json, task) != 0) {
                LOG_ERR("Plan halted at step %d due to previous error", step);
                fprintf(stderr, "Error: Plan halted at step %d\n", step);
                final_status = 1;
                break;
            }
        }
    } else {
        LOG_INFO("No plan array found. Executing root as a single task.");
        final_status = execute_verb(json, json);
    }

    if (final_status == 0) {
        LOG_INFO("=== PLAN COMPLETED SUCCESSFULLY ===");
    } else {
        LOG_ERR("=== PLAN FAILED ===");
    }

    cJSON_Delete(json);
    free(json_data);
    oa_close_log();
    return final_status;
}
