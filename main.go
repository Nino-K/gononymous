package main

import (
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/Nino-K/gononymous/server"
	"github.com/gorilla/websocket"
)

var sessonManager *server.SessionManager
var upgrader = websocket.Upgrader{}

func join(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Fatalf("websocket Upgrade: %s\n", err)
		return
	}

	peerId := r.Header.Get("CLIENT_ID")
	peer := server.NewPeer(peerId, conn)
	newSession := server.Session{
		Id:   sessionId(r.URL),
		Peer: peer,
	}

	sessonManager.Register(newSession)
	go peer.Listen()
	err = peer.Broadcast()
	if err != nil {
		log.Println(err)
	}
	sessonManager.Unregister(newSession)
}

func main() {
	sessonManager = server.NewSessionManager()
	http.HandleFunc("/", join)
	http.ListenAndServe("localhost:9988", nil)
}

func sessionId(url *url.URL) string {
	segments := strings.SplitAfterN(url.Path, "/", 2)
	return strings.Join(segments[1:], "/")

}
