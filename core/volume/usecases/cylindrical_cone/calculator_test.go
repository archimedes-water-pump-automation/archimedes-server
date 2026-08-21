package cylindrical_cone

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func validDimensions() map[string]any {
	return map[string]any{
		"bigger_radius":      2.0,
		"incline_angle":      15.0,
		"cylindrical_height": 5.0,
		"conical_height":     3.0,
	}
}

func TestNewCylindricalPartialConeShape(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	dimensions := validDimensions()
	shape := NewCylindricalPartialConeShape(dimensions)

	is.NotNil(shape)
	is.Equal(dimensions, shape.dimensions)
}

func TestCylindricalPartialConeShape_Calculate(t *testing.T) {
	tests := []struct {
		name          string
		dimensions    map[string]any
		fluidDistance float64
		wantVolume    float64
		wantErr       bool
	}{
		{
			name:          "fluid distance zero fills cylinder and full cone",
			dimensions:    validDimensions(),
			fluidDistance: 0,
			wantVolume:    87.40880089315004,
		},
		{
			name:          "fluid distance within conical section, tapered",
			dimensions:    validDimensions(),
			fluidDistance: 1,
			wantVolume:    81.83177979621324,
		},
		{
			name:          "fluid distance further into conical section",
			dimensions:    validDimensions(),
			fluidDistance: 2,
			wantVolume:    73.78983465864047,
		},
		{
			name:          "fluid distance at conical/cylindrical boundary",
			dimensions:    validDimensions(),
			fluidDistance: 3,
			wantVolume:    62.83185307179586,
		},
		{
			name:          "fluid distance equal to cylindrical height",
			dimensions:    validDimensions(),
			fluidDistance: 5,
			wantVolume:    37.69911184307752,
		},
		{
			name:          "fluid distance beyond cylindrical height",
			dimensions:    validDimensions(),
			fluidDistance: 6,
			wantVolume:    25.132741228718345,
		},
		{
			name:          "fluid distance at tank total height empties tank",
			dimensions:    validDimensions(),
			fluidDistance: 8,
			wantVolume:    0,
		},
		{
			name:          "fluid distance far beyond tank height clamps to zero",
			dimensions:    validDimensions(),
			fluidDistance: 100,
			wantVolume:    0,
		},
		{
			name: "invalid dimensions returns error and zero volume",
			dimensions: map[string]any{
				"bigger_radius": 2.0,
			},
			fluidDistance: 1,
			wantVolume:    0,
			wantErr:       true,
		},
		{
			name:          "empty dimensions returns error and zero volume",
			dimensions:    map[string]any{},
			fluidDistance: 1,
			wantVolume:    0,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			is := assert.New(t)

			shape := NewCylindricalPartialConeShape(tt.dimensions)
			volume, err := shape.Calculate(context.Background(), tt.fluidDistance)

			if tt.wantErr {
				is.Error(err)
			} else {
				is.NoError(err)
			}
			is.InDelta(tt.wantVolume, volume, 1e-9)
		})
	}
}
