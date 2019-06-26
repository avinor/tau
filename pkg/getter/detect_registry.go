package getter

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform/registry"
	"github.com/hashicorp/terraform/registry/regsrc"
)

// RegistryDetector implements detector to detect terraform registry
// Src have to be formatted with query parameter ?registryVersion=
// avinor/storage-account/azurerm?registryVersion=1.0
type RegistryDetector struct {
	httpClient *http.Client
}

func (d *RegistryDetector) Detect(src, _ string) (string, bool, error) {
	if len(src) == 0 {
		return "", false, nil
	}

	if strings.Contains(src, "?registryVersion=") {
		return d.DetectRegistry(src)
	}

	return "", false, nil
}

func (d *RegistryDetector) DetectRegistry(src string) (string, bool, error) {
	parts := strings.Split(src, "?registryVersion=")
	if len(parts) < 2 {
		return "", false, fmt.Errorf("Source not a valid registry path")
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
