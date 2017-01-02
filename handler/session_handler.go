package handler

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/Nino-K/gononymous/server"
)

//go:generate hel --type Upgrader --output mock_upgrader_test.go

type Upgrader interface {
	Upgrade(http.ResponseWriter, *http.Request, http.Header) (server.Conn, error)
}

var peerIdErr = errors.New("CLIENT_ID header must be provided")

type SessionHandler struct {
	SessionManager *server.SessionManager
	Upgrader       Upgrader
}

func NewSessionHandler(sessionManager *server.SessionManager, upgrader Upgrader) *SessionHandler {
	return &SessionHandler{
		SessionManager: sessionManager,
		Upgrader:       upgrader,
	}
}

func (s *SessionHandler) Join(w http.ResponseWriter, r *http.Request) {
	peerId := r.Header.Get("CLIENT_ID")
	if peerId == "" {
		http.Error(w, peerIdErr.Error(), http.StatusBadRequest)
		log.Printf("Join PeerId: %s \n", peerIdErr)
		return
	}
	conn, err := s.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Join upgrader: %s \n", err)
		return
	}
	peer := server.NewPeer(peerId, conn)
	newSession := server.Session{
		Id:   sessionId(r.URL),
		Peer: peer,
	}

	s.SessionManager.Register(newSession)
	go peer.Listen()
	err = peer.Broadcast()
	if err != nil {
		log.Println(err)
	}
	s.SessionManager.Unregister(newSession)
}

func sessionId(url *url.URL) string {
	segments := strings.SplitAfterN(url.Path, "/", 2)
	return strings.Join(segments[1:], "/")

}
