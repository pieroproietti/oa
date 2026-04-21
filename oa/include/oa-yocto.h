/*
 * oa: eggs in my dialect🥚🥚
 *
 * include/oa-yocto.h
 * Header per la classificazione utenti basata su OpenEmbedded-Core
 */

#ifndef OA_YOCTO_H
#define OA_YOCTO_H

#include <stdio.h>
#include <stdbool.h>
#include <stdint.h>
#include "cJSON.h"

// --- Definizioni per gli utenti standard (OE-Core) ---
#define OE_UID_HUMAN_MIN 1000
#define OE_UID_HUMAN_MAX 60000

// --- Prototipi delle funzioni ---

/**
 * @brief Scrive una riga in formato passwd (USER:x:UID:GID:GECOS:HOME:SHELL)
 */
void yocto_write_passwd(FILE *f, const char *user, int uid, int gid, const char *gecos, const char *home, const char *shell);

/**
 * @brief Scrive una riga in formato shadow (USER:PASS:LAST:MIN:MAX:WARN:INACT:EXP:RES)
 */
void yocto_write_shadow(FILE *f, const char *user, const char *enc_pass);

/**
 * @brief Scrive una riga in formato group (GROUP:x:GID:USERS)
 */
void yocto_write_group(FILE *f, const char *group, int gid, const char *users);

/**
 * @brief Filtra un file di testo (passwd/group) rimuovendo gli UID/GID umani
 */
int yocto_sanitize_file(const char *src_path, int min_id, int max_id);

/**
 * @brief Filtra il file shadow mantenendo solo gli utenti presenti in passwd
 */
int yocto_sanitize_shadow(const char *shadow_path, const char *passwd_path);

/**
 * @brief Decide se un utente dell'host deve essere processato.
 */
bool yocto_is_human_user(uint32_t uid, const char *home);

/**
 * @brief Aggiunge nativamente un utente a un elenco di gruppi secondari in /etc/group
 */
void yocto_add_user_to_groups(const char *group_file, const char *username, cJSON *groups_array);

#endif // OA_YOCTO_H
