package adapters

import (
	"archimedes-worker/core/log"
	"archimedes-worker/core/processor"
	"context"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	waitTime = 250
)

type streamConsumer struct {
	topic  string
	client mqtt.Client
	input  <-chan mqtt.Message
}

func NewStreamConsumer(client mqtt.Client, topic string, inputChannel <-chan mqtt.Message) *streamConsumer {
	token := client.Subscribe(topic, 1, nil)
	token.Wait()

	err := token.Error()
	if err != nil {
		log.Log("Failed to subscribe to topic: " + topic + ": " + err.Error())
		return nil
	}
	log.Log("Subscribed to topic: " + topic)

	return &streamConsumer{
		client: client,
		input:  inputChannel,
		topic:  topic,
	}
}

func (consumer *streamConsumer) Consume(ctx context.Context, processTankStream processor.IProcessTankStream) {
	for {
		select {
		case msg, ok := <-consumer.input:
			if !ok {
				log.Log("Input channel closed, stopping consumer")
				return
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
