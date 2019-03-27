package main

import (
	"fmt"

	r "github.com/dancannon/gorethink"
	"github.com/gorilla/websocket"
)

//Client is
type Client struct {
	socket       *websocket.Conn
	send         chan Message
	findHandler  FindHandler
	session      *r.Session
	stopChannels map[int]chan bool
	userName     string
	id           string
}

//FindHandler is
type FindHandler func(string) (Handler, bool)

//NewClient is
func NewClient(socket *websocket.Conn, findHandler FindHandler, session *r.Session) *Client {
	var user User
	user.Name = "anonymous"
	res, err := r.Table("user").Insert(user).RunWrite(session)
	if err != nil {
		fmt.Print("Error:", err)
	}
	var id string
	if len(res.GeneratedKeys) > 0 {
		id = res.GeneratedKeys[0]
	}
	return &Client{
		socket:       socket,
		send:         make(chan Message),
		findHandler:  findHandler,
		session:      session,
		stopChannels: make(map[int]chan bool),
		userName:     user.Name,
		id:           id,
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
	err := r.Table("user").Get(c.id).Delete().Exec(c.session)
	if err != nil {
		fmt.Print("Error:-", err)
		return
	}
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

func (client *Client) subscribe(cursor *r.Cursor, category string, stop chan bool) {

	result := make(chan r.ChangeResponse)
	// go func() {
	// 	var changeResponse r.ChangeResponse
	// 	for cursor.Next(&changeResponse) {
	// 		result <- changeResponse
	// 	}
	// }()
	cursor.Listen(result)
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
				client.send <- Message{Name: category + " remove", Data: changeResponse.OldValue}
				fmt.Println("Sent " + category + " remove message")
			}
		case <-stop:
			cursor.Close()
			return

		}
	}
}
