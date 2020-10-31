package mqtt_test

import (
	"context"
	"fmt"
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/elgohr/mqtt-to-influxdb/mqtt"
	"github.com/elgohr/mqtt-to-influxdb/shared"
	"github.com/fhmq/hmq/broker"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"net"
	"net/url"
	"os"
	"testing"
	"time"
)

func TestControl(t *testing.T) {
	b, err := broker.NewBroker(broker.DefaultConfig)
	require.NoError(t, err)
	b.Start()

	var con net.Conn
	for ; con == nil; con, _ = net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", "1883"), time.Second) {
		fmt.Println("waiting")
		time.Sleep(time.Second)
	}

	address := "tcp://127.0.0.1:1883"

	u, err := url.Parse(address)
	require.NoError(t, err)

	require.NoError(t, os.Setenv(mqtt.Url, address))
	defer os.Unsetenv(mqtt.Url)

	c := paho.NewClient(&paho.ClientOptions{
		Servers:        []*url.URL{u},
		ClientID:       "test-client",
		KeepAlive:      5,
		PingTimeout:    1 * time.Second,
		ConnectTimeout: 1 * time.Second,
	})

	defer c.Disconnect(10)
	if to := c.Connect(); to.Wait() && to.Error() != nil {
		require.NoError(t, to.Error())
	}

	col, err := mqtt.NewCollector()
	require.NoError(t, err)
	defer col.Shutdown()

	ic := col.Collect()

	expectedValue := uuid.NewV4().String()
	expectedTopic := uuid.NewV4().String()

	if to := c.Publish(expectedTopic, 0, true, expectedValue); to.Wait() && to.Error() != nil {
		require.NoError(t, to.Error())
	}

	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(1*time.Second))
	select {
	case calledWith := <-ic:
		require.Equal(t, shared.Message{
			Topic: expectedTopic,
			Value: []byte(expectedValue),
		}, calledWith)
	case <-ctx.Done():
		require.True(t, false, "got timeout")
	}
}
