package main

import (
	"archimedes-server/adapters/database/postgresql"
	"archimedes-server/adapters/endpoint/http"
	"archimedes-server/adapters/log/file"
	mqtt_adapter "archimedes-server/adapters/stream/mqtt"
	"archimedes-server/core/log"
	"archimedes-server/core/processor/usecases/pump"
	"archimedes-server/core/processor/usecases/tank"
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

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	clientIDTank = "archimedes-server-tank-subscriber"
	clientIDPump = "archimedes-server-pump-subscriber"
)

var (
	waterTankTopic  = os.Getenv("WATER_TANK_TOPIC")
	pumpStatusTopic = os.Getenv("PUMP_STATUS_TOPIC")
	server          = os.Getenv("MQTT_BROKER_URL")
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	var sigChan = make(chan os.Signal, 1)
	var tankStreamChannel = make(chan mqtt.Message)
	var pumpStatusChannel = make(chan mqtt.Message)

	logFile, err := os.OpenFile(os.Getenv("LOG_FILE"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer func() {
		log.Log("Closing log file...")
		logFile.Close()
	}()

	log.SetLogger(file.NewLogger(logFile))

	tankStreamClient := NewClient(tankStreamChannel, clientIDTank)
	pumpStatusClient := NewClient(pumpStatusChannel, clientIDPump)

	pool, err := NewPool(context.Background(), os.Getenv("DB_CONN_STRING"))
	if err != nil {
		panic(err)
	}
	defer func() {
		log.Log("Closing database connection pool...")
		pool.Close()
	}()

	readTankRepository := postgresql.NewReadTankRepository(pool)

	updateTankProcessor := tank.NewProcessTankUpdate(postgresql.NewUpdateTankRepository(pool), postgresql.NewGetVolumeType(pool))
	updatePumpProcessor := pump.NewProcessPumpStatusUpdate(postgresql.NewUpdatePumpStatusRepository(pool))

	waterTankConsumer := mqtt_adapter.NewStreamConsumer(tankStreamClient, waterTankTopic, tankStreamChannel)
	if waterTankConsumer == nil {
		panic(errors.New("Failed to create stream water tank consumer"))
	}

	pumpStatusConsumer := mqtt_adapter.NewStreamConsumer(pumpStatusClient, pumpStatusTopic, pumpStatusChannel)
	if pumpStatusConsumer == nil {
		panic(errors.New("Failed to create stream pump status consumer"))
	}

	ctx, cancel := context.WithCancel(context.Background())

	server := http.Serve("8080", readTankRepository, postgresql.NewReadPumpStatusRepository(pool))
	defer func() {
		log.Log("Shutting down HTTP server...")
		server.Close()
	}()

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
	log.Log("Received termination signal, shutting down...")

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
