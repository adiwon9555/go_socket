package main

import (
	"fmt"
	"time"

	r "github.com/dancannon/gorethink"
	"github.com/mitchellh/mapstructure"
)

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

func subscribeChannel(client *Client, data interface{}) {
	//Replace with changeFeed from Rethink Db that will look up for channels and then
	//block/wait until add,remove or edit operation in channels data in db
	// cursor, err := r.Table("channel").Changes(r.ChangesOpts{IncludeInitial: true}).Run(client.session)
	// if err != nil {
	// 	client.send <- Message{Name: "error", Data: err.Error()}
	// 	return
	// }
	// var changeResponse r.ChangeResponse
	// for cursor.Next(&changeResponse) {
	// 	fmt.Printf("%+v\n", changeResponse)
	// }
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
