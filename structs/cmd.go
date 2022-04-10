package structs

type CmdConfig struct {
	Name        string
	Tty         bool
	Interactive bool
	Detach      bool
	ImagePath   string
	Volume      []string
}
