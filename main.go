package main

import (
	"net/http"
)

func main() {
	router := NewRouter()
	router.Handle("channel add", addChannel)
	router.Handle("channel subscribe", subscribeChannel)
	http.Handle("/", router)
	http.ListenAndServe(":4001", nil)
}
