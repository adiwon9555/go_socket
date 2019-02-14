package backup

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
)

// //Message is
// type Message struct {
// 	Name string      `json:"name"`
// 	Data interface{} `json:"data"`
// }

// //Channel is
// type Channel struct {
// 	Id   int    `json:"id"`
// 	Name string `json:"name"`
// }

// var upgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// 	CheckOrigin:     func(r *http.Request) bool { return true },
// }

func main3() {

	http.HandleFunc("/", handle1)
	// http.ListenAndServe(":4001", hand{})
	http.ListenAndServe(":4001", nil)
	// http.ListenAndServe(":4001", http.FileServer(http.Dir("../")))
}
func handle1(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	for {
		var inMessage Message
		var outMessage Message
		if err := socket.ReadJSON(&inMessage); err != nil {
			fmt.Println("Error: ", err)
			break
		}
		// fmt.Printf("%+v\n", inMessage)

		switch inMessage.Name {
		case "channel add":
			err := addChannel1(inMessage.Data)
			if err != nil {
				fmt.Println("Error: ", err)
				outMessage = Message{"error", err}
				if err := socket.WriteJSON(outMessage); err != nil {
					fmt.Println("Error: ", err)
					break
				}
			}
		case "channel subscribe":
			go subscribeChannel(socket)
		}

	}
}
func addChannel1(data interface{}) error {
	var channel Channel
	err := mapstructure.Decode(data, &channel)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	channel.Id = 1
	fmt.Println("Channel added")
	fmt.Printf("%+v\n", channel)
	//Replace with adding channe to db
	return nil
}
func subscribeChannel(socket *websocket.Conn) {
	//Replace with changeFeed from Rethink Db that will look up for channels and then
	//block/wait until add,remove or edit operation in channels data in db
	for {
		time.Sleep(time.Second * 2)
		msg := Message{
			Name: "channel add",
			Data: Channel{
				Id:   1,
				Name: "Software Support",
			},
		}
		// socket.WriteJSON(msg)
		if err := socket.WriteJSON(msg); err != nil {
			fmt.Println("Error: ", err)
			break
		}

		fmt.Printf("Sent new channel")

	}

}
