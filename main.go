package main

import (
	"log"
	"net/http"

	r "github.com/dancannon/gorethink"
)

func main() {
	session, err := r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: "rtsupport",
	})
	if err != nil {
		log.Panic("Error:- ", err)
		return
	}
	router := NewRouter(session)
	router.Handle("channel add", addChannel)
	router.Handle("channel subscribe", subscribeChannel)
	http.Handle("/", router)
	http.ListenAndServe(":4001", nil)
}
