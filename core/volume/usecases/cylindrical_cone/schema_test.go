package cylindrical_cone

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name       string
		dimensions map[string]any
		wantErr    bool
	}{
		{
			name: "all required fields present with correct types",
			dimensions: map[string]any{
				"bigger_radius":      2.0,
				"incline_angle":      15.0,
				"cylindrical_height": 5.0,
				"conical_height":     3.0,
			},
			wantErr: false,
		},
		{
			name:       "empty dimensions",
			dimensions: map[string]any{},
			wantErr:    true,
		},
		{
			name: "missing bigger_radius",
			dimensions: map[string]any{
				"incline_angle":      15.0,
				"cylindrical_height": 5.0,
				"conical_height":     3.0,
			},
			wantErr: true,
		},
		{
			name: "missing conical_height",
			dimensions: map[string]any{
				"bigger_radius":      2.0,
				"incline_angle":      15.0,
				"cylindrical_height": 5.0,
			},
			wantErr: true,
		},
		{
			name: "missing incline_angle",
			dimensions: map[string]any{
				"bigger_radius":      2.0,
				"conical_height":     3.0,
				"cylindrical_height": 5.0,
			},
			wantErr: true,
		},
		{
			name: "missing cylindrical_height",
			dimensions: map[string]any{
				"bigger_radius":  2.0,
				"incline_angle":  15.0,
				"conical_height": 3.0,
			},
			wantErr: true,
		},
		{
			name: "wrong type for bigger_radius",
			dimensions: map[string]any{
				"bigger_radius":      "not-a-number",
				"incline_angle":      15.0,
				"cylindrical_height": 5.0,
				"conical_height":     3.0,
			},
			wantErr: true,
		},
		{
			name:       "nil dimensions",
			dimensions: nil,
			wantErr:    true,
		},
		{
			name: "wrong type for incline_angle",
			dimensions: map[string]any{
				"bigger_radius":      2.0,
				"incline_angle":      "not-a-number",
				"cylindrical_height": 5.0,
				"conical_height":     3.0,
			},
			wantErr: true,
		},
		{
			name: "wrong type for cylindrical_height",
			dimensions: map[string]any{
				"bigger_radius":      2.0,
				"incline_angle":      15.0,
				"cylindrical_height": "not-a-number",
				"conical_height":     3.0,
			},
			wantErr: true,
		},
		{
			name: "wrong type for conical_height",
			dimensions: map[string]any{
				"bigger_radius":      2.0,
				"incline_angle":      15.0,
				"cylindrical_height": 5.0,
				"conical_height":     "not-a-number",
			},
			wantErr: true,
		},
		{
			name: "incline_angle at lower bound",
			dimensions: map[string]any{
				"bigger_radius":      2.0,
				"incline_angle":      0.0,
				"cylindrical_height": 5.0,
				"conical_height":     3.0,
			},
			wantErr: false,
		},
		{
			name: "incline_angle at upper bound",
			dimensions: map[string]any{
				"bigger_radius":      2.0,
				"incline_angle":      90.0,
				"cylindrical_height": 5.0,
				"conical_height":     3.0,
			},
			wantErr: false,
		},
		{
			name: "incline_angle below lower bound",
			dimensions: map[string]any{
				"bigger_radius":      2.0,
				"incline_angle":      -0.1,
				"cylindrical_height": 5.0,
				"conical_height":     3.0,
			},
			wantErr: true,
		},
		{
			name: "incline_angle above upper bound",
			dimensions: map[string]any{
				"bigger_radius":      2.0,
				"incline_angle":      90.1,
				"cylindrical_height": 5.0,
				"conical_height":     3.0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			is := assert.New(t)

			err := validate(tt.dimensions)

			if tt.wantErr {
				is.Error(err)
			} else {
				is.NoError(err)
			}
		})
	}
}
