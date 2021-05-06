package main

import (
	"bufio"
	"log"
	"net"
	"strconv"
	"strings"
)

type Scale struct {
	weight    chan float64
	connected chan bool
}

func NewScale() *Scale {
	return &Scale{
		make(chan float64),
		make(chan bool),
	}
}

func (s *Scale) reader(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		data, err := reader.ReadString('\n')
		if err != nil {
			log.Println("dial read:", err)
			return
		}

		if strings.HasPrefix(data, "      ") {
			part := strings.Fields(data)
			if len(part) >= 3 && part[2] == "N" {
				if weight, err := strconv.ParseFloat(part[0], 64); err == nil {
					log.Println("weight:", weight)
					s.weight <- weight
				}
			}
		}
	}
}

func (s *Scale) Run(addr string) {
	for {
		if conn, err := net.Dial("tcp", addr); err == nil {
			log.Println("Scale connected")
			s.connected <- true

			s.reader(conn)

			log.Println("Scale disconnected")
			s.connected <- false
		}
	}
}
