#include "oa.h"
#include <sys/wait.h>
#include <unistd.h>

/**
 * @brief Esegue un comando tramite shell.
 * In questa versione semplificata, i mount sono già stati predisposti 
 * da remaster_prepare, quindi ci limitiamo al chroot e all'esecuzione.
 */
int sys_shell(OA_Context *ctx) {
    cJSON *cmd_obj = cJSON_GetObjectItemCaseSensitive(ctx->task, "run_command");
    cJSON *chroot_obj = cJSON_GetObjectItemCaseSensitive(ctx->task, "chroot");
    cJSON *path_obj = cJSON_GetObjectItemCaseSensitive(ctx->root, "pathLiveFs");

    if (!cJSON_IsString(cmd_obj)) {
        LOG_ERR("oa_shell: 'run_command' mancante o non valido.");
        return 1;
    }

    const char *command = cmd_obj->valuestring;
    bool use_chroot = cJSON_IsTrue(chroot_obj);

    if (use_chroot) {
        if (!cJSON_IsString(path_obj)) {
            LOG_ERR("oa_shell: pathLiveFs richiesto per il chroot.");
            return 1;
        }

        char target_root[PATH_SAFE];
        snprintf(target_root, sizeof(target_root), "%s/liveroot", path_obj->valuestring);

        pid_t pid = fork();
        if (pid == 0) { // Processo FIGLIO
            if (chroot(target_root) != 0 || chdir("/") != 0) {
                perror("Errore ingresso chroot");
                _exit(1);
            }
            
            // Eseguiamo il comando. Su Arch /bin/sh è presente grazie all'overlay di /usr
            execl("/bin/sh", "sh", "-c", command, (char *)NULL);
            
            // Se execl fallisce
            _exit(1); 
        }

        // Processo PADRE: attende la fine del comando
        int status;
        waitpid(pid, &status, 0);

        return WIFEXITED(status) ? WEXITSTATUS(status) : 1;
    } else {
        // Esecuzione standard sull'host (per task che non richiedono chroot)
        return system(command);
    }
}