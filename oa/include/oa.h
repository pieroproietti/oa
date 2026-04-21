/*
 * include/oa.h
 * La Fondamenta: Solo definizioni globali, zero logica di modulo.
 */
#ifndef OA_H
#define OA_H

// 1. Inclusioni di Sistema Universali (quelle che servono a TUTTI)
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <limits.h>
#include <stdbool.h>
#include <stdint.h>
#include <errno.h>

// 2. Librerie Esterne Globali
#include "cJSON.h"
#include "logger.h"
#include "helpers.h" // Se hai funzioni generiche tipo trim()

// 3. Costanti Globali
#define PATH_INPUT PATH_MAX
#define PATH_OUT   8192
#define CMD_MAX    32768
#define PATH_SAFE  8192

// 4. Tipi Globali (La Mente)
typedef struct {
    cJSON *root;    // Il JSON intero (configurazione globale) 
    cJSON *task;    // Il comando specifico nel plan (configurazione locale) 
} OA_Context;

// FINE. NON AGGIUNGERE NESSUN #include "oa_mount.h" O SIMILI QUI!

#endif
