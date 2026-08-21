// Package interfaces declares the contract stream consumers use to hand off
// a raw message body for processing, decoupling adapters/stream/mqtt from
// the tank and pump processor implementations.
package interfaces

import "context"

// IProcessStream handles one raw message from a subscribed MQTT topic.
// Implemented by core/processor/usecases/tank and
// core/processor/usecases/pump.
type IProcessStream interface {
	// Process unmarshals data and applies its effect (e.g. updating a
	// tank's volume or a pump's run status). The returned error is not
	// currently inspected by stream consumers, so implementations must log
	// failures themselves before returning.
	Process(ctx context.Context, data []byte) error
}
