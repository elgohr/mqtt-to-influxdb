package main

import (
	"context"
	"github.com/elgohr/mqtt-to-influxdb/influx"
	"github.com/elgohr/mqtt-to-influxdb/mqtt"
	"log"
	"os"
	"os/signal"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)

	col, err := mqtt.NewCollector()
	if err != nil {
		log.Fatalln(err)
	}

	sto, err := influx.NewStorage()
	if err != nil {
		log.Fatalln(err)
	}

	c := col.Collect()
	for {
		select {
		case msg := <-c:
			log.Println(msg)
			sto.Write(ctx, msg)
		case <-ctx.Done():
			log.Println("exiting")
			col.Shutdown()
			sto.Shutdown()
			return
		}
	}
}
