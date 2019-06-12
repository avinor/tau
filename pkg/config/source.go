package config

// Source for one file loaded
type Source struct {
	Hash         string
	File         string
	Content      []byte
	Dependencies map[string]*Module

	Config *Config
}

// ByDependencies sorts a list of sources by their dependencies
type ByDependencies []*Source

func (a ByDependencies) Len() int {
	return len(a)
}

func (a ByDependencies) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByDependencies) Less(i, j int) bool {

	for _, dep := range a[j].deps {
		if dep == a[i] {
			return true
		}
	}

	return false
}

func (src *Source) ModuleDirectory() string {
	return ""
}

func (src *Source) CreateBackendFile() error {
	return nil
}
