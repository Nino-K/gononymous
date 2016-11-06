package server

import (
	"fmt"
	"sync"
)

type SessionManager struct {
	peers     map[string][]*Peer
	peersLock sync.Mutex
	signals   chan string
}

func NewSessionManager() *SessionManager {
	sm := &SessionManager{
		peers:   make(map[string][]*Peer),
		signals: make(chan string),
	}
	go sm.signal()
	return sm
}

// this is going to send a signal to peer.Signal chan also
func (s *SessionManager) Register(sessionHandle string, p *Peer) error {
	if sessionHandle == "" {
		return fmt.Errorf("Register: sessionHandle can not be empty")
	}
	s.peersLock.Lock()
	defer s.peersLock.Unlock()
	peers, exist := s.peers[sessionHandle]
	if !exist {
		s.peers[sessionHandle] = []*Peer{p}
		return nil
	}
	for i, peer := range peers {
		// if the peer found, update it's connection
		if p.Id == peer.Id {
			peers = append(peers[:i], peers[i+1:]...)
		}
		peers = append(peers, p)
	}
	s.peers[sessionHandle] = peers

	s.signals <- sessionHandle
	return nil
}

func (s *SessionManager) Unregister(sessionHandle, peerId string) error {
	if sessionHandle == "" {
		return fmt.Errorf("Unregistr: sessionHandle can not be empty")
	}
	if peerId == "" {
		return fmt.Errorf("Unregister: peerId can not be empty")
	}
	s.peersLock.Lock()
	defer s.peersLock.Unlock()
	peers, exist := s.peers[sessionHandle]
	if exist {
		for i, peer := range peers {
			if peer.Id == peerId {
				peers = append(peers[:i], peers[i+1:]...)
			}
		}
		s.peers[sessionHandle] = peers
	}
	return nil
}

func (s *SessionManager) Peers(sessionHandle string) []*Peer {
	s.peersLock.Lock()
	defer s.peersLock.Unlock()
	return s.peers[sessionHandle]
}

func (s *SessionManager) signal() {
	for {
		select {
		case sessionHandle := <-s.signals:
			peers, exist := s.peers[sessionHandle]
			if exist && len(peers) > 1 {
				for _, p := range peers {
					s.notify(p, peers)
				}
			}

		}
	}
}

func (s *SessionManager) notify(peer *Peer, peers []*Peer) {
	for _, p := range peers {
		if p.Id != peer.Id {
			p.Connect(peer)
		}
	}
}
