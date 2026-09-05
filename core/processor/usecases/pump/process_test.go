package pump

import (
	"archimedes-server/core/pump/interfaces"
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

var _ interfaces.IUpdatePumpStatus = (*fakeUpdatePumpStatus)(nil)

type fakeUpdatePumpStatus struct {
	startCalls []startCall
	stopCalls  []stopCall
	startErr   error
	stopErr    error
}

type startCall struct {
	pumpID    string
	timestamp time.Time
}

type stopCall struct {
	pumpID     string
	timestamp  time.Time
	stopReason string
}

func (f *fakeUpdatePumpStatus) StartPump(ctx context.Context, pumpID string, timestamp time.Time) error {
	f.startCalls = append(f.startCalls, startCall{pumpID: pumpID, timestamp: timestamp})
	return f.startErr
}

func (f *fakeUpdatePumpStatus) StopPump(ctx context.Context, pumpID string, timestamp time.Time, stopReason string) error {
	f.stopCalls = append(f.stopCalls, stopCall{pumpID: pumpID, timestamp: timestamp, stopReason: stopReason})
	return f.stopErr
}

func TestNewProcessPumpStatusUpdate(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	repo := &fakeUpdatePumpStatus{}
	processor := NewProcessPumpStatusUpdate(repo)

	is.NotNil(processor)
	is.Equal(repo, processor.repository)
}

func TestProcessPumpUpdate_Process(t *testing.T) {
	timestamp := time.Date(2026, 8, 21, 12, 0, 0, 0, time.UTC)
	stamp := timestamp.Format(time.RFC3339Nano)

	tests := []struct {
		name      string
		data      []byte
		repoErr   error
		wantErr   bool
		wantStart []startCall
		wantStop  []stopCall
	}{
		{
			name: "pump on starts a run",
			data: []byte(`{"event":"pump","device":"pump-1","timestamp":"` + stamp +
				`","state":"on","reason":"flow_confirmed","flow_lpm":11.4,` +
				`"tank_state":"not_full","uptime_s":338}`),
			wantStart: []startCall{{pumpID: "pump-1", timestamp: timestamp}},
		},
		{
			name: "pump off stops the run with its reason",
			data: []byte(`{"event":"pump","device":"pump-2","timestamp":"` + stamp +
				`","state":"off","reason":"tank_full","flow_lpm":0.0,` +
				`"tank_state":"full","uptime_s":607}`),
			wantStop: []stopCall{{pumpID: "pump-2", timestamp: timestamp, stopReason: "tank_full"}},
		},
		{
			name: "event without a timestamp is stamped on receipt",
			data: []byte(`{"event":"pump","device":"pump-1","state":"on",` +
				`"reason":"flow_confirmed","flow_lpm":11.4,"uptime_s":338}`),
			wantStart: []startCall{{pumpID: "pump-1", timestamp: receivedAt}},
		},
		{
			name: "unknown state from the controller last will is a no-op",
			data: []byte(`{"event":"pump","device":"pump-3","state":"unknown",` +
				`"reason":"controller_offline"}`),
		},
		{
			name: "event of another type is a no-op",
			data: []byte(`{"event":"level","device":"tank-1","valid":true,"distance_cm":62.5}`),
		},
		{
			name:    "invalid json returns error",
			data:    []byte(`not-json`),
			wantErr: true,
		},
		{
			name: "repository error on start propagates",
			data: []byte(`{"event":"pump","device":"pump-1","timestamp":"` + stamp +
				`","state":"on","reason":"flow_confirmed"}`),
			repoErr:   errors.New("db down"),
			wantErr:   true,
			wantStart: []startCall{{pumpID: "pump-1", timestamp: timestamp}},
		},
		{
			name: "repository error on stop propagates",
			data: []byte(`{"event":"pump","device":"pump-2","timestamp":"` + stamp +
				`","state":"off","reason":"pipeline_dry"}`),
			repoErr:  errors.New("db down"),
			wantErr:  true,
			wantStop: []stopCall{{pumpID: "pump-2", timestamp: timestamp, stopReason: "pipeline_dry"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			is := assert.New(t)

			repo := &fakeUpdatePumpStatus{startErr: tt.repoErr, stopErr: tt.repoErr}
			processor := NewProcessPumpStatusUpdate(repo)

			err := processor.Process(context.Background(), tt.data)

			if tt.wantErr {
				is.Error(err)
			} else {
				is.NoError(err)
			}

			is.Equal(tt.wantStart, repo.startCalls)
			is.Equal(tt.wantStop, repo.stopCalls)
		})
	}
}
