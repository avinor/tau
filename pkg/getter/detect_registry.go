package getter

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform/registry"
	"github.com/hashicorp/terraform/registry/regsrc"
)

// RegistryDetector implements detector to detect terraform registry.
// Src have to be formatted with query parameter ?registryVersion=
// avinor/storage-account/azurerm?registryVersion=1.0
type RegistryDetector struct {
	httpClient *http.Client
}

// Detect implements the Detector interface and will check if this source is a terraform registry
// source. If src contains ?registryVersion parameter it will assume it is a registry source.
func (d *RegistryDetector) Detect(src, _ string) (string, bool, error) {
	if len(src) == 0 {
		return "", false, nil
	}

	if strings.Contains(src, "?registryVersion=") {
		return d.DetectRegistry(src)
	}

	return "", false, nil
}

// DetectRegistry tries to locate the download location for terraform registry source
func (d *RegistryDetector) DetectRegistry(src string) (string, bool, error) {
	parts := strings.Split(src, "?registryVersion=")
	if len(parts) < 2 {
		return "", false, fmt.Errorf("source not a valid registry path")
	}

	pkg := parts[0]
	version := parts[1]

	client := registry.NewClient(nil, d.httpClient)
	module, err := regsrc.ParseModuleSource(pkg)
	if err != nil {
		return "", false, err
	}

	location, err := client.ModuleLocation(module, version)
	if err != nil {
		return "", false, err
	}

	return location, true, nil
}
