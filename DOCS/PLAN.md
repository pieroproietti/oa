# 🗺️ PLAN: coa as a Universal Remastering Middleware

This roadmap outlines the transition of **coa** from a standalone CLI tool to a high-level **Orchestration Middleware (SDK)**. The goal is to provide a standardized YAML-based "Intent Contract" for external GUI incubators, while delegating low-level kernel operations to the **oa** C-engine.

---

### 🟢 Phase 1: SDK Refactoring & API Exposure (Short Term)
* **Decoupling Logic**: Modularize `internal/distro` and `internal/engine` to allow direct imports by external Go projects (e.g., Fyne or Wails-based GUIs).
* **Version Parity**: Implement a unified versioning system using Git tags and `ldflags` to ensure the "Brain" (Go) and "Body" (C) are always synchronized (e.g., v0.6.2).
* **Asset Extraction API**: Expose `internal/assets` as a library function to allow GUIs to trigger the extraction of distribution-specific configurations (dracut/mkinitcpio) into temporary workspaces.

---

### 🟡 Phase 2: The "Intent Contract" (YAML Integration)
* **YAML Schema Definition**: Create a human-readable YAML schema that allows users to define the "remastering intent" (Mode, User, Exclusions, ISO metadata) without technical jargon.
* **The Translator Layer**: Develop a translation engine in `coa` that parses the high-level YAML and enriches it with distro-specific data (Initramfs abstraction, bootloader paths) to generate the final `plan.json` for **oa**.
* **Pre-Flight Validation**: Implement a native validation layer to check for required dependencies (`xorriso`, `squashfs-tools`), disk space, and root privileges before handing execution to the engine.

---

### 🟠 Phase 3: Infrastructure Abstraction (Advanced)
* **Passepartout Automation**: Refine `utils.EnsureBootloaders()` to handle background downloads and verification of hybrid BIOS/UEFI payloads, providing real-time progress callbacks for GUI progress bars.
* **Bridge Configuration**: Automate the physical patching of host configurations (like `/etc/mkinitcpio.conf` for Arch-based distros) within the isolated OverlayFS environment, completely transparent to the GUI.
* **Advanced Exclusion Mapping**: Create a dynamic exclusion engine that combines hardcoded system safety paths (`/proc`, `/sys`, `/dev`) with user-defined patterns from the YAML contract.

---

### 🔴 Phase 4: Full Lifecycle Management (Hatching & Export)
* **Krill Rebirth (API Edition)**: Integrate the **Krill** installer logic into the middleware, allowing GUIs to trigger physical installations by simply providing target disk parameters in the YAML contract.
* **Artifact Export Orchestration**: Standardize the remote export of generated ISOs or packages to Proxmox/remote storages using SSH multiplexing, managed entirely by the middleware.
* **Documentation as a Contract**: Generate an automatically updated Wiki and Man pages that define the YAML specification as the "Universal Language" for any future Linux remastering GUI.

---

## 📜 Technical Philosophy

> **"coa acts as the Mind (Go), oa acts as the Body (C)."**
> The middleware provides **Turbo SquashFS** multi-core compression and **Zero-Copy OverlayFS** layering while shielding the developer from the complexity of kernel syscalls and distribution fragmentation.