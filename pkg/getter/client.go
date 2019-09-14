package getter

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/avinor/tau/pkg/helpers/paths"
	"github.com/avinor/tau/pkg/helpers/ui"
	"github.com/hashicorp/go-getter"
)

// Options for initialization a new getter client
type Options struct {
	Timeout          time.Duration
	WorkingDirectory string
}

// Client used to download or copy source files with. Support all features
// that go-getter supports, local files, http, git etc.
type Client struct {
	options   *Options
	detectors []getter.Detector
	getters   map[string]getter.Getter
}

var (
	// defaultTimeout for context to retrieve the sources
	defaultTimeout = 10 * time.Second
)

// New creates a new getter client. It configures all the detectors and getters itself to make
// sure they are configured correctly.
func New(options *Options) *Client {
	if options == nil {
		options = &Options{}
	}

	if options.WorkingDirectory == "" {
		options.WorkingDirectory = paths.WorkingDir
	}

	if options.Timeout == 0 {
		options.Timeout = defaultTimeout
	}

	httpClient := &http.Client{
		Timeout: options.Timeout,
	}

	registryDetector := &RegistryDetector{
		httpClient: httpClient,
	}

	detectors := []getter.Detector{
		registryDetector,
		new(getter.GitHubDetector),
		new(getter.GitDetector),
		new(getter.BitBucketDetector),
		new(getter.S3Detector),
		new(getter.GCSDetector),
		new(getter.FileDetector),
	}

	httpGetter := &getter.HttpGetter{
		Netrc: true,
	}

	getters := map[string]getter.Getter{
		"file":  &LocalGetter{FileGetter: getter.FileGetter{Copy: true}},
		"git":   new(getter.GitGetter),
		"gcs":   new(getter.GCSGetter),
		"hg":    new(getter.HgGetter),
		"s3":    new(getter.S3Getter),
		"http":  httpGetter,
		"https": httpGetter,
	}

	return &Client{
		options:   options,
		detectors: detectors,
		getters:   getters,
	}
}

// Clone creates a clone of the client with an alternative new working directory.
// If workingDir is set to "" it will just reuse same directory in client
func (c *Client) Clone(workingDir string) *Client {
	if workingDir == "" {
		workingDir = c.options.WorkingDirectory
	}

	return &Client{
		options: &Options{
			WorkingDirectory: workingDir,
			Timeout:          c.options.Timeout,
		},
		detectors: c.detectors,
		getters:   c.getters,
	}
}

// Get retrieves sources from src and load them into dst folder. If version is set it will try to
// download from terraform registry. Set to nil to disable this feature.
func (c *Client) Get(src, dst string, version *string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.options.Timeout)
	defer cancel()

	if version != nil && *version != "" {
		src = fmt.Sprintf("%s?registryVersion=%s", src, *version)
	}

	ui.Info("- %v", src)

	client := &getter.Client{
		Ctx:       ctx,
		Src:       src,
		Dst:       dst,
		Pwd:       c.options.WorkingDirectory,
		Mode:      getter.ClientModeAny,
		Detectors: c.detectors,
		Getters:   c.getters,
	}

	return client.Get()
}

// GetFile retrieves a single file from src destination. It implements almost same
// functionallity that Get function, but does only allow a single file to be downloaded.
func (c *Client) GetFile(src, dst string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.options.Timeout)
	defer cancel()

	ui.Info("- %v", src)

	client := &getter.Client{
		Ctx:       ctx,
		Src:       src,
		Dst:       dst,
		Pwd:       c.options.WorkingDirectory,
		Mode:      getter.ClientModeFile,
		Detectors: c.detectors,
		Getters:   c.getters,
	}

	return client.Get()
}

// Detect is a wrapper on go-getter detect and will return a new source string
// that is the parsed url using correct getter
func (c *Client) Detect(src string) (string, error) {
	return getter.Detect(src, c.options.WorkingDirectory, c.detectors)
}
