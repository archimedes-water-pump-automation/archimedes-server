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
	server       = "tcp://localhost:1883"
	clientIDTank = "archimedes-server-tank-subscriber"
	clientIDPump = "archimedes-server-pump-subscriber"
)

var (
	waterTankTopic  = os.Getenv("WATER_TANK_TOPIC")
	pumpStatusTopic = os.Getenv("PUMP_STATUS_TOPIC")
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	var sigChan = make(chan os.Signal, 1)
	var tankStreamChannel = make(chan mqtt.Message)
	var pumpStatusChannel = make(chan mqtt.Message)

	log.SetLogger(adapters.NewLogger())

	tankStreamClient := NewClient(tankStreamChannel, clientIDTank)
	pumpStatusClient := NewClient(pumpStatusChannel, clientIDPump)

	pool, err := NewPool(context.Background(), os.Getenv("DB_CONN_STRING"))
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	updateTankProcessor := processor.NewProcessTankUpdate(adapters.NewUpdateTankRepository(pool))
	updatePumpProcessor := processor.NewProcessPumpStatusUpdate(adapters.NewUpdatePumpStatusRepository(pool))

	waterTankConsumer := adapters.NewStreamConsumer(tankStreamClient, waterTankTopic, tankStreamChannel)
	if waterTankConsumer == nil {
		panic(errors.New("Failed to create stream water tank consumer"))
	}

	pumpStatusConsumer := adapters.NewStreamConsumer(pumpStatusClient, pumpStatusTopic, pumpStatusChannel)
	if pumpStatusConsumer == nil {
		panic(errors.New("Failed to create stream pump status consumer"))
	}

	ctx, cancel := context.WithCancel(context.Background())

	go webserver.Serve(adapters.NewReadTankRepository(pool), adapters.NewReadPumpStatusRepository(pool))

	go func() {
		defer wg.Done()
		waterTankConsumer.Consume(ctx, updateTankProcessor)
	}()

	go func() {
		defer wg.Done()
		pumpStatusConsumer.Consume(ctx, updatePumpProcessor)
	}()

	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	cancel()

	wg.Wait()
	log.Log("Archimedes terminated, exiting...")
}

func NewClient(inputChannel chan mqtt.Message, clientID string) mqtt.Client {

	opts := mqtt.NewClientOptions()
	opts.AddBroker(server)
	opts.SetClientID(clientID)
	opts.SetDefaultPublishHandler(func(tankStreamClient mqtt.Client, msg mqtt.Message) {
		inputChannel <- msg
	})
	opts.OnConnect = func(tankStreamClient mqtt.Client) {
		log.Log("Connected to MQTT Broker")
	}
	opts.OnConnectionLost = func(tankStreamClient mqtt.Client, err error) {
		log.Log("Connection lost: " + err.Error())
	}

	tankStreamClient := mqtt.NewClient(opts)
	if token := tankStreamClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return tankStreamClient
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
