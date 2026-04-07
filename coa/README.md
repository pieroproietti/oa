# coa 🥚

**coa** is the intelligent orchestrator written in Go, designed to be the "Mind" behind the GNU/Linux system remastering process. It acts as a high-level interface for the native **oa** engine, handling system analysis, flight plan generation, and abstracting distribution-specific differences.

## 🚀 Key Features

* **Distro-Agnostic Intelligence**: Automatically recognizes Debian, Arch, Fedora, and their derivatives.
* **Initramfs Abstraction**: Delegates initrd generation to the engine via dynamic templates, supporting `mkinitcpio`, 9`dracut`, and `mkinitramfs`.
* **Passepartout Bootloaders**: Manages the automatic download and injection of hybrid BIOS/UEFI bootloaders to ensure successful booting on any hardware.
* **Bridge Configuration**: Implements physical patching of host configurations, such as `/etc/mkinitcpio.conf`, directly within the isolated working environment.
* **JSON-Driven**: Compiles complex workflows into a standardized JSON format for atomic execution by the C engine.

## 🛠 Prerequisites

* **Engine**: Requires the `oa` binary installed, typically in `/usr/local/bin/oa`.
* **Tools**: `xorriso`, `squashfs-tools`, and `cryptsetup` for encrypted mode.

## 📂 Usage

```bash
# Start a standard ISO production
sudo coa produce --mode standard

# Clean up the workspace and unmount filesystems
sudo coa kill

# Detect host distribution information
coa detect
```

## 🏗 Architecture

**coa** follows a "Mind and Body" philosophy. While the Go orchestrator handles the logic and planning, the C-native engine (**oa**) performs the heavy lifting using kernel syscalls, OverlayFS for zero-copy mirroring, and high-performance compression.

## 🐧 Part of the penguins-eggs Ecosystem

**coa** is a core component of the next-generation *penguins-eggs* infrastructure, aiming for maximum performance and native integration with modern Linux distributions.

---
*Developed with the efficiency of Go and the precision of a remastering artisan.*
