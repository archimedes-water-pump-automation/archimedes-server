package main

import (
	"context"
	"crypto/tls"
	"errors"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"archimedes-server/adapters"
	"archimedes-server/cmd/webserver"
	"archimedes-server/core/log"
	"archimedes-server/core/processor"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	server   = "tcp://localhost:1883"
	clientID = "archimedes-server-subscriber"
)

var (
	topic = os.Getenv("TOPIC")
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	var sigChan = make(chan os.Signal, 1)
	var mqttMsgChan = make(chan mqtt.Message)

	log.SetLogger(adapters.NewLogger())

	client := NewClient(mqttMsgChan)

	pool, err := NewPool(context.Background(), os.Getenv("DB_CONN_STRING"))
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	updateRepository := adapters.NewUpdateTankRepository(pool)
	readRepository := adapters.NewReadTankRepository(pool)
	processor := processor.NewProcessTankStreamUseCase(updateRepository)

	consumer := adapters.NewStreamConsumer(client, topic, mqttMsgChan)
	if consumer == nil {
		panic(errors.New("Failed to create stream consumer"))
	}

	ctx, cancel := context.WithCancel(context.Background())

	go webserver.Serve(readRepository)

	go func() {
		defer wg.Done()
		consumer.Consume(ctx, processor)
	}()

	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	cancel()

	wg.Wait()
	log.Log("Archimedes terminated, exiting...")
}

func NewClient(inputChannel chan mqtt.Message) mqtt.Client {

	opts := mqtt.NewClientOptions()
	opts.AddBroker(server)
	opts.SetClientID(clientID)
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		inputChannel <- msg
	})
	opts.OnConnect = func(client mqtt.Client) {
		log.Log("Connected to MQTT Broker")
	}
	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		log.Log("Connection lost: " + err.Error())
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return client
}

func NewPool(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	if enabled, _ := strconv.ParseBool(os.Getenv("DB_TLS_ENABLED")); enabled {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
		}
		config.ConnConfig.TLSConfig = tlsConfig
	}

	config.ConnConfig.RuntimeParams["timezone"] = "UTC"

	config.MaxConns = 3
	config.MinConns = 1
	config.MaxConnLifetime = 30 * time.Minute
	config.MaxConnIdleTime = 10 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
