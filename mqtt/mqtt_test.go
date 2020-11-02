package mqtt_test

import (
	"context"
	"fmt"
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/elgohr/mqtt-to-influxdb/mqtt"
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

	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(20*time.Second))

	dyn := uuid.NewV4().String()
	for _, scenario := range []struct {
		name   string
		topic  string
		input  string
		output interface{}
	}{
		{
			name:   "dynamic",
			topic:  dyn,
			input:  dyn + "t",
			output: dyn + "t",
		},
		{
			name:   "strings",
			topic:  "expected-string",
			input:  `["periodic"]`,
			output: `["periodic"]`,
		},
		{
			name:   "integers",
			topic:  "expected-int",
			input:  "1",
			output: 1,
		},
		{
			name:   "floats",
			topic:  "expected-int",
			input:  "1.1",
			output: 1.1,
		},
		{
			name:  "json",
			topic: "expected-json",
			input: `{"clientID":"test-client","online":true,"timestamp":"2020-11-01T18:35:52Z"}`,
			output: map[string]interface{}{
				"clientID":  "test-client",
				"online":    true,
				"timestamp": "2020-11-01T18:35:52Z",
			},
		},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			if to := c.Publish(scenario.topic, 0, true, scenario.input); to.Wait() && to.Error() != nil {
				require.NoError(t, to.Error())
			}

			select {
			case calledWith := <-ic:
				require.Equal(t, scenario.topic, calledWith.Topic)
				require.Equal(t, scenario.output, calledWith.Value)
				require.WithinDuration(t, time.Now(), calledWith.Time, time.Second)
			case <-ctx.Done():
				require.True(t, false, "got timeout")
			}
		})
	}
}
