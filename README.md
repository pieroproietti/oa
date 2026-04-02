# oa-tools - The Next Generation Remastering Suite 🐧

Welcome to **oa-tools**, the successor to the `penguins-eggs` architecture. This monorepo hosts a split-responsibility system designed for high-performance Linux remastering, following the "Universal Strategy" of absolute portability.

The project is divided into two distinct entities: **oa** (The Arm) and **coa** (The Mind).

---

## 🏗 Project Architecture

We have transitioned to a monorepo structure to ensure perfect synchronization between the engine (oa) and the orchestrator (coa).

### 🦾 [oa (eggs in my dialect)](./oa) - The arm
**Language: C**
`oa` is the low-level engine. It handles the "heavy lifting" of the system:
- Managing OverlayFS and mount points.
- Executing SquashFS compression.
- Building the ISO structures (ISOLINUX/UEFI).
- Interacting directly with the Linux Kernel and system binaries.
- **Philosophy:** Performance, stability, and zero-dependency execution.

### 🧠 [coa (brooding or hatching in my dialect)](./coa) - The brain

The name derives from the dialect word coa, referring to the act of nesting or incubating the eggs until they are ready to hatch.

**Language: Go**
`coa` is the high-level orchestrator, the tools users will use:
- **Discovery:** Automatically identifies the host distribution and family (Debian, Arch, Fedora, etc.).
- **Strategy:** Resolves derivatives via YAML and decides which `initrd` or `bootloader` to use.
- **Planning:** Generates dynamic "Flight Plans" (JSON) based on user intent.
- **Interface:** Manages the CLI experience, sub-commands (`produce`, `kill`, `status`), and future TUI helpers.

---

## 🚀 Getting Started

### Prerequisites
- A Linux system (Debian-based, Arch, or Fedora).
- `gcc` and `make` (for `oa`).
- `golang` 1.21+ (for `coa`).

### Build Everything
From this root directory, simply run:
```bash
make
```

This will compile both binaries:
- `./oa/oa` (The Engine)
- `./coa/coa` (The Orchestrator)

---

## 📜 Philosophy
The **oa-tools** project aims to provide a "Passepartout" for Linux remastering. By separating the **Mind** (Go) from the **Arm** (C), we achieve a clean, maintainable, and incredibly fast workflow that can adapt to any distribution without changing the user experience.

---
*Created with passion by Piero Proietti.*