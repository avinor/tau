package loader

import (
	"fmt"
	"testing"

	"github.com/avinor/tau/pkg/config"
	"github.com/stretchr/testify/assert"
)

var (
	modA = &ParsedFile{File: &config.File{Name: "A"}}
	modB = &ParsedFile{File: &config.File{Name: "B"}}
	modC = &ParsedFile{File: &config.File{Name: "C"}}
	modD = &ParsedFile{File: &config.File{Name: "D"}, Dependencies: map[string]*ParsedFile{"modA": modA}}
	modE = &ParsedFile{File: &config.File{Name: "E"}, Dependencies: map[string]*ParsedFile{"modA": modA}}
	modG = &ParsedFile{File: &config.File{Name: "G"}, Dependencies: map[string]*ParsedFile{"modA": modA}}
	modI = &ParsedFile{File: &config.File{Name: "I"}, Dependencies: map[string]*ParsedFile{"modG": modG}}
	modK = &ParsedFile{File: &config.File{Name: "K"}, Dependencies: map[string]*ParsedFile{"modA": modA, "modG": modG}}
)

func TestDependencySorting(t *testing.T) {
	tests := []struct {
		Input   []*ParsedFile
		Expects []*ParsedFile
	}{
		{
			[]*ParsedFile{modA, modG},
			[]*ParsedFile{modA, modG},
		},
		{
			[]*ParsedFile{modG, modA},
			[]*ParsedFile{modA, modG},
		},
		{
			[]*ParsedFile{modA, modB, modC},
			[]*ParsedFile{modA, modB, modC},
		},
		{
			[]*ParsedFile{modC, modB, modA},
			[]*ParsedFile{modA, modB, modC},
		},
		{
			[]*ParsedFile{modA, modD, modE},
			[]*ParsedFile{modA, modD, modE},
		},
		{
			[]*ParsedFile{modD, modE, modA},
			[]*ParsedFile{modA, modD, modE},
		},
		{
			[]*ParsedFile{modA, modG, modI},
			[]*ParsedFile{modA, modG, modI},
		},
		{
			[]*ParsedFile{modG, modA, modI},
			[]*ParsedFile{modA, modG, modI},
		},
		{
			[]*ParsedFile{modI, modG, modA},
			[]*ParsedFile{modA, modG, modI},
		},
		{
			[]*ParsedFile{modI, modG, modA, modK},
			[]*ParsedFile{modA, modG, modI, modK},
		},
		{
			[]*ParsedFile{modA, modG, modI, modK},
			[]*ParsedFile{modA, modG, modI, modK},
		},
		{
			[]*ParsedFile{modA, modK, modI, modG},
			[]*ParsedFile{modA, modG, modI, modK},
		},
		{
			[]*ParsedFile{modK, modI, modG, modA},
			[]*ParsedFile{modA, modG, modI, modK},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			sortFiles(test.Input)

			assert.Equal(t, test.Expects, test.Input)
		})
	}
}
