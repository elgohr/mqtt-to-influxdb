package main

import (
	"github.com/elgohr/mqtt-to-influxdb/influx"
	"github.com/elgohr/mqtt-to-influxdb/mqtt"
	"log"
	"os"
	"os/signal"
)

func main() {
	q := make(chan os.Signal)
	signal.Notify(q, os.Kill, os.Interrupt)

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
			sto.Write(msg)
		case <-q:
			log.Println("exiting")
			col.Shutdown()
			sto.Shutdown()
			return
		}
	}
}
