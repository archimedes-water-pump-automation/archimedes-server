package tank

import (
	tankInterfaces "archimedes-server/core/tank/interfaces"
	volumeInterfaces "archimedes-server/core/volume/interfaces"
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// receivedAt is what the processor stamps an event with when the device
// published none of its own. Pinned before any test runs so the fallback
// path is assertable and races cannot see a half-written clock.
var receivedAt = time.Date(2026, 8, 21, 12, 30, 0, 0, time.UTC)

func TestMain(m *testing.M) {
	timeNow = func() time.Time { return receivedAt }
	os.Exit(m.Run())
}

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
	// gotDistance is the distance the processor passed in, so a test can
	// tell "never called" from "called with a zero distance".
	gotDistance float64
}

func (f *fakeVolumeCalculator) Calculate(ctx context.Context, fluidDistance float64) (float64, error) {
	f.gotDistance = fluidDistance
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
	validData := []byte(`{"event":"level","device":"tank-1","timestamp":"` +
		timestamp.Format(time.RFC3339Nano) +
		`","valid":true,"distance_cm":62.5,"uptime_s":360}`)

	tests := []struct {
		name           string
		data           []byte
		calculator     *fakeVolumeCalculator
		getVolumeErr   error
		updateErr      error
		wantErr        bool
		wantDistance   float64
		wantUpdateCall *updateVolumeCall
	}{
		{
			name:           "valid event updates volume with calculated value",
			data:           validData,
			calculator:     &fakeVolumeCalculator{volume: 42.5},
			wantDistance:   62.5,
			wantUpdateCall: &updateVolumeCall{tankID: "tank-1", volume: 42.5, updatedAt: timestamp},
		},
		{
			name: "event without a timestamp is stamped on receipt",
			data: []byte(`{"event":"level","device":"tank-1","valid":true,` +
				`"distance_cm":62.5,"uptime_s":360}`),
			calculator:     &fakeVolumeCalculator{volume: 42.5},
			wantDistance:   62.5,
			wantUpdateCall: &updateVolumeCall{tankID: "tank-1", volume: 42.5, updatedAt: receivedAt},
		},
		{
			name: "invalid reading stores nothing",
			data: []byte(`{"event":"level","device":"tank-1","timestamp":"` +
				timestamp.Format(time.RFC3339Nano) +
				`","valid":false,"distance_cm":null,"reason":"sensor_unreadable"}`),
			calculator: &fakeVolumeCalculator{volume: 42.5},
		},
		{
			name: "node offline last will stores nothing",
			data: []byte(`{"event":"level","device":"tank-1","valid":false,` +
				`"distance_cm":null,"reason":"node_offline"}`),
			calculator: &fakeVolumeCalculator{volume: 42.5},
		},
		{
			name: "reading flagged valid with a null distance stores nothing",
			data: []byte(`{"event":"level","device":"tank-1","valid":true,` +
				`"distance_cm":null}`),
			calculator: &fakeVolumeCalculator{volume: 42.5},
		},
		{
			name:       "event of another type is a no-op",
			data:       []byte(`{"event":"pump","device":"pump-1","state":"on"}`),
			calculator: &fakeVolumeCalculator{volume: 42.5},
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
			name:         "error from calculator propagates",
			data:         validData,
			calculator:   &fakeVolumeCalculator{err: errors.New("bad dimensions")},
			wantDistance: 62.5,
			wantErr:      true,
		},
		{
			name:           "error from repository propagates",
			data:           validData,
			calculator:     &fakeVolumeCalculator{volume: 10.0},
			updateErr:      errors.New("db down"),
			wantErr:        true,
			wantDistance:   62.5,
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

			if tt.calculator != nil {
				is.Equal(tt.wantDistance, tt.calculator.gotDistance)
			}
		})
	}
}
