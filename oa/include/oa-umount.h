/*
 * include/actions/oa_mount.h
 */
#ifndef OA_UMOUNT_H
#define OA_UMOUNT_H

// 1. Prende la struttura OA_Context e le basi
#include "oa.h" 

// 2. Dichiara SOLO le funzioni che questo modulo offre agli altri
int oa_umount(OA_Context *ctx);

#endif
