package handler

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/Nino-K/gononymous/server"
	"github.com/gorilla/websocket"
)

//go:generate hel --type Upgrader --output mock_upgrader_test.go

type Upgrader interface {
	Upgrade(http.ResponseWriter, *http.Request, http.Header) (*websocket.Conn, error)
}

const CLIENTID = "CLIENT_ID"

var (
	peerIdErr    = errors.New("CLIENT_ID header must be provided")
	sessionIDErr = errors.New("sessionId must be provided")
)

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
	peerId := r.Header.Get(CLIENTID)
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
	sessionId := sessionId(r.URL)
	if sessionId == "" {
		http.Error(w, sessionIDErr.Error(), http.StatusBadRequest)
		log.Printf("Join SessionId: %s \n", sessionIDErr)
		return
	}
	newSession := server.Session{
		Id:   sessionId,
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
