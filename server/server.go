package server

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

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
	defer func() {
		conn.WriteControl(websocket.CloseAbnormalClosure, nil, time.Time{})
	}()

	clientId := r.Header.Get("CLIENT_ID")
	fmt.Println(clientId, "CONNECTED")
	peer := NewPeer(clientId, conn)
	err = s.Register(sessionHandle, peer)
	if err != nil {
		panic(err)
	}

	err = peer.Listen()
	if err != nil {
		s.Unregister(sessionHandle, clientId)
		// remove all cached instances of this peer, signal other peers
		//TODO: add s.Unregister and signal peers if err is EOF
		fmt.Println("Home Handler: %v", err)
		return
	}
}

func sessionHandle(url *url.URL) string {
	segments := strings.SplitAfterN(url.Path, "/", 2)
	return strings.Join(segments[1:], "/")

}
