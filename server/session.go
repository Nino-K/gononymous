package server

import (
	"fmt"
	"sync"
)

const (
	connect = iota + 1
	disconnect
)

type signal struct {
	code          int
	sessionHandle string
	clientId      string
}

type SessionManager struct {
	peers      map[string][]*Peer
	peersMux   sync.RWMutex
	register   chan signal
	unregister chan signal
}

func NewSessionManager() *SessionManager {
	sm := &SessionManager{
		peers:      make(map[string][]*Peer),
		register:   make(chan signal),
		unregister: make(chan signal),
	}
	go sm.signal()
	return sm
}

// this is going to send a signal to peer.Signal chan also
func (s *SessionManager) Register(sessionHandle string, p *Peer) error {
	if sessionHandle == "" {
		return fmt.Errorf("Register: sessionHandle can not be empty")
	}
	s.peersMux.Lock()
	defer s.peersMux.Unlock()
	peers, exist := s.peers[sessionHandle]
	if !exist {
		s.peers[sessionHandle] = []*Peer{p}
		return nil
	}
	for i := len(peers) - 1; i >= 0; i-- {
		if p.Id == peers[i].Id {
			peers = append(peers[:i], peers[i+1:]...)
		}
	}
	peers = append(peers, p)
	s.peers[sessionHandle] = peers
	fmt.Println("about to send signal")
	s.register <- signal{
		code:          connect,
		sessionHandle: sessionHandle,
		clientId:      p.Id,
	}
	fmt.Printf("total %d registered for session %s", len(peers), sessionHandle)
	return nil
}

func (s *SessionManager) Unregister(sessionHandle, peerId string) error {
	if sessionHandle == "" {
		return fmt.Errorf("Unregistr: sessionHandle can not be empty")
	}
	if peerId == "" {
		return fmt.Errorf("Unregister: peerId can not be empty")
	}
	s.peersMux.Lock()
	defer s.peersMux.Unlock()
	peers, exist := s.peers[sessionHandle]
	if exist {
		for i, peer := range peers {
			if peer.Id == peerId {
				peers = append(peers[:i], peers[i+1:]...)
			}
		}
		s.peers[sessionHandle] = peers
		// TODO: do we need to signal if !exist
		s.unregister <- signal{
			code:          disconnect,
			sessionHandle: sessionHandle,
			clientId:      peerId,
		}
	}
	return nil
}

func (s *SessionManager) Peers(sessionHandle string) []*Peer {
	s.peersMux.RLock()
	defer s.peersMux.RUnlock()
	return s.peers[sessionHandle]
}

func (s *SessionManager) signal() {
	for {
		select {
		case connectSig := <-s.register:
			s.notify(connectSig)
		case disconnectSig := <-s.unregister:
			s.notify(disconnectSig)
		default:
		}
	}
}

//TODO: refactor this thing, looks very scary :-|
func (s *SessionManager) notify(sig signal) {
	s.peersMux.RLock()
	defer s.peersMux.RUnlock()
	peers, exist := s.peers[sig.sessionHandle]
	// there is no point to notify yourself about yourself
	if sig.code == connect && len(peers) < 2 {
		return
	}
	if exist {
		for _, peer := range peers {
			if sig.code == disconnect {
				s.signalDisconnect(sig.clientId, peers)
			} else {
				s.signalConnect(peer, peers)
			}
		}
	}
}

func (s *SessionManager) signalConnect(peer *Peer, peers []*Peer) {
	for _, p := range peers {
		if p.Id != peer.Id {
			p.Connect(peer)
		}
	}
}

func (s *SessionManager) signalDisconnect(peerId string, peers []*Peer) {
	for _, p := range peers {
		p.Disconnect(peerId)
	}
}
