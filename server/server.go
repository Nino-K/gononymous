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

	err = s.Register(sessionHandle, "change_me", conn)
	if err != nil {
		fmt.Println("error")
		panic(err)
	}

	fmt.Println(s.Clients(sessionHandle))
	fmt.Println(len(s.Clients(sessionHandle)))
	conn.WriteMessage(websocket.BinaryMessage, []byte("called Home"+sessionHandle))

}

func sessionHandle(url *url.URL) string {
	segments := strings.SplitAfterN(url.Path, "/", 2)
	return strings.Join(segments[1:], "/")

}
