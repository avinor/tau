// Package terraform initializes an engine that execute terraform commands and
// generate terraform code. It only deals with interfaces describing what terraform
// should do, and not the specific implementation. Each version of terraform should
// have its own implementation, or reuse older one if no changes are made to file
// format.
package terraform