package engine

import (
	"coa/src/internal/distro"
	"coa/src/internal/pilot"
	"fmt"
	"os"
	"path/filepath"
)

// GenerateBootMenus scrive i menu di avvio ricavando i parametri dal Cervello (YAML)
func GenerateBootMenus(workPath string, d *distro.Distro, profile *pilot.BrainProfile, volId string) error {
	grubCfgPath := filepath.Join(workPath, "isodir", "boot", "grub", "grub.cfg")
	isolinuxCfgPath := filepath.Join(workPath, "isodir", "isolinux", "isolinux.cfg")

	os.MkdirAll(filepath.Dir(grubCfgPath), 0755)
	os.MkdirAll(filepath.Dir(isolinuxCfgPath), 0755)

	// =================================================================
	// 1. LOGICA DI ESTRAZIONE PARAMETRI (Ora pulitissima grazie ad Areas)
	// =================================================================

	// Default: Debian Style
	// bootParams := "boot=live components quiet splash"
	bootParams := "boot=live components "

	// La regola è semplice: lo YAML vince su tutto.
	// Se non c'è nello YAML, applichiamo i fallback per famiglia.
	if profile.Areas.Boot.Params != "" {
		bootParams = profile.Areas.Boot.Params
	} else if d.FamilyID == "archlinux" || d.FamilyID == "arch" {
		bootParams = "archisobasedir=arch archisolabel=" + volId + " rw"
	}

	// =================================================================
	// 2. GENERAZIONE MENU GRUB (UEFI/Legacy)
	// =================================================================
	grubContent := fmt.Sprintf(`set timeout=5
set default=0

# Caricamento moduli essenziali per lettura filesystem e partizioni
insmod all_video
insmod part_gpt
insmod part_msdos
insmod fat
insmod iso9660
insmod ext2
insmod search_fs_label
insmod search_fs_uuid
insmod loopback

# Ricerca della partizione tramite Label per evitare il rescue prompt
search --no-floppy --set=root --label %s

menuentry "Start OA Live (%s)" {
    linux /live/vmlinuz %s
    initrd /live/initrd.img
}

menuentry "Start OA Live (%s) - RAM mode" {
    linux /live/vmlinuz %s toram
    initrd /live/initrd.img
}
`, volId, d.FamilyID, bootParams, d.FamilyID, bootParams)

	// =================================================================
	// 3. GENERAZIONE MENU ISOLINUX (BIOS)
	// =================================================================
	isolinuxContent := fmt.Sprintf(`UI vesamenu.c32
TIMEOUT 50
DEFAULT live
PROMPT 0

LABEL live
    MENU LABEL Start OA Live (%s)
    LINUX /live/vmlinuz
    APPEND %s
    INITRD /live/initrd.img

LABEL ram
    MENU LABEL Start OA Live (%s) - RAM mode
    LINUX /live/vmlinuz
    APPEND %s toram
    INITRD /live/initrd.img
`, d.FamilyID, bootParams, d.FamilyID, bootParams)

	// =================================================================
	// 4. SCRITTURA SU DISCO
	// =================================================================
	err := os.WriteFile(grubCfgPath, []byte(grubContent), 0644)
	if err != nil {
		return fmt.Errorf("errore scrittura grub.cfg: %w", err)
	}

	err = os.WriteFile(isolinuxCfgPath, []byte(isolinuxContent), 0644)
	if err != nil {
		return fmt.Errorf("errore scrittura isolinux.cfg: %w", err)
	}

	// Trampolino EFI in /EFI/BOOT/grub.cfg (Fisso e strutturale)
	efiTrampolinePath := filepath.Join(workPath, "isodir", "EFI", "BOOT", "grub.cfg")
	os.MkdirAll(filepath.Dir(efiTrampolinePath), 0755)
	trampolino := "search --set=root --file /live/filesystem.squashfs\nset prefix=($root)/boot/grub\nconfigfile $prefix/grub.cfg\n"
	os.WriteFile(efiTrampolinePath, []byte(trampolino), 0644)
	return nil
}
