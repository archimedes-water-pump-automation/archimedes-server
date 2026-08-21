package tank

import (
	tankInterfaces "archimedes-server/core/tank/interfaces"
	volumeInterfaces "archimedes-server/core/volume/interfaces"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	_ tankInterfaces.IUpdateTank         = (*fakeUpdateTank)(nil)
	_ volumeInterfaces.IGetVolumeType    = (*fakeGetVolumeType)(nil)
	_ volumeInterfaces.IVolumeCalculator = (*fakeVolumeCalculator)(nil)
)

type updateVolumeCall struct {
	tankID    string
	volume    float64
	updatedAt time.Time
}

type fakeUpdateTank struct {
	calls []updateVolumeCall
	err   error
}

func (f *fakeUpdateTank) UpdateVolume(ctx context.Context, tankID string, newVolume float64, updatedAt time.Time) error {
	f.calls = append(f.calls, updateVolumeCall{tankID: tankID, volume: newVolume, updatedAt: updatedAt})
	return f.err
}

type fakeGetVolumeType struct {
	calculator volumeInterfaces.IVolumeCalculator
	err        error
}

func (f *fakeGetVolumeType) GetVolumeFromShape(ctx context.Context, tankID string) (volumeInterfaces.IVolumeCalculator, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.calculator, nil
}

type fakeVolumeCalculator struct {
	volume float64
	err    error
}

func (f *fakeVolumeCalculator) Calculate(ctx context.Context, fluidDistance float64) (float64, error) {
	return f.volume, f.err
}

func TestNewProcessTankUpdate(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	repo := &fakeUpdateTank{}
	getVolumeType := &fakeGetVolumeType{}
	processor := NewProcessTankUpdate(repo, getVolumeType)

	is.NotNil(processor)
	is.Equal(repo, processor.repository)
	is.Equal(getVolumeType, processor.getVolumeType)
}

func TestProcessTankUpdate_Process(t *testing.T) {
	timestamp := time.Date(2026, 8, 21, 12, 0, 0, 0, time.UTC)
	validData := []byte(`{"tank_id":"tank-1","event_type":"reading","distance":1.5,"timestamp":"` + timestamp.Format(time.RFC3339Nano) + `"}`)

	tests := []struct {
		name           string
		data           []byte
		calculator     *fakeVolumeCalculator
		getVolumeErr   error
		updateErr      error
		wantErr        bool
		wantUpdateCall *updateVolumeCall
	}{
		{
			name:           "valid event updates volume with calculated value",
			data:           validData,
			calculator:     &fakeVolumeCalculator{volume: 42.5},
			wantUpdateCall: &updateVolumeCall{tankID: "tank-1", volume: 42.5, updatedAt: timestamp},
		},
		{
			name:    "invalid json returns error",
			data:    []byte(`not-json`),
			wantErr: true,
		},
		{
			name:         "error getting volume calculator propagates",
			data:         validData,
			getVolumeErr: errors.New("unknown shape"),
			wantErr:      true,
		},
		{
			name:       "error from calculator propagates",
			data:       validData,
			calculator: &fakeVolumeCalculator{err: errors.New("bad dimensions")},
			wantErr:    true,
		},
		{
			name:           "error from repository propagates",
			data:           validData,
			calculator:     &fakeVolumeCalculator{volume: 10.0},
			updateErr:      errors.New("db down"),
			wantErr:        true,
			wantUpdateCall: &updateVolumeCall{tankID: "tank-1", volume: 10.0, updatedAt: timestamp},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			is := assert.New(t)

			repo := &fakeUpdateTank{err: tt.updateErr}
			getVolumeType := &fakeGetVolumeType{calculator: tt.calculator, err: tt.getVolumeErr}
			processor := NewProcessTankUpdate(repo, getVolumeType)

			err := processor.Process(context.Background(), tt.data)

			if tt.wantErr {
				is.Error(err)
			} else {
				is.NoError(err)
			}

			if tt.wantUpdateCall != nil {
				is.Equal([]updateVolumeCall{*tt.wantUpdateCall}, repo.calls)
			} else {
				is.Empty(repo.calls)
			}
		})
	}
}
