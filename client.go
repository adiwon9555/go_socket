package main

import (
	"fmt"

	r "github.com/dancannon/gorethink"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
)

type Msg struct {
	ID        string `json:"id" gorethink:"id,omitempty"`
	Body      string `json:"body" gorethink:"body"`
	Author    string `json:"author" gorethink:"author"`
	CreatedAt int64  `json:"createdAt" gorethink:"createdAt"`
	ChannelId string `json:"channelId" gorethink:"channelId"`
}

//Message is
type Message struct {
	Name string      `json:"name" gorethink:"name"`
	Data interface{} `json:"data" gorethink:"data"`
}

//Channel is
type Channel struct {
	ID   string `json:"id" gorethink:"id,omitempty"`
	Name string `json:"name" gorethink:"name"`
}
type User struct {
	ID   string `json:"id" gorethink:"id,omitempty"`
	Name string `json:"name" gorethink:"name"`
}

//Client is
type Client struct {
	socket       *websocket.Conn
	send         chan Message
	findHandler  FindHandler
	session      *r.Session
	stopChannels map[int]chan bool
}

//FindHandler is
type FindHandler func(string) (Handler, bool)

//NewClient is
func NewClient(socket *websocket.Conn, findHandler FindHandler, session *r.Session) *Client {
	return &Client{
		socket:       socket,
		send:         make(chan Message),
		findHandler:  findHandler,
		session:      session,
		stopChannels: make(map[int]chan bool),
	}

}

func (c *Client) NewStopChannel(i int) chan bool {
	c.StopForKey(i)
	stop := make(chan bool)
	c.stopChannels[i] = stop
	return stop

}

func (c *Client) CloseConnections() {
	for _, stop := range c.stopChannels {
		stop <- true
	}
	close(c.send)
}

func (c *Client) StopForKey(key int) {
	if ch, found := c.stopChannels[key]; found {
		ch <- true
		delete(c.stopChannels, key)
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

func (client *Client) subscribe(data interface{}, category string, stop chan bool) {
	var cursor *r.Cursor
	var err error
	if category == "message" {
		var msg Msg
		err := mapstructure.Decode(data, &msg)
		if err != nil {
			client.send <- Message{Name: "error", Data: err.Error()}
			return
		}
		cursor, err = r.Table(category).Filter(r.Row.Field("channelId").Eq(msg.ChannelId)).Changes(r.ChangesOpts{IncludeInitial: true}).Run(client.session)
	} else {
		cursor, err = r.Table(category).Changes(r.ChangesOpts{IncludeInitial: true}).Run(client.session)
	}
	if err != nil {
		client.send <- Message{Name: "error", Data: err.Error()}
		return
	}
	result := make(chan r.ChangeResponse)
	go func() {
		var changeResponse r.ChangeResponse
		for cursor.Next(&changeResponse) {
			result <- changeResponse
		}
	}()
	for {
		select {
		case changeResponse := <-result:
			if changeResponse.NewValue != nil && changeResponse.OldValue == nil {
				client.send <- Message{Name: category + " add", Data: changeResponse.NewValue}
				fmt.Println("Sent " + category + " add message")
			} else if changeResponse.NewValue != nil && changeResponse.OldValue != nil {
				client.send <- Message{Name: category + " edit", Data: changeResponse.NewValue}
				fmt.Println("Sent " + category + " edit message")
			} else if changeResponse.NewValue == nil && changeResponse.OldValue != nil {
				client.send <- Message{Name: category + " remove", Data: changeResponse.NewValue}
				fmt.Println("Sent " + category + " remove message")
			}
		case <-stop:
			cursor.Close()
			return

		}
	}
}
