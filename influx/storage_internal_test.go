package influx

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestLoadEnvironmentWriteOptionConfiguration(t *testing.T) {
	for _, scenario := range []struct {
		Key      string
		Value    string
		Expected *influxdb2.Options
	}{
		{
			Key: RetryInterval, Value: "1",
			Expected: influxdb2.DefaultOptions().SetRetryInterval(1),
		},
		{
			Key: RetryInterval, Value: "wrong",
			Expected: influxdb2.DefaultOptions(),
		},
		{
			Key: MaxRetries, Value: "1",
			Expected: influxdb2.DefaultOptions().SetMaxRetries(1),
		},
		{
			Key: MaxRetries, Value: "wrong",
			Expected: influxdb2.DefaultOptions(),
		},
		{
			Key: BatchSize, Value: "1",
			Expected: influxdb2.DefaultOptions().SetBatchSize(1),
		},
		{
			Key: BatchSize, Value: "wrong",
			Expected: influxdb2.DefaultOptions(),
		},
	} {
		t.Run(scenario.Key+":"+scenario.Value, func(t *testing.T) {
			require.NoError(t, os.Setenv(scenario.Key, scenario.Value))
			defer os.Unsetenv(scenario.Key)

			s, err := NewStorage()
			require.NoError(t, err)
			require.Equal(t, scenario.Expected.WriteOptions(), s.client.Options().WriteOptions())
		})
	}
}
