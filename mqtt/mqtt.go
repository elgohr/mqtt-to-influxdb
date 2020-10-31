package mqtt

import (
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/elgohr/mqtt-to-influxdb/shared"
	uuid "github.com/satori/go.uuid"
	"log"
	"net/url"
	"os"
	"time"
)

const (
	ClientName = "MQTT_NAME"
	Url        = "MQTT_URL"
	Username   = "MQTT_USERNAME"
	Password   = "MQTT_PASSWORD"
)

type Collector struct {
	mqtt paho.Client
}

func NewCollector() (*Collector, error) {
	u, err := url.Parse(os.Getenv(Url))
	if err != nil {
		return nil, err
	}

	mqtt := paho.NewClient(&paho.ClientOptions{
		Servers:              []*url.URL{u},
		ClientID:             ClientId(),
		Username:             os.Getenv(Username),
		Password:             os.Getenv(Password),
		KeepAlive:            30,
		PingTimeout:          10 * time.Second,
		ConnectTimeout:       30 * time.Second,
		MaxReconnectInterval: 10 * time.Minute,
		AutoReconnect:        true,
	})
	if err := check(mqtt.Connect()); err != nil {
		return nil, err
	}
	return &Collector{mqtt: mqtt}, nil
}

func (c *Collector) Collect() <-chan shared.Message {
	messages := make(chan shared.Message)

	const allMessages = "#"
	if err := check(c.mqtt.Subscribe(allMessages, 0, func(c paho.Client, m paho.Message) {
		messages <- shared.Message{
			Topic: m.Topic(),
			Value: m.Payload(),
		}
	})); err != nil {
		log.Println(err)
	}

	return messages
}

func (c *Collector) Shutdown() {
	log.Println("Shutting down MQTT client")
	c.mqtt.Disconnect(200)
}

func ClientId() string {
	if name := os.Getenv(ClientName); name != "" {
		return name
	}
	return uuid.NewV4().String()
}

func check(t paho.Token) error {
	if err := t.Error(); t.Wait() && err != nil {
		return err
	}
	return nil
}
