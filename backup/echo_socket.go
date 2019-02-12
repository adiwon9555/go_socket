package backup

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func main() {

	http.HandleFunc("/", handle)
	// http.ListenAndServe(":4001", hand{})
	http.ListenAndServe(":4001", nil)
	// http.ListenAndServe(":4001", http.FileServer(http.Dir("../")))
}
func handle(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	for {
		msgType, msg, err := socket.ReadMessage()
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		fmt.Println(string(msg))
		if err := socket.WriteMessage(msgType, msg); err != nil {
			fmt.Println("Error: ", err)
			return
		}
	}
	fmt.Fprintf(w, "hello from Server")
}
