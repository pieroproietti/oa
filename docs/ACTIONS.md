# 🐧 oa: Action Reference Manual

Every operation in **oa** is driven by a JSON "Plan." This document defines the available actions, their parameters, and their expected behavior on the system.

---

## 🏗️ 1. action_prepare
**Purpose**: Initializes the Zero-Copy environment using OverlayFS and bind mounts.

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `pathLiveFs` | String | The base directory for the remastering process. |

**Behavior**:
1. Creates the base directory structure: `liveroot/` and `.overlay/` (with `lowerdir`, `upperdir`, and `workdir` inside).
2. Performs a physical copy of `/etc` to the `liveroot`.
3. Bind-mounts root entries (e.g., `/bin`, `/sbin`, `/lib`) in read-only mode using `MS_PRIVATE` propagation.
4. Projects `/usr` and `/var` using **OverlayFS** to allow modifications without touching the host.
5. Bind-mounts kernel API filesystems: `/proc`, `/sys`, `/run`, `/dev` into `liveroot/`.

---

## 👤 2. action_users
**Purpose**: Creates the Live user identity within the `liveroot` independently, without relying on host binaries.

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `pathLiveFs` | String | The base directory for the remastering process. |
| `users` | Array of Objects | Contains user definitions (`login`, `password`, `home`, `shell`, `gecos`). |
| `mode` | String | Operation mode: `""` (default) or `"clone"`. |

**Behavior**:
1. If `mode` is not `"clone"`, purges host identities by sanitizing `/etc/passwd` and `/etc/group` (removing UIDs between 1000 and 59999).
2. Opens `liveroot/etc/passwd` and `liveroot/etc/shadow` directly via C file streams.
3. Writes the new user identities and passwords natively using Yocto-inspired helper functions.
4. Creates the user's home directory and sets ownership to 1000:1000.

---

## ⚙️ 3. action_initrd
**Purpose**: Generates the Initial RAM Disk for the live session via template substitution.

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `pathLiveFs` | String | The base directory for the remastering process. |
| `initrd_cmd` | String | Shell template to generate the initrd (e.g., `mkinitramfs -o {{out}} {{ver}}`). |

**Behavior**:
1. Detects the host's kernel version using the `uname` syscall.
2. Replaces the `{{out}}` placeholder with the target `initrd.img` path.
3. Replaces the `{{ver}}` placeholder with the detected kernel version.
4. Executes the finalized command to build the initramfs.

---

## 💀 4. action_remaster
**Purpose**: Prepares the boot environment and populates Isolinux binaries.

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `pathLiveFs` | String | The base directory for the remastering process. |
| `mode` | String | Operation mode, used to determine user persistence logic logging. |

**Behavior**:
1. Creates `iso/live` and `iso/isolinux` directories.
2. Detects the kernel version and copies `vmlinuz` from `/boot` into the live directory.
3. Copies ISOLINUX binaries and BIOS modules into `iso/isolinux/`.
4. Generates a default `isolinux.cfg` boot menu if it does not already exist.

---

## 📦 5. action_squash
**Purpose**: Compresses the `liveroot` into a high-performance SquashFS image.

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `pathLiveFs`| String | The base directory for the remastering process. |
| `compression` | String | Algorithm (`zstd`, `xz`, `gzip`). Default: `zstd`. |
| `compression_level` | Integer | Compression level (e.g., 1-22 for zstd). Defaults to 3. |
| `exclude_list` | String | Path to a custom exclusion list file. |
| `mode` | String | `""`, `"clone"`, or `"crypted"`. |

**Behavior**:
1. Detects available online CPU cores to pass to `mksquashfs`.
2. Applies session exclusions including `/proc`, `/sys`, `/dev`, `/run`, and `/tmp`.
3. **Logic by Mode**:
   * If `mode` is **NOT** `"clone"`, automatically excludes `home/*` and `root/*`.
4. Uses the specified `exclude_list` if valid, otherwise falls back to `/usr/share/oa/exclusion.list`.
5. Generates the `filesystem.squashfs` with the specified compression options.

---

## 💿 6. action_iso
**Purpose**: Masters the final bootable ISO image.

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `pathLiveFs`| String | The base directory for the remastering process. |
| `volid` | String | The label of the ISO (e.g., `OA_LIVE`). |
| `output_iso` | String | The output filename (e.g., `live-system.iso`). |

**Behavior**:
1. Definitively constructs the `xorriso` command using a large `CMD_MAX` buffer.
2. Configures the ISO with hybrid boot capabilities (`-isohybrid-mbr`) and an ISOLINUX bootloader.
3. Writes the output file to the root of `pathLiveFs`.

---

## 🚀 7. action_run
**Purpose**: Safely executes commands inside the `liveroot` chroot environment.

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `pathLiveFs`| String | The base directory for the remastering process. |
| `run_command`| String | The command binary to