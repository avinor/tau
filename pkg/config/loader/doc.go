// Package loader can load and parse configuration files. Unlike the File structure in config
// package that will just load a single file this will also parse the file and load all
// dependencies and load _auto load files that will be merged together with the main configuration.
package loader