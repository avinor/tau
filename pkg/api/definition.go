package api

// // Definition of the api
// type Definition struct {
// 	tempDir string
// 	config  *Config
// }

// // Config parameters for configuring api definition
// type Config struct {
// 	Source           string
// 	WorkingDirectory string
// 	ExtraArguments   []string
// }

// // New returns a new api definition
// func New(config *Config) (*Definition, error) {
// 	if config.Source == "" {
// 		return nil, errors.Errorf("Source is empty")
// 	}

// 	if config.WorkingDirectory == "" {
// 		pwd, err := os.Getwd()
// 		if err != nil {
// 			return nil, err
// 		}
// 		config.WorkingDirectory = pwd
// 		log.Debugf("Current working directory: %v", pwd)
// 	}

// 	return &Definition{
// 		tempDir: filepath.Join(config.WorkingDirectory, ".tau", hash(config.Source)),
// 		config:  config,
// 	}, nil
// }

// func (d *Definition) LoadSources(maxDependencyDepth int) error {
// 	return nil
// }

// func (d *Definition) LoadSettings() error {
// 	// Save to file
// 	return nil
// }

// func (d *Definition) SaveSettings() error {
// 	// Load file, if exists
// 	return nil
// }

// func (d *Definition) InitTerraform() error {
// 	// Travers modules and run terraform init
// 	return nil
// }

// func (d *Definition) RunTerraform(cmd string, args ...string) error {
// 	// Run terraform <cmd>
// 	return nil
// }

// func (d *Definition) CreateValues() error {
// 	return nil
// }

// // Run a terraform command on loaded modules, init is handled special
// func (d *Definition) Run(cmd string) error {
// 	runner, err := commands.GetRunner(cmd)
// 	if err != nil {
// 		return err
// 	}

// 	return runner.Run()
// }

// func (d *Definition) runInit() error {
// 	loader, err := newLoader(d.config)
// 	if err != nil {
// 		return err
// 	}

// 	modules, err := loader.load()
// 	if err != nil {
// 		return err
// 	}

// 	if modules == nil {
// 		return nil
// 	}

// 	return nil
// }
