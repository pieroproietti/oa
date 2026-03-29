CC = gcc
# -Iinclude: fondamentale per trovare cJSON.h e gli altri
CFLAGS = -Wall -Wextra -Iinclude -D_GNU_SOURCE
LDFLAGS = -lm

# 1. Trova tutti i file .c ricorsivamente dentro src/
SOURCES = $(shell find src -name '*.c')

# 2. Trasforma i percorsi dei sorgenti in percorsi per gli oggetti
# Es: src/actions/action_iso.c -> obj/actions/action_iso.o
OBJECTS = $(SOURCES:src/%.c=obj/%.o)

TARGET = oa

all: $(TARGET)

# Regola per il linking finale
$(TARGET): $(OBJECTS)
	@echo "  LD    $@"
	@$(CC) $(OBJECTS) -o $@ $(LDFLAGS)
	@echo "--------------------------------------"
	@echo "Build completata con successo: ./$(TARGET)"

# Regola magica per la compilazione
# @mkdir -p crea le sottocartelle in obj/ se non esistono
obj/%.o: src/%.c
	@mkdir -p $(dir $@)
	@echo "  CC    $<"
	@$(CC) $(CFLAGS) -c $< -o $@

clean:
	@echo "  Cleaning up..."
	@rm -rf obj $(TARGET)

.PHONY: all clean
