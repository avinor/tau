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

// Options for initialization
type Options struct {
	HttpClient       *http.Client //nolint:golint
	WorkingDirectory string
}

// Client to retrieve sources
type Client struct {
	httpClient *http.Client
	pwd        string
	detectors  []getter.Detector
	getters    map[string]getter.Getter
}

const (
	defaultTimeout = 10 * time.Second
)

// New creates a new getter client. It configures all the detectors and getters itself to make
// sure they are configured correctly.
func New(options *Options) *Client {
	if options == nil {
		options = &Options{}
	}

	if options.HttpClient == nil {
		options.HttpClient = &http.Client{
			Timeout: defaultTimeout,
		}
	}

	if options.WorkingDirectory == "" {
		options.WorkingDirectory = paths.WorkingDir
	}

	registryDetector := &RegistryDetector{
		httpClient: options.HttpClient,
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
		httpClient: options.HttpClient,
		pwd:        options.WorkingDirectory,
		detectors:  detectors,
		getters:    getters,
	}
}

// Get retrieves sources from src and load them into dst
func (c *Client) Get(src, dst string, version *string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()

	if version != nil && *version != "" {
		src = fmt.Sprintf("%s?registryVersion=%s", src, *version)
	}

	ui.Info("- %v", src)

	client := &getter.Client{
		Ctx:       ctx,
		Src:       src,
		Dst:       dst,
		Pwd:       c.pwd,
		Mode:      getter.ClientModeAny,
		Detectors: c.detectors,
		Getters:   c.getters,
	}

	return client.Get()
}

// Detect is a wrapper on go-getter detect and will return the location for source
func (c *Client) Detect(src string) (string, error) {
	return getter.Detect(src, c.pwd, c.detectors)
}
