package pilot

// Task rappresenta l'unità di lavoro che finirà nel piano di remaster
type Task struct {
	Name        string            `json:"Name"`
	Files       map[string]string `json:"Files"`
	Commands    []string          `json:"Commands"`
	Chroot      bool              `json:"Chroot"`
	Description string            `json:"Description"`
}

type BrainProfile struct {
	Tasks []Task
}

// Strutture di configurazione (mappano i file YAML piatti)

type IdentityConfig struct {
	AdminGroup string   `yaml:"admin_group"`
	UserGroups []string `yaml:"user_groups"`
}

type InitrdConfig struct {
	Command string            `yaml:"command"`
	Files   map[string]string `yaml:"setup_files"` // Usa setup_files per tutti
}

type BootConfig struct {
	Params string `yaml:"params"`
}

type LayoutConfig struct {
	Links map[string]string `yaml:"links"`
}
