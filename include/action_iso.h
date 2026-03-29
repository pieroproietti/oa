#ifndef ACTION_ISO_H
#define ACTION_ISO_H

#include "cJSON.h"

/**
 * @brief Finalizza la ISO avviabile usando xorriso
 * @param json L'oggetto JSON contenente i parametri dell'azione
 * @return 0 in caso di successo, 1 in caso di errore
 */
int action_iso(cJSON *json);

#endif
