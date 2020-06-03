package v012

import (
	"encoding/json"

	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"

	"github.com/avinor/tau/pkg/shell/processors"
)

// OutputProcessor processes output from terrafrom and parses them into variables.
// Implements the def.OutputProcessor interface
type OutputProcessor struct {
	processors.Buffer

	// decodeNames is set when the output names are encoded first. Internal attribute
	// only as output values should not be encoded normally
	decodeNames bool
}

// GetOutput takes the output from terraform command and parses the output into
// a map of string -> cty.Value
func (op *OutputProcessor) GetOutput() (map[string]cty.Value, error) {
	type OutputMeta struct {
		Sensitive bool            `json:"sensitive"`
		Type      json.RawMessage `json:"type"`
		Value     json.RawMessage `json:"value"`
	}
	outputs := map[string]OutputMeta{}
	values := map[string]cty.Value{}

	if err := json.Unmarshal([]byte(op.String()), &outputs); err != nil {
		return nil, err
	}

	for name, meta := range outputs {
		ctyType, err := ctyjson.UnmarshalType(meta.Type)
		if err != nil {
			return nil, err
		}

		ctyValue, err := ctyjson.Unmarshal(meta.Value, ctyType)
		if err != nil {
			return nil, err
		}

		if op.decodeNames {
			name = decodeName(name)
		}

		values[name] = ctyValue
	}

	return values, nil
}
