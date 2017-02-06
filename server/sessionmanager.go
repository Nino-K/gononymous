package server

import "log"

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
		register:   make(chan Session),
		unregister: make(chan Session),
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
			peers, exist := sessions[session.Id]
			if !exist {
				sessions[session.Id] = []*Peer{session.Peer}
				log.Printf("%s joined, total peers %+v \n", session.Peer.Id, sessions)
				continue
			}
			if peerExist(peers, session.Peer.Id) {
				continue
			}
			notifyPeers(peers, session.Peer)
			peers = append(peers, session.Peer)
			sessions[session.Id] = peers
			log.Printf("%s joined, total peers %+v \n", session.Peer.Id, sessions)
		case session := <-s.unregister:
			peers, exist := sessions[session.Id]
			if exist {
				for i, p := range peers {
					if p.Id == session.Peer.Id {
						peers = append(peers[:i], peers[i+1:]...)
						p.Conn.Close()
					}
				}
				sessions[session.Id] = peers
				if len(peers) == 0 {
					delete(sessions, session.Id)
				}
			}
			log.Printf("%s left, remaining peers %+v \n", session.Peer.Id, sessions)
		}
	}
}

func notifyPeers(peers []*Peer, peer *Peer) {
	for _, p := range peers {
		// let others know about the new peer
		p.Connect(peer)
		// let the new peer know about others
		peer.Connect(p)
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
