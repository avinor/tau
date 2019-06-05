package api

// ByDependencies sorts a list of modules by their dependencies
type ByDependencies []*Module

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
