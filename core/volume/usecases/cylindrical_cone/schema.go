package cylindrical_cone

import (
	"errors"

	"github.com/google/jsonschema-go/jsonschema"
)

var schema = &jsonschema.Schema{
	Type: "object",
	Properties: map[string]*jsonschema.Schema{
		"bigger_radius":      {Type: "number"},
		"incline_angle":      {Type: "number"},
		"cylindrical_height": {Type: "number"},
		"conical_height":     {Type: "number"},
	},
	Required: []string{"bigger_radius", "incline_angle", "cylindrical_height", "conical_height"},
}

func validate(dimensions map[string]any) (err error) {
	rs, err := schema.Resolve(nil)
	if err != nil {
		return errors.New("error on schema resolve: " + err.Error())
	}

	err = rs.Validate(dimensions)
	if err != nil {
		return errors.New("error on schema validation: " + err.Error())
	}

	return nil
}
