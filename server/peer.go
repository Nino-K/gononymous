package server

import "fmt"

//go:generate hel --type Conn --output mock_conn_test.go

type Conn interface {
	WriteMessage(msgType int, data []byte) error
	ReadMessage() (int, []byte, error)
}

// expose a signal chan to pass to sessionManager
// session manager's register will pass the newly registered peers down signal chan
type Peer struct {
	Id             string
	Conn           Conn
	msgs           chan []byte
	connectedPeers []*Peer
}

func NewPeer(id string, conn Conn) *Peer {
	return &Peer{
		Id:   id,
		Conn: conn,
		msgs: make(chan []byte),
	}
}

func (p *Peer) Listen() error {
	for {
		fmt.Println("connected peers are: ", p.connectedPeers)
		_, b, err := p.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println(string(b))
		//p.msgs <- b
	}
}

func (p *Peer) Messages() chan []byte {
	return p.msgs
}

func (p *Peer) Connect(peer *Peer) {
	if !p.peerExist(peer) {
		p.connectedPeers = append(p.connectedPeers, peer)
	}
}

func (p *Peer) peerExist(peer *Peer) bool {
	for _, p := range p.connectedPeers {
		return p.Id == peer.Id
	}
	return false
}

//TODO Send only upon other client connection
//func (p *peer) Send(msg []byte) {
//	for {
//		select {
//		case peer := <-p.peers:
//			peer.Send(msg)
//		}
//	}
//}
