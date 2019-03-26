package main

import (
	"fmt"
	"time"

	r "github.com/dancannon/gorethink"
	"github.com/mitchellh/mapstructure"
)

const (
	channelStop = iota
	messageStop
	userStop
)

func subscribeChannel(client *Client, data interface{}) {
	stop := client.NewStopChannel(channelStop)
	client.subscribe(data, "channel", stop)
	//Replace with changeFeed from Rethink Db that will look up for channels and then
	//block/wait until add,remove or edit operation in channels data in db
	// stop := client.NewStopChannel(channelStop)
	// cursor, err := r.Table("channel").Changes(r.ChangesOpts{IncludeInitial: true}).Run(client.session)
	// if err != nil {
	// 	client.send <- Message{Name: "error", Data: err.Error()}
	// 	return
	// }
	// result := make(chan r.ChangeResponse)
	// go func() {
	// 	var changeResponse r.ChangeResponse
	// 	for cursor.Next(&changeResponse) {
	// 		result <- changeResponse
	// 	}
	// }()
	// for {
	// 	select {
	// 	case changeResponse := <-result:
	// 		if changeResponse.NewValue != nil && changeResponse.OldValue == nil {
	// 			client.send <- Message{Name: "channel add", Data: changeResponse.NewValue}
	// 			fmt.Println("Sent channel add message")
	// 		}
	// 	case <-stop:
	// 		cursor.Close()
	// 		return

	// 	}
	// }

	// data := changeResponse.NewValue
	// var channel Channel
	// err := mapstructure.Decode(data, &channel)
	// if err != nil {
	// 	client.send <- Message{Name: "error", Data: err.Error()}
	// 	return
	// }
	// fmt.Printf("%+v\n", channel)
	// client.send <- Message{Name: "channel add", Data: channel}
	//we dont need to convert data to channel as socket.WriteJson can automatically convert map type to json

	// for {
	// 	time.Sleep(time.Second * 2)
	// 	msg := Message{
	// 		Name: "channel add",
	// 		Data: Channel{
	// 			ID:   1,
	// 			Name: "Software Support",
	// 		},
	// 	}
	// 	client.send <- msg
	// 	fmt.Printf("Sent new channel")
	// }

}

func unsubscribeChannel(client *Client, data interface{}) {
	client.StopForKey(channelStop)
}
func subscribeUser(client *Client, data interface{}) {
	stop := client.NewStopChannel(userStop)
	addUser(client)
	client.subscribe(data, "user", stop)
}
func unsubscribeUser(client *Client, data interface{}) {
	// go removeUser(client, data)
	client.StopForKey(userStop)
}
func subscribeMessage(client *Client, data interface{}) {
	stop := client.NewStopChannel(messageStop)
	client.subscribe(data, "message", stop)
}
func unsubscribeMessage(client *Client, data interface{}) {
	client.StopForKey(messageStop)
}
func editUser(client *Client, data interface{}) {
	var user User
	err := mapstructure.Decode(data, &user)
	if err != nil {
		client.send <- Message{Name: "error", Data: err.Error()}
		return
	}
	// channel.ID = 1
	fmt.Println("User edited")
	fmt.Printf("%+v\n", user)
	err = r.Table("user").Get(user.ID).Update(user).Exec(client.session)
	if err != nil {
		client.send <- Message{Name: "error", Data: err.Error()}
		return
	}
}
func addMessage(client *Client, data interface{}) {
	var msg Msg
	err := mapstructure.Decode(data, &msg)
	if err != nil {
		client.send <- Message{Name: "error", Data: err.Error()}
		return
	}
	msg.CreatedAt = time.Now().Unix()
	// channel.ID = 1
	fmt.Println("Message added")
	fmt.Printf("%+v\n", msg)
	err = r.Table("message").Insert(msg).Exec(client.session)
	if err != nil {
		client.send <- Message{Name: "error", Data: err.Error()}
		return
	}
}
func addChannel(client *Client, data interface{}) {
	var channel Channel
	err := mapstructure.Decode(data, &channel)
	if err != nil {
		client.send <- Message{Name: "error", Data: err.Error()}
		return
	}
	// channel.ID = 1
	fmt.Println("Channel added")
	fmt.Printf("%+v\n", channel)
	err = r.Table("channel").Insert(channel).Exec(client.session)
	if err != nil {
		client.send <- Message{Name: "error", Data: err.Error()}
		return
	}
	//Replace with adding channel to db
}
func addUser(client *Client) {
	// var user User
	// err := mapstructure.Decode(data, &user)
	// if err != nil {
	// 	client.send <- Message{Name: "error", Data: err.Error()}
	// 	return
	// }
	// channel.ID = 1
	var user = User{Name: "anonymous"}
	fmt.Println("User added")
	fmt.Printf("%+v\n", user)
	err := r.Table("user").Insert(user).Exec(client.session)
	if err != nil {
		client.send <- Message{Name: "error", Data: err.Error()}
		return
	}
	//Replace with adding channel to db
}
func removeUser(client *Client, data interface{}) {
	var user User
	err := mapstructure.Decode(data, &user)
	if err != nil {
		client.send <- Message{Name: "error", Data: err.Error()}
		return
	}

	fmt.Println("User Removed")
	fmt.Printf("%+v\n", user)
	err = r.Table("user").Get(user.ID).Delete().Exec(client.session)
	if err != nil {
		client.send <- Message{Name: "error", Data: err.Error()}
		return
	}
	//Replace with adding channel to db
}
