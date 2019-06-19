package loader

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"context"

	"github.com/hashicorp/go-getter"
	log "github.com/sirupsen/logrus"
)

// GetterOptions for initialization
type GetterOptions struct {
	HttpClient *http.Client
}

// GetterClient to retrieve sources
type GetterClient struct {
	httpClient *http.Client
	detectors []getter.Detector
}

const (
	defaultTimeout = 10 * time.Second
)

// NewGetter creates a new getter client
func NewGetter(options *GetterOptions) *GetterClient {
	if options == nil {
		options = &GetterOptions{}
	}

	if options.HttpClient == nil {
		options.HttpClient = &http.Client{
			Timeout: defaultTimeout,
		}
	}

	registryDetector := &RegistryDetector{
		httpClient: options.HttpClient,
	}

	return &GetterClient{
		httpClient: options.HttpClient,
		detectors: append([]getter.Detector{registryDetector}, getter.Detectors...),
	}
}

// Get retrieves sources from src and load them into dst
func (c *GetterClient) Get(src, dst string, pwd, version *string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()

	if pwd == nil {
		cpwd, err := os.Getwd()
		if err != nil {
			return nil
		}

		pwd = &cpwd
	}

	if version != nil {
		src = fmt.Sprintf("%s?registryVersion=%s", src, *version)
	}

	log.Debugf("Getting sources from %v", src)

	client := &getter.Client{
		Ctx:  ctx,
		Src:  src,
		Dst:  dst,
		Pwd:  *pwd,
		Mode: getter.ClientModeAny,
		Detectors: c.detectors,
	}

	return client.Get()
}
