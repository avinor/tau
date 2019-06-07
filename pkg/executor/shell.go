package executor

type Shell struct{}

type Config struct{}

func NewShell(config *Config) (*Shell, error) {
	return nil, nil
}

func (shell *Shell) ExecuteTerraform(cmd string, args ...string) error {
	return nil
}
