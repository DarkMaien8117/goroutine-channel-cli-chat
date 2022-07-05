package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
)

// ClientMode : run on client mode
func ClientMode(host string, roomID int) {
	conn, err := net.Dial("tcp", host)
	defer conn.Close()
	checkError(err)
	log.Println("Client connect on", host)

	client := &Client{
		RoomID: roomID,
		Conn:   conn,
	}

	// Prepare, send roomID to server
	client.prepare() //first send to server -> RoomID

	// Goroutine to receive message
	go client.receive() //Show the message from network

	// Ready to send message,so catch inputed text and broadcast to network this text
	reader := bufio.NewReader(os.Stdin)
	for {
		msg, err := reader.ReadString('\n')
		checkError(err)
		_, err = client.Conn.Write([]byte(msg))
		checkError(err)
	}
}

// Prepare : prepare client including send room ID to server
func (client *Client) prepare() {
	// Send room ID first
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint16(bs, uint16(client.RoomID))
	_, err := client.Conn.Write(bs)
	checkError(err)
	log.Println("Join on room", client.RoomID)
}

// Receive : receive any message from another client
func (client *Client) receive() {
	buff := make([]byte, 4096)
	for {
		_, err := client.Conn.Read(buff)
		checkError(err)
		fmt.Printf("%v: %s", client, string(bytes.Trim(buff, "\x00")))
	}
}
