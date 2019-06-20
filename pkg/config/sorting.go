package config

// ByDependencies sorts a list of sources by their dependencies
type ByDependencies []*Source

func (a ByDependencies) Len() int {
	return len(a)
}

func (a ByDependencies) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByDependencies) Less(i, j int) bool {

	for _, dep := range a[j].Dependencies {
		if dep == a[i] {
			return true
		}
	}

	return false
}