package main

import (
	"flag"
	"log"
	"net"
	"os"
)

// ClientManager datastruct
type ClientManager struct {
	ClientsRoom map[int]map[*Client]bool
	Register    chan *Client
	Deregister  chan *Client
}

// Client datastruct
type Client struct {
	RoomID int
	Conn   net.Conn
}

func main() {
	mode := flag.String("mode", "server", "server | client")
	host := flag.String("host", "localhost:5000", "<address>:<port>")
	roomID := flag.Int("room", 0, "room id for client only")
	flag.Parse()
	if *mode == "server" {
		ServerMode(*host)
	} else {
		ClientMode(*host, *roomID)
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
