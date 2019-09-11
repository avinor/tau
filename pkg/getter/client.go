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

// Client to retrieve sources
type Client struct {
	workingDir string
	detectors  []getter.Detector
	getters    map[string]getter.Getter
}

var (
	// Timeout for context to retrieve the sources
	Timeout = 10 * time.Second
)

// New creates a new getter client. It configures all the detectors and getters itself to make
// sure they are configured correctly.
func New(workingDir string) *Client {
	if workingDir == "" {
		workingDir = paths.WorkingDir
	}

	httpClient := &http.Client{
		Timeout: Timeout,
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
		workingDir: workingDir,
		detectors:  detectors,
		getters:    getters,
	}
}

// Clone creates a clone of the client with an alternative new working directory.
// If workingDir is set to "" it will just reuse same directory in client
func (c *Client) Clone(workingDir string) *Client {
	if workingDir == "" {
		workingDir = c.workingDir
	}

	return &Client{
		workingDir: workingDir,
		detectors:  c.detectors,
		getters:    c.getters,
	}
}

// Get retrieves sources from src and load them into dst folder. If version is set it will try to
// download from terraform registry. Set to nil to disable this feature.
func (c *Client) Get(src, dst string, version *string) error {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	if version != nil && *version != "" {
		src = fmt.Sprintf("%s?registryVersion=%s", src, *version)
	}

	ui.Info("- %v", src)

	client := &getter.Client{
		Ctx:       ctx,
		Src:       src,
		Dst:       dst,
		Pwd:       c.workingDir,
		Mode:      getter.ClientModeAny,
		Detectors: c.detectors,
		Getters:   c.getters,
	}

	return client.Get()
}

// Detect is a wrapper on go-getter detect and will return the location for source
func (c *Client) Detect(src string) (string, error) {
	return getter.Detect(src, c.workingDir, c.detectors)
}
