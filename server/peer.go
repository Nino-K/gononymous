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

func (p *Peer) broadcast(msg []byte) {
	p.connectedPeersMux.Lock()
	defer p.connectedPeersMux.Unlock()
	fmt.Println("connected peers are: ", p.connectedPeers)
	for _, peer := range p.connectedPeers {
		err := peer.Write(websocket.BinaryMessage, msg)
		if err != nil {
			//TODO: do something smart with this err
			// perhapes remove peer, retry, etc
			fmt.Println("something bad happened", err)
		}
	}
}

func (p *Peer) peerExist(peer *Peer) bool {
	for _, p := range p.connectedPeers {
		return p.Id == peer.Id
	}
	return false
}
