package main

import (
	"net/http"

	"github.com/Nino-K/gononymous/server"
)

func main() {
	server := server.New()
	http.HandleFunc("/", server.Home)
	http.ListenAndServe("localhost:9988", nil)
}
