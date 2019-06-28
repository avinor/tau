package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModuleRegexp(t *testing.T) {
	tests := []struct {
		Name  string
		Match bool
	}{
		{"test.hcl", true},
		{"/tmp/test.hcl", true},
		{"test.tau", true},
		{"/tmp/test.tau", true},
		{"test_auto.hcl", false},
		{"test_auto.tau", false},
		{"test.hclr", false},
		{"test.taur", false},
		{"/tmp/hcl", false},
		{"/tmp/hcl.mp3", false},
		{"/tmp/test.HCL", true},
		{"/tmp/test.TAU", true},
		{"/tmp/TEST.TAU", true},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			assert.Equal(t, test.Match, moduleMatchFunc(test.Name), test.Name)
		})
	}
}

func TestAutoRegexp(t *testing.T) {
	tests := []struct {
		Name  string
		Match bool
	}{
		{"test.hcl", false},
		{"/tmp/test.hcl", false},
		{"test.tau", false},
		{"/tmp/test.tau", false},
		{"test_auto.hcl", true},
		{"test_auto.tau", true},
		{"test.hclr", false},
		{"test.taur", false},
		{"/tmp/hcl", false},
		{"/tmp/hcl.mp3", false},
		{"/tmp/test.HCL", false},
		{"/tmp/test.TAU", false},
		{"/tmp/TEST.TAU", false},
		{"TEST_AUTO.tau", true},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			assert.Equal(t, test.Match, autoMatchFunc(test.Name), test.Name)
		})
	}
}

var (
	modA = &Source{Name: "A"}
	modB = &Source{Name: "B"}
	modC = &Source{Name: "C"}
	modD = &Source{Name: "D", Dependencies: map[string]*Source{"modA": modA}}
	modE = &Source{Name: "E", Dependencies: map[string]*Source{"modA": modA}}
	modG = &Source{Name: "G", Dependencies: map[string]*Source{"modA": modA}}
	modI = &Source{Name: "I", Dependencies: map[string]*Source{"modG": modG}}
	modK = &Source{Name: "K", Dependencies: map[string]*Source{"modA": modA, "modG": modG}}
)

func TestDependencySorting(t *testing.T) {
	tests := []struct {
		Sources []*Source
		Expects []*Source
	}{
		{
			[]*Source{modA, modG},
			[]*Source{modA, modG},
		},
		{
			[]*Source{modG, modA},
			[]*Source{modA, modG},
		},
		{
			[]*Source{modA, modB, modC},
			[]*Source{modA, modB, modC},
		},
		{
			[]*Source{modC, modB, modA},
			[]*Source{modA, modB, modC},
		},
		{
			[]*Source{modA, modD, modE},
			[]*Source{modA, modD, modE},
		},
		{
			[]*Source{modD, modE, modA},
			[]*Source{modA, modD, modE},
		},
		{
			[]*Source{modA, modG, modI},
			[]*Source{modA, modG, modI},
		},
		{
			[]*Source{modG, modA, modI},
			[]*Source{modA, modG, modI},
		},
		{
			[]*Source{modI, modG, modA},
			[]*Source{modA, modG, modI},
		},
		{
			[]*Source{modI, modG, modA, modK},
			[]*Source{modA, modG, modI, modK},
		},
		{
			[]*Source{modA, modG, modI, modK},
			[]*Source{modA, modG, modI, modK},
		},
		{
			[]*Source{modA, modK, modI, modG},
			[]*Source{modA, modG, modI, modK},
		},
		{
			[]*Source{modK, modI, modG, modA},
			[]*Source{modA, modG, modI, modK},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			sortSources(test.Sources)

			assert.Equal(t, test.Expects, test.Sources)
		})
	}
}
