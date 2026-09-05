// Package domain defines the payloads carried over the MQTT streams that
// the processor use cases unmarshal and act on.
//
// Their shape is not this server's to choose: it is fixed by
// MQTT_CONTRACT.md, mirrored in every repository of the system, and the
// tank-node and pump-ctl firmware publish exactly these fields. Renaming
// one here does not translate a device's message, it stops the device
// being understood.
package domain

import "time"

// Event names carried in the envelope's event field. A message whose
// event does not match the topic it arrived on is not for the processor
// reading that topic.
const (
	EventLevel = "level"
	EventPump  = "pump"
)

// States a pump event reports. Unknown is the pump controller's last
// will: the broker sends it when the controller drops off uncleanly, so
// it means "the controller is unreachable", not "the pump stopped".
const (
	PumpStateOn      = "on"
	PumpStateOff     = "off"
	PumpStateUnknown = "unknown"
)

// Envelope is the set of fields every device event shares, whichever
// topic it arrives on.
type Envelope struct {
	// Event names the payload shape: one of the Event* constants.
	Event string `json:"event"`
	// Device is the id of the publisher, and the key this server stores
	// the event against: it must match the tank's or pump's id in the
	// database.
	Device string `json:"device"`
	// Timestamp is when the device observed the event, in UTC. It is
	// absent while a device's clock is unsynced and in last wills, which
	// the broker publishes on the device's behalf; see At.
	Timestamp *time.Time `json:"timestamp,omitempty"`
	// UptimeS is seconds since the publisher booted, absent in last
	// wills. Carried for diagnosis; nothing here acts on it.
	UptimeS *uint64 `json:"uptime_s,omitempty"`
}

// At returns the time the event should be recorded at: the device's own
// timestamp when it had a clock, and fallback otherwise. The boards have
// no battery-backed RTC, so an event published before SNTP lands carries
// no timestamp at all rather than a 1970 one, and the moment this server
// received it is the closest true answer available.
func (envelope Envelope) At(fallback time.Time) time.Time {
	if envelope.Timestamp == nil || envelope.Timestamp.IsZero() {
		return fallback
	}

	return envelope.Timestamp.UTC()
}

// WaterTankEvent is a tank-level sensor reading published to the water
// tank MQTT topic by tank-node.
type WaterTankEvent struct {
	Envelope
	// Valid reports whether FluidDistance is a reading at all.
	Valid bool `json:"valid"`
	// FluidDistance is the sensor-reported distance in centimetres from
	// the sensor face down to the fluid surface, consumed by
	// IVolumeCalculator.Calculate. It is null whenever Valid is false.
	// The tank's stored dimensions must be in centimetres too, since the
	// calculator works in whatever unit it is given.
	FluidDistance *float64 `json:"distance_cm"`
	// Reason explains an invalid reading: "sensor_unreadable" or
	// "node_offline".
	Reason string `json:"reason,omitempty"`
}

// Reading returns the event's distance and whether it can be used. A
// false second return means the node could not read its sensor, and the
// caller must record nothing: a missing distance is not a distance of
// zero, which would read as a tank filled to the sensor.
func (event WaterTankEvent) Reading() (float64, bool) {
	if !event.Valid || event.FluidDistance == nil {
		return 0, false
	}

	return *event.FluidDistance, true
}

// PumpEvent is a pump transition published to the pump status MQTT topic
// by pump-ctl. State is one of the PumpState* constants; Reason carries
// what caused the transition and is stored as the run's stop reason.
type PumpEvent struct {
	Envelope
	State  string `json:"state"`
	Reason string `json:"reason,omitempty"`
	// FlowLPM is the pipeline inflow at the transition, absent in the
	// controller's last will. Carried for diagnosis; nothing here acts
	// on it.
	FlowLPM *float64 `json:"flow_lpm,omitempty"`
	// FluidDistance is the last known tank level in centimetres, null
	// when the controller had no valid level and absent in its last
	// will. Carried for diagnosis; the stored volume comes from the tank
	// topic, not from here.
	FluidDistance *float64 `json:"distance_cm,omitempty"`
}
