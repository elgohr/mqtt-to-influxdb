package influx_test

import (
	"context"
	"fmt"
	"github.com/elgohr/mqtt-to-influxdb/influx"
	"github.com/elgohr/mqtt-to-influxdb/shared"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"testing"
)

const (
	TestOrganization = "org"
	TestBucket       = "buck"
)

func TestStorage(t *testing.T) {
	cmd, err := RunInflux()
	require.NoError(t, err)
	defer cmd.Process.Kill()

	require.NoError(t, os.Setenv(influx.ServerUrl, "http://localhost:8086"))
	defer os.Unsetenv(influx.ServerUrl)

	tc := influxdb2.NewClient("http://localhost:8086", "")
	defer tc.Close()

	home, err := os.UserHomeDir()
	require.NoError(t, err)
	dbPath := path.Join(home, ".influxdbv2")
	defer os.RemoveAll(dbPath)

	res, err := tc.Setup(context.Background(), "admin", "admin", TestOrganization, TestBucket, 1)
	require.NoError(t, err)

	require.NoError(t, os.Setenv(influx.Token, *res.Auth.Token))
	defer os.Unsetenv(influx.Token)

	require.NoError(t, os.Setenv(influx.Organization, TestOrganization))
	defer os.Unsetenv(influx.Organization)

	require.NoError(t, os.Setenv(influx.Bucket, TestBucket))
	defer os.Unsetenv(influx.Bucket)

	s, err := influx.NewStorage()
	require.NoError(t, err)

	topic := "topic"
	s.Write(shared.Message{
		Topic: topic,
		Value: []byte("test-string"),
	})

	t.Run("don't write same values twice", func(t *testing.T) {
		s.Write(shared.Message{
			Topic: topic,
			Value: []byte("test-string"),
		})
	})

	s.Shutdown() // for flushing

	qry := fmt.Sprintf(`from(bucket:"%s")|> range(start: -1h) |> filter(fn: (r) => r._measurement == "%s")`, TestBucket, topic)
	qres, err := tc.QueryAPI(TestOrganization).Query(context.Background(), qry)
	require.NoError(t, err)
	require.NoError(t, qres.Err())

	records := []testRecord{}
	for n := true; n == true; n = qres.Next() {
		record := qres.Record()
		if record != nil {
			records = append(records, testRecord{
				measurement: record.Measurement(),
				value:       record.Values()["_value"],
			})
		}
	}

	require.Equal(t, 1, len(records))
	require.Equal(t, topic, records[0].measurement)
	require.Equal(t, "test-string", fmt.Sprintf("%v", records[0].value))
}

type testRecord struct {
	measurement string
	value       interface{}
}
