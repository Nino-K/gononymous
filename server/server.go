package server

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
)

type Server struct {
	*SessionManager
	upgrader websocket.Upgrader
}

func New() *Server {
	return &Server{
		NewSessionManager(),
		websocket.Upgrader{},
	}
}

func (s *Server) Home(w http.ResponseWriter, r *http.Request) {
	sessionHandle := sessionHandle(r.URL)
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}

	clientId := r.Header.Get("CLIENT_ID")
	fmt.Println(clientId, "CONNECTED")
	peer := NewPeer(clientId, conn)
	err = s.Register(sessionHandle, peer)
	if err != nil {
		fmt.Println("error")
		panic(err)
	}

	err = peer.Listen()
	if err != nil {
		fmt.Errorf("Home Handler: %v", err)
		return
	}
}

func sessionHandle(url *url.URL) string {
	segments := strings.SplitAfterN(url.Path, "/", 2)
	return strings.Join(segments[1:], "/")

}
