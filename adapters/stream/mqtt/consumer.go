// Package mqtt adapts an MQTT topic subscription into the
// core/processor/interfaces.IProcessStream use cases, decoupling the
// tank and pump processors from the paho MQTT client.
package mqtt

import (
	"archimedes-server/core/log"
	"archimedes-server/core/processor/interfaces"
	"context"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	waitTime = 250
)

// streamConsumer subscribes to one MQTT topic and feeds every message it
// receives on inputChannel to an IProcessStream.
type streamConsumer struct {
	topic  string
	client mqtt.Client
	input  <-chan mqtt.Message
}

// NewStreamConsumer subscribes client to topic at QoS 1 and returns a
// consumer that reads inputChannel for Consume. inputChannel is expected to
// be fed by client's default publish handler; NewStreamConsumer does not
// wire that up itself. Returns nil if the subscription fails.
func NewStreamConsumer(client mqtt.Client, topic string, inputChannel <-chan mqtt.Message) *streamConsumer {
	token := client.Subscribe(topic, 1, nil)
	token.Wait()

	err := token.Error()
	if err != nil {
		log.Log(fmt.Sprintf("Failed to subscribe to topic %q: %q", topic, err.Error()))
		return nil
	}
	log.Log(fmt.Sprintf("Subscribed to topic: %q", topic))

	return &streamConsumer{
		client: client,
		input:  inputChannel,
		topic:  topic,
	}
}

// Consume blocks, dispatching each message from the input channel to
// processTankStream.Process (its return value is not inspected; failures
// must be logged inside Process) until the input channel is closed or ctx
// is cancelled, at which point it unsubscribes and disconnects the client.
//
// Retained messages are dropped rather than processed. The pump topic is
// retained so a dashboard connecting late can see what the pump is doing,
// but every processor behind this consumer turns a message into a record
// of something that happened: the client connects with a clean session,
// so the broker redelivers that retained message on every reconnect, and
// processing it again would insert a second pump run for a start that
// happened once. The cost is that a server starting mid-run does not
// learn the pump is already running until its next transition.
func (consumer *streamConsumer) Consume(ctx context.Context, processTankStream interfaces.IProcessStream) {
	for {
		select {
		case msg, ok := <-consumer.input:
			if !ok {
				log.Log("Input channel closed, stopping consumer")
				return
			}
			if msg.Retained() {
				log.Log(fmt.Sprintf("Ignoring retained message on topic %q", msg.Topic()))
				continue
			}
			processTankStream.Process(ctx, msg.Payload())

		case <-ctx.Done():
			log.Log("Unsubscribing and disconnecting...")
			consumer.client.Unsubscribe(consumer.topic)
			consumer.client.Disconnect(waitTime)
			return
		}
	}
}
