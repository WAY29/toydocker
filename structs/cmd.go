package structs

type CmdConfig struct {
	Tty         bool
	Interactive bool
	Detach      bool
	ImagePath   string
	Volume      []string
}
