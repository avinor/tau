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

	// Reproduce bug in sorting
	modHub   = &ParsedFile{File: &config.File{Name: "Hub"}}
	modLogs  = &ParsedFile{File: &config.File{Name: "Logs"}}
	modSp    = &ParsedFile{File: &config.File{Name: "Sp"}}
	modKV    = &ParsedFile{File: &config.File{Name: "KV"}, Dependencies: map[string]*ParsedFile{"hub": modHub, "logs": modLogs, "sp": modSp}}
	modReg   = &ParsedFile{File: &config.File{Name: "Reg"}}
	modSpoke = &ParsedFile{File: &config.File{Name: "Spoke"}, Dependencies: map[string]*ParsedFile{"logs": modLogs, "hub": modHub}}
	modGw    = &ParsedFile{File: &config.File{Name: "GW"}, Dependencies: map[string]*ParsedFile{"logs": modLogs, "spoke": modSpoke}}
	modAKS   = &ParsedFile{File: &config.File{Name: "AKS"}, Dependencies: map[string]*ParsedFile{"logs": modLogs, "spoke": modSpoke, "reg": modReg, "gw": modGw, "sp": modSp}}
)

// TestCollectionVisit tests that all nodes are visisted correctly,
// but ignores testing order. That is assumed handled by the graph
// library from terraform. Should be processed in correct order.
func TestCollectionVisit(t *testing.T) {
	tests := []struct {
		Input   ParsedFileCollection
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
		{
			[]*ParsedFile{modKV, modAKS, modGw, modReg, modSpoke},
			[]*ParsedFile{modKV, modReg, modSpoke, modGw, modAKS},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			visited := []*ParsedFile{}

			err := test.Input.Walk(func(file *ParsedFile) error {
				visited = append(visited, file)
				return nil
			})

			assert.NoError(t, err)
			assert.ElementsMatch(t, test.Expects, visited)
		})
	}
}
