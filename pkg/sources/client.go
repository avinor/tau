package sources

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"context"

	"github.com/hashicorp/go-getter"
	log "github.com/sirupsen/logrus"
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
}

// Detectors is the list of detectors used
var Detectors []getter.Detector

const (
	defaultTimeout = 10 * time.Second
)

func init() {
	Detectors = append([]getter.Detector{new(RegistryDetector)}, getter.Detectors...)
}

// New creates a new client
func New(options *Options) *Client {
	if options.HttpClient == nil {
		options.HttpClient = &http.Client{
			Timeout: defaultTimeout,
		}
	}

	if options.WorkingDirectory == "" {
		pwd, err := os.Getwd()
		if err != nil {
			return nil
		}

		options.WorkingDirectory = pwd
	}

	return &Client{
		httpClient: options.HttpClient,
		pwd: options.WorkingDirectory,
	}
}

// Get retrieves sources from src and load them into dst
func (c *Client) Get(src, dst string, version *string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()

	if version != nil {
		src = fmt.Sprintf("%s?registryVersion=%s", src, *version)
	}

	log.Debugf("Getting sources from %v", src)

	client := &getter.Client{
		Ctx:  ctx,
		Src:  src,
		Dst:  dst,
		Pwd:  c.pwd,
		Mode: getter.ClientModeAny,
		Detectors: Detectors,
	}

	return client.Get()
}
