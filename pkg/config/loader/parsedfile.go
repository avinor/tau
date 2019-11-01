package loader

import (
	"path/filepath"

	"github.com/avinor/tau/pkg/config"
	"github.com/avinor/tau/pkg/helpers/paths"
	"github.com/pkg/errors"
	"github.com/zclconf/go-cty/cty"
)

var (
	// filePathMustBeAbsError is returned when a file path is relative
	filePathMustBeAbsError = errors.Errorf("file path must be absolute")
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

	moduleDir string
}

// NewParsedFile creates a new parsed file from input parameters. It does not try to read the file
// on disk, but filename has to be an absolute path to file.
func NewParsedFile(filename string, content []byte, tauDir, cacheDir string) (*ParsedFile, error) {
	if !filepath.IsAbs(filename) {
		return nil, filePathMustBeAbsError
	}

	configFile, err := config.NewFile(filename, content)
	if err != nil {
		return nil, err
	}
	tempDir := paths.Join(tauDir, configFile.Name)
	moduleDir := paths.Join(tempDir, "module")

	configFile.AddToContext("module", cty.ObjectVal(map[string]cty.Value{
		"path": cty.StringVal(moduleDir),
	}))

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

	env["TF_PLUGIN_CACHE_DIR"] = paths.JoinAndCreate(cacheDir, "_plugins")

	if ok, err := cfg.Validate(); !ok {
		return nil, err
	}

	return &ParsedFile{
		File:         configFile,
		TempDir:      tempDir,
		Config:       cfg,
		Env:          env,
		Dependencies: map[string]*ParsedFile{},
		moduleDir:    moduleDir,
	}, nil
}

// ModuleDir returns the module directory where source module is downloaded
func (p ParsedFile) ModuleDir() string {
	return p.moduleDir
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

// IsInitialized returns true if the module has been initialized already
func (p ParsedFile) IsInitialized() bool {
	return paths.IsDir(p.ModuleDir())
}
