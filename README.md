# oa-tools - The Next Generation Remastering Suite 🐧

Welcome to **oa-tools**, born from `penguins-eggs` experience. This monorepo hosts a split-responsibility system designed for high-performance Linux remastering, following the "Universal Strategy" of absolute portability.

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

### 🧠 [coa (brooding in my dialect)](./coa) - The brain
It manages the full lifecycle: from laying the ISO to the final installation.

The name derives from the dialect word coa, referring to the act of brooding or incubating eggs until they are ready to hatch.

**Language: Go**

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