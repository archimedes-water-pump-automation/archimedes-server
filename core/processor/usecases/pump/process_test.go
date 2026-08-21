package pump

import (
	"archimedes-server/core/pump/interfaces"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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

	tests := []struct {
		name      string
		data      []byte
		repoErr   error
		wantErr   bool
		wantStart []startCall
		wantStop  []stopCall
	}{
		{
			name:      "start event calls StartPump",
			data:      []byte(`{"pump_id":"pump-1","event_type":"start","timestamp":"` + timestamp.Format(time.RFC3339Nano) + `"}`),
			wantStart: []startCall{{pumpID: "pump-1", timestamp: timestamp}},
		},
		{
			name:     "stop event calls StopPump with reason",
			data:     []byte(`{"pump_id":"pump-2","event_type":"stop","timestamp":"` + timestamp.Format(time.RFC3339Nano) + `","stop_reason":"manual"}`),
			wantStop: []stopCall{{pumpID: "pump-2", timestamp: timestamp, stopReason: "manual"}},
		},
		{
			name:    "unknown event type is a no-op",
			data:    []byte(`{"pump_id":"pump-3","event_type":"unknown"}`),
			wantErr: false,
		},
		{
			name:    "invalid json returns error",
			data:    []byte(`not-json`),
			wantErr: true,
		},
		{
			name:      "repository error on start propagates",
			data:      []byte(`{"pump_id":"pump-1","event_type":"start","timestamp":"` + timestamp.Format(time.RFC3339Nano) + `"}`),
			repoErr:   errors.New("db down"),
			wantErr:   true,
			wantStart: []startCall{{pumpID: "pump-1", timestamp: timestamp}},
		},
		{
			name:     "repository error on stop propagates",
			data:     []byte(`{"pump_id":"pump-2","event_type":"stop","timestamp":"` + timestamp.Format(time.RFC3339Nano) + `"}`),
			repoErr:  errors.New("db down"),
			wantErr:  true,
			wantStop: []stopCall{{pumpID: "pump-2", timestamp: timestamp, stopReason: ""}},
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
