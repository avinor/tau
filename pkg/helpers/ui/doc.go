// Package ui is an abstraction on dealing with writing and reading from ui / cli.
// Do not want to clutter the pkg package with log entries as it makes it difficult
// to reuse. So ui package abstracts it away and leaves it to those use packages to
// implement a handler
//
// Concept is based on github.com/mitchellh/cli package, but changed a bit to include
// headers etc
package ui
