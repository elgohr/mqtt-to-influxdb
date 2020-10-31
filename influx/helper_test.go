package influx_test

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"
)

func RunInflux() (*exec.Cmd, error) {
	cmd := exec.Command("testdata/influxd")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	var con net.Conn
	for ; con == nil; con, _ = net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", "8086"), time.Second) {
		fmt.Println("waiting")
		time.Sleep(time.Second)
	}

	return cmd, nil
}
