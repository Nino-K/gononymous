package server

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Id   string
	Conn *websocket.Conn
}

type SessionManager struct {
	clients     map[string][]Client
	clientsLock sync.Mutex
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		clients: make(map[string][]Client),
	}
}

func (s *SessionManager) Register(sessionHandle string, clientId string, conn *websocket.Conn) error {
	if sessionHandle == "" {
		return fmt.Errorf("Register: sessionHandle can not be empty")
	}
	if clientId == "" {
		return fmt.Errorf("Register: clientId can not be empty")
	}
	s.clientsLock.Lock()
	defer s.clientsLock.Unlock()
	clients, exist := s.clients[sessionHandle]
	if !exist {
		s.clients[sessionHandle] = []Client{Client{Id: clientId, Conn: conn}}
		return nil
	}
	for _, client := range clients {
		if clientId != client.Id {
			clients = append(clients, Client{Id: clientId, Conn: conn})
		}
	}
	s.clients[sessionHandle] = clients
	return nil
}

func (s *SessionManager) Clients(sessionHandle string) []Client {
	s.clientsLock.Lock()
	defer s.clientsLock.Unlock()
	return s.clients[sessionHandle]
}
