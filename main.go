package main

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/Nino-K/gononymous/server"
)

var sessonManager *server.SessionManager
var i int

func join(w http.ResponseWriter, r *http.Request) {
	peerId := r.Header.Get("CLIENT_ID")
	newSession := server.Session{
		Id: sessionId(r.URL),
		Peer: server.Peer{
			Id: peerId,
		},
	}
	sessonManager.Register(newSession)
	i++
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
