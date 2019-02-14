package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)

//Message is
type Message struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

//Channel is
type Channel struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

//Client is
type Client struct {
	socket      *websocket.Conn
	send        chan Message
	findHandler FindHandler
}

//FindHandler is
type FindHandler func(string) (Handler, bool)

//NewClient is
func NewClient(socket *websocket.Conn, findHandler FindHandler) *Client {
	return &Client{
		socket:      socket,
		send:        make(chan Message),
		findHandler: findHandler,
	}

}
func (client *Client) Read() {
	var msg Message
	for {
		err := client.socket.ReadJSON(&msg)
		if err != nil {
			fmt.Println("Error: ", err)
			break
		}
		if handler, found := client.findHandler(msg.Name); found {
			go handler(client, msg.Data)
		}
	}

	client.socket.Close()
}

func (client *Client) Write() {
	for msg := range client.send {
		err := client.socket.WriteJSON(msg)
		if err != nil {
			fmt.Println("Error: ", err)
			break
		}
	}
	client.socket.Close()
}
