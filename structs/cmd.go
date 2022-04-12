package structs

type CmdRunConfig struct {
	Name        string
	Tty         bool
	Interactive bool
	Detach      bool
	ImagePath   string
	Volume      []string
	Env         []string
}
