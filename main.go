package main

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/Nino-K/gononymous/handler"
	"github.com/Nino-K/gononymous/server"
	"github.com/gorilla/websocket"
)

func main() {
	upgrader := websocket.Upgrader{}
	sessonManager := server.NewSessionManager()
	sessionHandler := handler.NewSessionHandler(sessonManager, &upgrader)
	http.HandleFunc("/", sessionHandler.Join)
	http.ListenAndServe("localhost:9988", nil)
}

func sessionId(url *url.URL) string {
	segments := strings.SplitAfterN(url.Path, "/", 2)
	return strings.Join(segments[1:], "/")

}
