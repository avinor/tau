package v012

import (
	"github.com/zclconf/go-cty/cty"
)

type Processor struct {
}

func (p *Processor) ProcessOutput(output []byte) (map[string]cty.Value, error) {
	return nil, nil
}
