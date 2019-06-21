package getter

import (
	"github.com/avinor/tau/pkg/dir"
	"fmt"
	"net/http"
	"time"
	"context"

	"github.com/hashicorp/go-getter"
	"github.com/apex/log"
)

// Options for initialization
type Options struct {
	HttpClient *http.Client
	WorkingDirectory string
}

// Client to retrieve sources
type Client struct {
	httpClient *http.Client
	pwd string
	detectors []getter.Detector
}

const (
	defaultTimeout = 10 * time.Second
)

// New creates a new getter client
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
		options.WorkingDirectory = dir.Working
	}

	registryDetector := &RegistryDetector{
		httpClient: options.HttpClient,
	}

	return &Client{
		httpClient: options.HttpClient,
		pwd : options.WorkingDirectory,
		detectors: append([]getter.Detector{registryDetector}, getter.Detectors...),
	}
}

// Get retrieves sources from src and load them into dst
func (c *Client) Get(src, dst string, version *string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()

	if version != nil {
		src = fmt.Sprintf("%s?registryVersion=%s", src, *version)
	}

	log.Infof("- %v", src)

	client := &getter.Client{
		Ctx:  ctx,
		Src:  src,
		Dst:  dst,
		Pwd:  c.pwd,
		Mode: getter.ClientModeAny,
		Detectors: c.detectors,
	}

	return client.Get()
}
