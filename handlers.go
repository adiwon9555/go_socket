package main

import (
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
)

func addChannel(client *Client, data interface{}) {
	var channel Channel
	err := mapstructure.Decode(data, &channel)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	channel.ID = 1
	fmt.Println("Channel added")
	fmt.Printf("%+v\n", channel)
	//Replace with adding channe to db
}
func subscribeChannel(client *Client, data interface{}) {
	//Replace with changeFeed from Rethink Db that will look up for channels and then
	//block/wait until add,remove or edit operation in channels data in db
	for {
		time.Sleep(time.Second * 2)
		msg := Message{
			Name: "channel add",
			Data: Channel{
				ID:   1,
				Name: "Software Support",
			},
		}
		client.send <- msg
		fmt.Printf("Sent new channel")

	}

}
