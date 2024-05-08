package containers

type Templates struct {
	Local   LocalTemplates
	Offsite *OffsiteTemplates
}

type LocalTemplates struct {
	Backup  Template `yaml:"backup"`
	Restore Template `yaml:"restore"`
	FileExt string   `yaml:"file-ext"`
}

type OffsiteTemplates struct {
	OffsitePush Template `yaml:"offsite-push"`
	OffsitePull Template `yaml:"offsite-pull"`
}

type Template struct {
	Image   string            `yaml:"image"`
	Cmd     []string          `yaml:"cmd"`
	Env     []string          `yaml:"env"`
	Volumes []string          `yaml:"volumes"`
	Labels  map[string]string `yaml:"labels"`
}
