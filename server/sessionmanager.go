package server

import "fmt"

type Session struct {
	Id   string
	Peer *Peer
}

type SessionManager struct {
	register   chan Session
	unregister chan Session
}

func NewSessionManager() *SessionManager {
	sm := &SessionManager{
		register: make(chan Session),
	}
	go sm.run()
	return sm
}

func (s *SessionManager) Register(session Session) {
	s.register <- session
}

func (s *SessionManager) Unregister(session Session) {
	s.unregister <- session
}

func (s *SessionManager) run() {
	sessions := make(map[string][]*Peer)
	for {
		select {
		case session := <-s.register:
			fmt.Println(session.Id)
			peers, exist := sessions[session.Id]
			if !exist {
				sessions[session.Id] = []*Peer{session.Peer}
				continue
			}
			if peerExist(peers, session.Peer.Id) {
				continue
			}
			peers = append(peers, session.Peer)
			sessions[session.Id] = peers
			notifyPeers(peers)
			fmt.Println("notified all peers")
		}
		fmt.Printf("%#v \n", sessions)
	}
}

func notifyPeers(peers []*Peer) {
	for _, p := range peers {
		for _, eachPeer := range peers {
			p.Connect(eachPeer)
		}
	}
}

func peerExist(peers []*Peer, id string) bool {
	for _, p := range peers {
		if p.Id == id {
			return true
		}
	}
	return false
}
