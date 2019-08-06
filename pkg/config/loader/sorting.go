package loader

import "sort"

// sortFiles sorts the files in order of dependencies.
func sortFiles(files []*ParsedFile) {
	sort.SliceStable(files, func(i, j int) bool {
		return lessCompareFiles(files[i], files[j])
	})
}

// lessCompareFiles takes 2 files and compare them to see if j is less than
// i. It runs recursively to check dependencies of dependency against i. Otherwise
// it will sort them incorrectly
//
// If there is no dependency that it can sort by it sorts alphabetically on name.
// This is just to get a consistent order of all elements.
func lessCompareFiles(i, j *ParsedFile) bool {
	for _, dep := range j.Dependencies {
		if dep == i || lessCompareFiles(i, dep) {
			return true
		}
	}

	return i.Name < j.Name
}
