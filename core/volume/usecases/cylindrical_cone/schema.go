package cylindrical_cone

import (
	"errors"
	"sync"

	"github.com/google/jsonschema-go/jsonschema"
)

var (
	rs     *jsonschema.Resolved
	rsOnce sync.Once
	rsErr  error
)

func float64Ptr(f float64) *float64 {
	return &f
}

// schema describes the shape of the dimensions map a cylindrical-cone tank
// must store: four required numeric fields, matching the fields read by
// Calculate.
var schema = &jsonschema.Schema{
	Type: "object",
	Properties: map[string]*jsonschema.Schema{
		"bigger_radius":      {Type: "number"},
		"incline_angle":      {Type: "number", Minimum: float64Ptr(0.0), Maximum: float64Ptr(90.0)},
		"cylindrical_height": {Type: "number"},
		"conical_height":     {Type: "number"},
	},
	Required: []string{"bigger_radius", "incline_angle", "cylindrical_height", "conical_height"},
}

// resolveSchema resolves schema once and caches the result, since Resolve
// does the same static work on every call.
func resolveSchema() (*jsonschema.Resolved, error) {
	rsOnce.Do(func() {
		rs, rsErr = schema.Resolve(nil)
	})

	return rs, rsErr
}

// validate checks dimensions against schema, catching malformed or
// incomplete tank configuration before Calculate attempts unchecked type
// assertions on the map's values.
func validate(dimensions map[string]any) (err error) {
	rs, err := resolveSchema()
	if err != nil {
		return errors.New("error on schema resolve: " + err.Error())
	}

	err = rs.Validate(dimensions)
	if err != nil {
		return errors.New("error on schema validation: " + err.Error())
	}

	return nil
}
