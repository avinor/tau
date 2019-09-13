package def

import "github.com/avinor/tau/pkg/getter"

// Options sent to New function when making a new Runner.
type Options struct {
	Getter   *getter.Client
	CacheDir string
}
