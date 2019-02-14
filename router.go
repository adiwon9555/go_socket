package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

//Handler is
type Handler func(client *Client, data interface{})

//Router is
type Router struct {
	rules map[string]Handler
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

//NewRouter is Router initializer
func NewRouter() *Router {
	return &Router{
		rules: make(map[string]Handler),
	}
}
func (router *Router) findHandler(msgName string) (Handler, bool) {
	handler, found := router.rules[msgName]
	return handler, found

}

//Handle is
func (router *Router) Handle(msgName string, handler Handler) {
	router.rules[msgName] = handler
}
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	client := NewClient(socket, router.findHandler)
	go client.Write()
	client.Read()
}
