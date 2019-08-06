package loader

import (
	"path/filepath"

	"github.com/avinor/tau/pkg/config"
	"github.com/avinor/tau/pkg/helpers/paths"
	"github.com/pkg/errors"
)

var (
	// filePathMustBeAbsError is returned when a file path is relative
	filePathMustBeAbsError = errors.Errorf("file path must be absolute")

	// loaded is a map of already loaded ParsedFile. Will always be checked so same file is
	// not loaded twice. Map key is absolute path of file
	loaded = map[string]*ParsedFile{}
)

// ParsedFile is a parsed configuration file. It is a composite of config.File so includes also
// Config() function to return configuration. For parsed file the Config attribute should
// be used instead as that prevents it from parsing the config file multiple times.
// They will both return same result though.
type ParsedFile struct {
	*config.File

	TempDir      string
	Config       *config.Config
	Env          map[string]string
	Dependencies map[string]*ParsedFile
}

// GetParsedFile checks if the file has already been parsed and returns previous parsed file
// or loads the file if not already loaded. Next time this is called with same source file
// it will return a reference to the previous loaded file.
func GetParsedFile(file string, tauDir string) (*ParsedFile, error) {
	if !filepath.IsAbs(file) {
		return nil, filePathMustBeAbsError
	}

	if _, already := loaded[file]; already {
		return loaded[file], nil
	}

	parsed, err := parseFile(file, tauDir)
	if err != nil {
		return nil, err
	}

	loaded[file] = parsed
	return parsed, nil
}

// ModuleDir returns the module directory where source module is downloaded
func (p ParsedFile) ModuleDir() string {
	return paths.Join(p.TempDir, "module")
}

// DependencyDir returns the dependency directory for dependency `dep`
func (p ParsedFile) DependencyDir(dep string) string {
	return paths.JoinAndCreate(p.TempDir, "dep", dep)
}

// OverrideFile returns name of override file that writes backend configuration etc
// into module directory.
func (p ParsedFile) OverrideFile() string {
	return paths.Join(p.ModuleDir(), "tau_override.tf")
}

// PlanFile returns name of plan file when running `terraform plan`.
func (p ParsedFile) PlanFile() string {
	return paths.Join(p.ModuleDir(), "tau.tfplan")
}

// VariableFile returns name of input variable file
func (p ParsedFile) VariableFile() string {
	return paths.Join(p.ModuleDir(), "terraform.tfvars")
}

// parseFile parses the file and returns a newly created ParsedFile. It will create a new
// struct for every call to function.
func parseFile(file string, tauDir string) (*ParsedFile, error) {
	configFile, err := config.NewFile(file)
	if err != nil {
		return nil, err
	}

	if err := AddAutoImports(configFile); err != nil {
		return nil, err
	}

	cfg, err := configFile.Config()
	if err != nil {
		return nil, err
	}

	env, err := cfg.Environment.Parse(configFile.EvalContext())
	if err != nil {
		return nil, err
	}

	if ok, err := cfg.Validate(); !ok {
		return nil, err
	}

	return &ParsedFile{
		File:         configFile,
		TempDir:      paths.Join(tauDir, configFile.Name),
		Config:       cfg,
		Env:          env,
		Dependencies: map[string]*ParsedFile{},
	}, nil
}
