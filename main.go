package main

import (
	"log"
	"net/http"

	_ "net/http/pprof"

	"github.com/Nino-K/gononymous/server"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	server := server.New()
	http.HandleFunc("/", server.Home)
	http.ListenAndServe("localhost:9988", nil)
}
