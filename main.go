package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"mascot/src/app"
	"mascot/src/config"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	c := flag.String("c", "config.json", "Specify the configuration file.")
	flag.Parse()

	bytes, err := ioutil.ReadFile(*c)
	if err != nil {
		log.Fatalf("config file: %v\n", err)
	}

	conf := config.Config{}
	if err := json.Unmarshal(bytes, &conf); err != nil {
		log.Fatalf("config file: %v\n", err)
	}

	app, err := app.New(&conf)
	if err != nil {
		log.Fatalf("could not init application: %s", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	app.Run()

	<-stop
	app.Stop()
}
