package influx_test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"net"
	"os"
	"os/exec"
	"testing"
	"time"
)

func RunInflux(t *testing.T) *exec.Cmd {
	cmd := exec.Command("testdata/influxd")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	require.NoError(t, cmd.Start())
	var con net.Conn
	for ; con == nil; con, _ = net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", "8086"), time.Second) {
		fmt.Println("waiting")
		time.Sleep(time.Second)
	}
	return cmd
}
