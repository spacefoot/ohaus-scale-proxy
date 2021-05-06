package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Address string
	Bind    string
}

var config = Config{
	Address: "127.0.0.1",
	Bind:    "127.0.0.1:23193",
}

func init() {
	addr := flag.String("addr", "", "Scale address")
	bind := flag.String("bind", "", "Bind address")
	configFile := flag.String("config", "", "Config file")

	flag.Parse()

	if *configFile != "" {
		data, err := ioutil.ReadFile(*configFile)
		if err != nil {
			log.Panic("File reading error", err)
		}
		yaml.Unmarshal(data, &config)
	}

	if *addr != "" {
		config.Address = *addr
	}

	if *bind != "" {
		config.Bind = *bind
	}
}

func main() {
	scale := NewScale()
	manager := NewManager(scale)

	http.HandleFunc("/ws", manager.Handler())

	go manager.Run()
	go scale.Run(fmt.Sprintf("%s:9761", config.Address))

	log.Printf("Use scale %s", config.Address)
	log.Printf("Start server on %s", config.Bind)
	log.Fatal(http.ListenAndServe(config.Bind, nil))
}
