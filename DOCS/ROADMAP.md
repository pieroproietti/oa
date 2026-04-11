# ROADMAP: coa (Orchestrator) 🗺️

The evolution of **coa** to reach functional parity with `penguins-eggs`.

## ✅ Phase 1: The Nest (Completed)
- [x] Base architecture in Go (monorepo).
- [x] Integration with the C engine (`oa`).
- [x] Distribution discovery (Mothers and Derivatives via YAML).
- [x] Dynamic Flight Plan generation (in-memory JSON).

## 🚧 Phase 2: The Brooding (In Progress - v0.5.x)
- [x] **Advanced CLI**: Implementation of `produce`, `kill`, and `status` sub-commands.
- [ ] **Validation Layer**: Pre-flight checks for disk space and root permissions.
- [ ] **Log Streaming**: Clean visualization of logs originating from the C engine.
- [ ] **Custom Excludes**: Management of dynamic exclusion lists.

## 🥚 Phase 3: The Hatching (v0.6.x - v0.8.x)
- [ ] **Wardrobe Integration**: Management of "costumes" (configurations) via Git/YAML.
- [ ] **TUI (Terminal User Interface)**: Integration of `coa dad` and `coa mom` (visual configurators).
- [ ] **Export Tools**: Automation for ISO uploading and checksum generation.

## 🐧 Phase 4: Free Flight (v1.0.0)
- [ ] **Krill Rebirth**: Integration of the new system installer.
- [ ] **Functional Parity**: Achieving all core features of `penguins-eggs`.
- [ ] **Documentation**: Automatically generated Wiki and man pages.