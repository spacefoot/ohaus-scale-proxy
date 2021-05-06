package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var bind = flag.String("bind", "127.0.0.1:23193", "Bind address")
var addr = flag.String("addr", "127.0.0.1", "Scale address")

func init() {
	flag.Parse()
}

func main() {
	scale := NewScale()
	manager := NewManager(scale)

	http.HandleFunc("/ws", manager.Handler())

	go manager.Run()
	go scale.Run(fmt.Sprintf("%s:9761", *addr))

	log.Printf("Start server on %s", *bind)
	log.Fatal(http.ListenAndServe(*bind, nil))
}
