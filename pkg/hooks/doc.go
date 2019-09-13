// Package hooks contains a runner that can run hook configured in config file.
// It uses an abstraction for different executors to make it possible to run
// different hooks in different executors.
//
// Standard implementation contains an executor for simple commands and scripts.
// Scripts will use go-getter to download the script and run locally.
package hooks
