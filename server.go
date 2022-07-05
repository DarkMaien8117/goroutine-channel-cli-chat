package main

import (
	"encoding/binary"
	"log"
	"net"
)

// ServerMode : run on server mode
func ServerMode(host string) {
	listener, err := net.Listen("tcp", host)
	checkError(err)
	log.Println("Server listen on", host)

	clientManager := ClientManager{
		ClientsRoom: make(map[int]map[*Client]bool, 8),
		Register:    make(chan *Client),
		Deregister:  make(chan *Client),
	}

	// start client manager
	go clientManager.start()

	// listen if any client join
	for {
		//accept some client's connection
		conn, err := listener.Accept()
		checkError(err)

		//create new temp client and matching accepted connection
		client := &Client{
			Conn: conn,
		}

		//and set client's roomID from the accepted connection
		clientManager.prepareClient(client)

		//resgister accepted client at clientManager
		clientManager.Register <- client

		//and broadcast accepted client's message to all clients in roomID
		go clientManager.serveClient(client)
	}
}

// Start : Start client manager for register and deregister any client
func (clientManager *ClientManager) start() {
	for {
		select {
		case client := <-clientManager.Register:
			log.Printf("Register client %v\n", *client)
			if _, ok := clientManager.ClientsRoom[client.RoomID]; !ok {
				clientManager.ClientsRoom[client.RoomID] = make(map[*Client]bool, 4)
			}
			clientManager.ClientsRoom[client.RoomID][client] = true
		case client := <-clientManager.Deregister:
			log.Printf("Deregister client %v\n", *client)
			client.Conn.Close()
			delete(clientManager.ClientsRoom[client.RoomID], client)
		}
	}
}

//set client's roomID for accepted client's roomID
func (clientManager *ClientManager) prepareClient(client *Client) {
	// Get roomID first
	bs := make([]byte, 4)
	_, err := client.Conn.Read(bs)
	checkError(err)
	roomID := int(binary.LittleEndian.Uint16(bs))
	client.RoomID = roomID
	log.Printf("Client %v join\n", *client)
}

//broadcast accepted conn's message to all clients in conn's roomID
func (clientManager *ClientManager) serveClient(client *Client) {
	log.Printf("Serving client %v\n", *client)
	for {
		buff := make([]byte, 4096)
		_, err := client.Conn.Read(buff)
		if err != nil {
			clientManager.Deregister <- client
			break
		}

		// Broadcast to every client in roomID
		for anotherClient := range clientManager.ClientsRoom[client.RoomID] {
			if client != anotherClient {
				anotherClient.Conn.Write(buff)
			}
		}
	}
}
