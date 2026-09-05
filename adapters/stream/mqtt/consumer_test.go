package mqtt

import (
	"context"
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/stretchr/testify/assert"
)

var _ mqtt.Message = (*fakeMessage)(nil)

// fakeMessage is the part of mqtt.Message the consumer reads.
type fakeMessage struct {
	topic    string
	payload  []byte
	retained bool
}

func (m *fakeMessage) Duplicate() bool   { return false }
func (m *fakeMessage) Qos() byte         { return 1 }
func (m *fakeMessage) Retained() bool    { return m.retained }
func (m *fakeMessage) Topic() string     { return m.topic }
func (m *fakeMessage) MessageID() uint16 { return 0 }
func (m *fakeMessage) Payload() []byte   { return m.payload }
func (m *fakeMessage) Ack()              {}

type recordingProcessor struct {
	processed chan []byte
}

func (p *recordingProcessor) Process(ctx context.Context, data []byte) error {
	p.processed <- data
	return nil
}

// A retained message is a message the broker replays on every reconnect.
// Processing it again would record a second pump run for a start that
// happened once, so the consumer must drop it.
func TestStreamConsumer_Consume_SkipsRetainedMessages(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	input := make(chan mqtt.Message, 2)
	processor := &recordingProcessor{processed: make(chan []byte, 2)}
	consumer := &streamConsumer{topic: "watertank/pump-01/pump", input: input}

	input <- &fakeMessage{topic: "watertank/pump-01/pump", payload: []byte(`{"state":"on"}`), retained: true}
	input <- &fakeMessage{topic: "watertank/pump-01/pump", payload: []byte(`{"state":"off"}`)}

	// Closing the input is what stops the consumer here: cancelling the
	// context would send it down the shutdown path, which unsubscribes
	// through a client this consumer was never given.
	close(input)

	go consumer.Consume(context.Background(), processor)

	select {
	case data := <-processor.processed:
		// The live message arrives; the retained one ahead of it does not.
		is.Equal(`{"state":"off"}`, string(data))
	case <-time.After(time.Second):
		t.Fatal("processor was never called")
	}

	is.Empty(processor.processed)
}

// The consumer stops when its input channel closes, without touching the
// client it was never given.
func TestStreamConsumer_Consume_StopsOnClosedChannel(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	input := make(chan mqtt.Message)
	processor := &recordingProcessor{processed: make(chan []byte, 1)}
	consumer := &streamConsumer{topic: "watertank/tank-01/level", input: input}

	done := make(chan struct{})
	go func() {
		consumer.Consume(context.Background(), processor)
		close(done)
	}()

	close(input)

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("consumer did not stop when the input channel closed")
	}

	is.Empty(processor.processed)
}
