package influx

import (
	"context"
	"fmt"
	"github.com/elgohr/mqtt-to-influxdb/shared"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"os"
	"sync"
	"time"
)

const (
	ServerUrl    = "INFLUX_URL"
	Token        = "INFLUX_TOKEN"
	Organization = "INFLUX_ORGANIZATION"
	Bucket       = "INFLUX_BUCKET"
)

type Storage struct {
	client         influxdb2.Client
	writer         api.WriteAPIBlocking
	cache          sync.Map
	lastValueCache sync.Map
}

func NewStorage() (*Storage, error) {
	serverUrl := os.Getenv(ServerUrl)
	if serverUrl == "" {
		serverUrl = "http://localhost:8086"
	}
	token := os.Getenv(Token)
	org := os.Getenv(Organization)
	bucket := os.Getenv(Bucket)

	client := influxdb2.NewClientWithOptions(serverUrl, token, influxdb2.DefaultOptions())
	return &Storage{
		client:         client,
		writer:         client.WriteAPIBlocking(org, bucket),
		cache:          sync.Map{},
		lastValueCache: sync.Map{},
	}, nil
}

func (s *Storage) Write(ctx context.Context, msg shared.Message) {
	if val, exists := s.lastValueCache.Load(msg.Topic); exists && val.(string) == fmt.Sprintf("%v", msg.Value) {
		return
	}
	if err := s.write(ctx, msg); err != nil {
		s.cache.Store(msg.Hash(), msg)
		return
	}
	s.lastValueCache.Store(msg.Topic, fmt.Sprintf("%v", msg.Value))
	s.cache.Range(func(key, value interface{}) bool {
		msg := value.(shared.Message)
		if err := s.write(ctx, msg); err == nil {
			s.cache.Delete(key)
		}
		return true
	})
}

func (s *Storage) write(ctx context.Context, msg shared.Message) error {
	tCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	return s.writer.WritePoint(tCtx, influxdb2.NewPoint(
		msg.Topic,
		map[string]string{},
		map[string]interface{}{"value": msg.Value},
		msg.Time))
}

func (s *Storage) Shutdown() {
	s.client.Close()
}
