package domain

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWaterTankEvent_Unmarshal(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Exactly what tank-node publishes; see MQTT_CONTRACT.md.
	data := []byte(`{"event":"level","device":"tank-01",` +
		`"timestamp":"2026-09-05T03:10:12Z","valid":true,` +
		`"distance_cm":62.5,"uptime_s":360}`)

	var event WaterTankEvent
	is.NoError(json.Unmarshal(data, &event))

	is.Equal(EventLevel, event.Event)
	is.Equal("tank-01", event.Device)
	is.Equal(time.Date(2026, 9, 5, 3, 10, 12, 0, time.UTC), *event.Timestamp)
	is.Equal(uint64(360), *event.UptimeS)

	distance, ok := event.Reading()
	is.True(ok)
	is.Equal(62.5, distance)
}

func TestWaterTankEvent_Reading(t *testing.T) {
	t.Parallel()

	distance := 62.5

	tests := []struct {
		name         string
		event        WaterTankEvent
		wantDistance float64
		wantOK       bool
	}{
		{
			name:         "valid reading is usable",
			event:        WaterTankEvent{Valid: true, FluidDistance: &distance},
			wantDistance: distance,
			wantOK:       true,
		},
		{
			name:  "unreadable sensor yields nothing, not zero centimetres",
			event: WaterTankEvent{Valid: false, Reason: "sensor_unreadable"},
		},
		{
			name:  "node offline last will yields nothing",
			event: WaterTankEvent{Valid: false, Reason: "node_offline"},
		},
		{
			name:  "valid flag with a null distance yields nothing",
			event: WaterTankEvent{Valid: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			is := assert.New(t)

			distance, ok := tt.event.Reading()

			is.Equal(tt.wantOK, ok)
			is.Equal(tt.wantDistance, distance)
		})
	}
}

func TestPumpEvent_Unmarshal(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Exactly what pump-ctl publishes; see MQTT_CONTRACT.md.
	data := []byte(`{"event":"pump","device":"pump-01",` +
		`"timestamp":"2026-09-05T03:10:12Z","state":"on",` +
		`"reason":"flow_confirmed","flow_lpm":11.40,` +
		`"tank_state":"not_full","uptime_s":338}`)

	var event PumpEvent
	is.NoError(json.Unmarshal(data, &event))

	is.Equal(EventPump, event.Event)
	is.Equal("pump-01", event.Device)
	is.Equal(PumpStateOn, event.State)
	is.Equal("flow_confirmed", event.Reason)
	is.Equal(11.40, *event.FlowLPM)
	is.Equal("not_full", event.TankState)
}

func TestPumpEvent_UnmarshalLastWill(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// The broker publishes this on the controller's behalf, so it has
	// neither a timestamp nor an uptime.
	data := []byte(`{"event":"pump","device":"pump-01","state":"unknown",` +
		`"reason":"controller_offline"}`)

	var event PumpEvent
	is.NoError(json.Unmarshal(data, &event))

	is.Equal(PumpStateUnknown, event.State)
	is.Nil(event.Timestamp)
	is.Nil(event.UptimeS)
}

func TestEnvelope_At(t *testing.T) {
	t.Parallel()

	deviceTime := time.Date(2026, 9, 5, 3, 10, 12, 0, time.UTC)
	receivedAt := time.Date(2026, 9, 5, 3, 10, 15, 0, time.UTC)
	zero := time.Time{}
	offset := deviceTime.In(time.FixedZone("UTC-3", -3*60*60))

	tests := []struct {
		name     string
		envelope Envelope
		want     time.Time
	}{
		{
			name:     "device timestamp wins",
			envelope: Envelope{Timestamp: &deviceTime},
			want:     deviceTime,
		},
		{
			name:     "missing timestamp falls back to receipt time",
			envelope: Envelope{},
			want:     receivedAt,
		},
		{
			name:     "zero timestamp falls back to receipt time",
			envelope: Envelope{Timestamp: &zero},
			want:     receivedAt,
		},
		{
			name:     "offset timestamp is normalised to UTC",
			envelope: Envelope{Timestamp: &offset},
			want:     deviceTime,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			is := assert.New(t)

			got := tt.envelope.At(receivedAt)

			is.True(got.Equal(tt.want), "got %s, want %s", got, tt.want)
			is.Equal(time.UTC, got.Location())
		})
	}
}
