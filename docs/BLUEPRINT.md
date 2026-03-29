# oa: Technical Blueprint & Architectural Design

**Project Philosophy**: A minimalist, high-performance C engine that provides a "Zero-Copy" remastering environment using OverlayFS. Logic is hardcoded in C for speed and safety; distribution-specific data is abstracted into JSON profiles.

---

## 1. Core Action Pipeline

| Action | JSON Parameters | Debian/Devuan (Default) | Divergence (Arch/Fedora/Suse) |
| :--- | :--- | :--- | :--- |
| **`prepare`** | `pathLiveFs`, `mounts[]` | Mounts `/dev, /proc, /sys, /run` | Standard POSIX (Universal) |
| **`users`** | `name`, `pass`, `groups[]`, `uid` | Admin group: `sudo` | Admin group: `wheel` |
| **`skeleton`** | `kernel_path`, `initrd_template` | `/vmlinuz`, `mkinitramfs` | `vmlinuz-linux`, `dracut/mkinitcpio` |
| **`customize`**| `hooks_path`, `env[]` | Bash scripts in `hooks.d/` | Distro-specific config triggers |
| **`squash`** | `comp`, `level`, `exclude_list` | APT/DPKG exclusions | Pacman/DNF/Zypper exclusions |
| **`iso`** | `volid`, `filename`, `boot_mode` | Isolinux/Grub (BIOS/UEFI) | Systemd-boot or custom EFI paths |
| **`cleanup`** | - | Umount OverlayFS & Bind | Standard POSIX (Universal) |

---

## 2. The "Agnostic" Strategy (The Cut)

To keep the C code clean and "immortal," we avoid hardcoding distro-specific paths. We use **Variable Injection** in the JSON plan.

### A. Kernel & Initrd Abstraction
Instead of `oa` guessing where the kernel is, the `plan.json` provides the template:
* **Debian Template**: `"initrd_cmd": "mkinitramfs -o {{output}} {{version}}"`
* **Arch Template**: `"initrd_cmd": "mkinitcpio -k {{version}} -g {{output}}"`
* **oa's Job**: The C engine only performs string replacement for `{{output}}` (the ISO path) and `{{version}}` (detected via `uname` or provided).

### B. User & Group Management (Yocto Style)
* **The Logic**: `oa` creates the user inside the `chroot` using standard binaries (`useradd`, `chpasswd`).
* **The Data**: The list of groups is passed as an array. `oa` doesn't need to know if `sudo` or `wheel` exists; it simply attempts to add the user to whatever groups are listed in the JSON.

---

## 3. Global Configuration Structure (Example)

```json
{
  "project": "oa",
  "version": "0.1-alpha",
  "profile": "debian-stable",
  "globals": {
    "work_dir": "/home/eggs",
    "kernel_source": "/vmlinuz",
    "initrd_tool": "mkinitramfs -o {{out}} {{ver}}",
    "admin_groups": ["sudo", "audio", "video", "netdev"]
  },
  "plan": [
    { "command": "action_prepare" },
    { 
      "command": "action_users", 
      "user": "oa", 
      "pass": "live", 
      "groups": "{{admin_groups}}" 
    },
    { "command": "action_skeleton" },
    { "command": "action_squash", "compression": "zstd", "level": 3 },
    { "command": "action_iso", "filename": "oa-live.iso" },
    { "command": "action_cleanup" }
  ]
}
```

---

## 4. Implementation Constraints (The "Veteran" Rules)

1. **Zero-Copy Rule**: Never use `cp -a` for the root filesystem. Always use **OverlayFS** to project the host onto the `liveroot`.
2. **Minimalist Dependencies**: Only `libc` and `cJSON`. No high-level wrappers.
3. **Execution Safety**: Every `system()` or `exec()` call must be logged in JSON format to `stdout` for the orchestrator (like eggs) to parse.
4. **Surgical Cleanup**: The `cleanup` action must be resilient. Even if a previous action fails, `cleanup` must ensure no mounts are left "leaking" on the host.

---
*oa architecture is under exploration - 2026*