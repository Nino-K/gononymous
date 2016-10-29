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
	//TODO: go peer read and go peer write

	//conn.WriteMessage(websocket.BinaryMessage, []byte("called Home"+sessionHandle))

	go func() {
		i := 0
		for {
			peer.Write([]byte(fmt.Sprintf("sending %d from %s", i, clientId)))
			i++
		}
	}()
	err = peer.Listen()
	if err != nil {
		panic(err)
	}
}

func sessionHandle(url *url.URL) string {
	segments := strings.SplitAfterN(url.Path, "/", 2)
	return strings.Join(segments[1:], "/")

}
