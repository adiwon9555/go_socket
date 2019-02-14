package backup

import (
	"fmt"
	"math/rand"
	"time"
)

type Message struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

//Client is for Reading and Writing messages to web socket
type Client struct {
	send chan Message
}

func (client *Client) write() {
	for data := range client.send {
		//socket WriteJson
		fmt.Printf("%+v\n", data)
	}
}
func (client *Client) subscribeChannels() {
	for {
		time.Sleep(r())
		outMessage := Message{
			Name: "channel add",
			Data: "",
		}
		client.send <- outMessage
	}

}
func (client *Client) subscribeMessages() {
	for {
		time.Sleep(r())
		outMessage := Message{
			Name: "Message add",
			Data: "",
		}
		client.send <- outMessage
	}
}
func r() time.Duration {
	return time.Millisecond * time.Duration(rand.Intn(1000))
}

func NewClient() *Client {
	return &Client{
		send: make(chan Message),
	}

}

func main2() {
	client := NewClient()
	go client.subscribeChannels()
	go client.subscribeMessages()
	client.write()

}
