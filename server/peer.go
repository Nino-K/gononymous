package server

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

//go:generate hel --type Conn --output mock_conn_test.go

type Conn interface {
	WriteMessage(msgType int, data []byte) error
	ReadMessage() (int, []byte, error)
}

type Peer struct {
	Id                string
	conn              Conn
	msgs              chan []byte
	connectedPeers    []*Peer
	connectedPeersMux sync.Mutex
}

func NewPeer(id string, conn Conn) *Peer {
	peer := &Peer{
		Id:   id,
		conn: conn,
		msgs: make(chan []byte, 1000),
	}
	go peer.send()
	return peer
}

func (p *Peer) Listen() error {
	var b []byte
	var err error
	for {
		_, b, err = p.conn.ReadMessage()
		if err != nil {
			break
		}
		//fmt.Println(string(b))
		p.msgs <- b
	}
	return err
}

func (p *Peer) Connect(peer *Peer) {
	p.connectedPeersMux.Lock()
	defer p.connectedPeersMux.Unlock()
	if !p.peerExist(peer) {
		p.connectedPeers = append(p.connectedPeers, peer)
	}
}

//TODO: this thing is behaving so odd
//expose a method that returns p.connected peers
//and test this
func (p *Peer) Disconnect(peerId string) {
	fmt.Println("going to remove", peerId, "own peer id is", p.Id)
	p.connectedPeersMux.Lock()
	defer p.connectedPeersMux.Unlock()
	for i := len(p.connectedPeers) - 1; i >= 0; i-- {
		if p.connectedPeers[i].Id == peerId {
			p.connectedPeers = append(p.connectedPeers[:i], p.connectedPeers[i+1:]...)
		}
	}
}

func (p *Peer) Write(msgType int, msg []byte) error {
	return p.conn.WriteMessage(msgType, msg)
}

func (p *Peer) send() {
	for {
		select {
		case msg := <-p.msgs:
			p.broadcast(msg)
		}
	}
}

// update the connected peers
func (p *Peer) broadcast(msg []byte) {
	p.connectedPeersMux.Lock()
	defer p.connectedPeersMux.Unlock()
	fmt.Println("connected peers are: ", p.connectedPeers)
	fmt.Println("connected peers total len: ", len(p.connectedPeers))
	for _, peer := range p.connectedPeers {
		err := peer.Write(websocket.BinaryMessage, msg)
		if err != nil {
			fmt.Println("something bad happened", err)

		}
	}
}

func (p *Peer) Peers() []*Peer {
	p.connectedPeersMux.Lock()
	defer p.connectedPeersMux.Unlock()
	return p.connectedPeers
}

func (p *Peer) peerExist(peer *Peer) bool {
	for _, p := range p.connectedPeers {
		return p.Id == peer.Id
	}
	return false
}
