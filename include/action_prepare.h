#ifndef ACTION_PREPARE_H
#define ACTION_PREPARE_H

#include "cJSON.h"

// La funzione che prepara i mount
int action_prepare(cJSON *json);

// La funzione che pulisce tutto (Aggiungi questa riga!)
int action_cleanup(cJSON *json);

#endif
