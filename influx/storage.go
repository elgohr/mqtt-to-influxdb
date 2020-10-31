package influx

import (
	"github.com/elgohr/mqtt-to-influxdb/shared"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"os"
	"time"
)

const (
	ServerUrl    = "INFLUX_URL"
	Token        = "INFLUX_TOKEN"
	Organization = "INFLUX_ORGANIZATION"
	Bucket       = "INFLUX_BUCKET"
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

	config := influxdb2.DefaultOptions().SetBatchSize(10)
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
		time.Now()),
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
