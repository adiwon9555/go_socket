package main

import (
	"encoding/json"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

//Message is
type Message struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

//Channel is
type Channel struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func main1() {
	rawMsg := []byte(`{"name":"channel add",` +
		`"data":{"name":"Hardware Support"}}`)
	var stdata Message
	err := json.Unmarshal(rawMsg, &stdata)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fmt.Printf("%+v\n", stdata)
	if stdata.Name == "channel add" {
		channel, err := addChannel(stdata.Data)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		var sendMessage Message
		sendMessage.Name = "channel add"
		sendMessage.Data = channel
		jsonMessageToSend, err := json.Marshal(sendMessage)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		fmt.Println("sent message is", string(jsonMessageToSend))

	}
	// jsonData, err := json.Marshal(stdata)
	// fmt.Println("JSON data is ", string(jsonData))
}
func addChannel(data interface{}) (Channel, error) {
	var channel Channel
	// channelMap := data.(map[string]interface{})
	// channel.Name = channelMap["name"].(string)
	err := mapstructure.Decode(data, &channel)
	if err != nil {
		fmt.Println("Error: ", err)
		return channel, err
	}
	channel.Id = 1
	fmt.Printf("%+v\n", channel)
	return channel, nil
}
