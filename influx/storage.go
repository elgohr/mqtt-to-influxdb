package influx

import (
	"github.com/elgohr/mqtt-to-influxdb/shared"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"log"
	"os"
	"strconv"
)

const (
	ServerUrl    = "INFLUX_URL"
	Token        = "INFLUX_TOKEN"
	Organization = "INFLUX_ORGANIZATION"
	Bucket       = "INFLUX_BUCKET"

	RetryInterval = "INFLUX_RETRY_INTERVAL"
	MaxRetries    = "INFLUX_MAX_RETRIES"
	BatchSize     = "INFLUX_BATCH_SIZE"
)

type Storage struct {
	client influxdb2.Client
	writer api.WriteAPI
}

func NewStorage() (*Storage, error) {
	serverUrl := GetEnvDefault(ServerUrl, "http://localhost:8086")
	token := os.Getenv(Token)

	org := os.Getenv(Organization)
	bucket := os.Getenv(Bucket)

	config := loadEnvironmentConfiguration(influxdb2.DefaultOptions())
	client := influxdb2.NewClientWithOptions(serverUrl, token, config)

	return &Storage{
		client: client,
		writer: client.WriteAPI(org, bucket),
	}, nil
}

func (s Storage) Write(msg shared.Message) {
	s.writer.WritePoint(influxdb2.NewPoint(
		msg.Topic,
		map[string]string{},
		map[string]interface{}{"value": msg.Value},
		msg.Time),
	)
}

func (s *Storage) Shutdown() {
	s.writer.Flush()
	s.client.Close()
}

func GetEnvDefault(key string, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}

func loadEnvironmentConfiguration(config *influxdb2.Options) *influxdb2.Options {
	config = setWhenPresent(config, RetryInterval, func(config *influxdb2.Options, value string) *influxdb2.Options {
		ui, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			log.Printf("%s : %v \n", RetryInterval, err)
			return config
		}
		return config.SetRetryInterval(uint(ui))
	})
	config = setWhenPresent(config, MaxRetries, func(config *influxdb2.Options, value string) *influxdb2.Options {
		ui, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			log.Printf("%s : %v \n", MaxRetries, err)
			return config
		}
		return config.SetMaxRetries(uint(ui))
	})
	config = setWhenPresent(config, BatchSize, func(config *influxdb2.Options, value string) *influxdb2.Options {
		ui, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			log.Printf("%s : %v \n", BatchSize, err)
			return config
		}
		return config.SetBatchSize(uint(ui))
	})
	return config
}

func setWhenPresent(config *influxdb2.Options, key string, changer envConfigChanger) *influxdb2.Options {
	val := os.Getenv(key)
	if val != "" {
		return changer(config, val)
	}
	return config
}

type envConfigChanger func(config *influxdb2.Options, value string) *influxdb2.Options
