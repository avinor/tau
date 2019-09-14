// Package getter is used to download or copy source files from remote / local
// location. It is a wrapper around go-getter and supports all the same features.
// It also implements a special detector for terraform registry support.
// This allows it to use the official terraform registry to download modules.
package getter
