package server

import (
	"fmt"
	"os"
)

//go:generate hel --type Conn --output mock_conn_test.go

type Conn interface {
	ReadMessage() (int, []byte, error)
	WriteMessage(int, []byte) error
}

type message struct {
	messageType int
	content     []byte
}

type Peer struct {
	Id             string
	Conn           Conn
	msg            chan message
	connectedPeers chan *Peer
}

func NewPeer(id string, conn Conn) *Peer {
	return &Peer{
		Id:             id,
		Conn:           conn,
		msg:            make(chan message),
		connectedPeers: make(chan *Peer),
	}
}

func (p *Peer) Listen() {
	for {
		mt, content, err := p.Conn.ReadMessage()
		if err != nil {
			//TODO: how are we going to deal with this error
			fmt.Fprintf(os.Stderr, "Listen: %s\n", err)
			continue
		}
		p.msg <- message{messageType: mt, content: content}
	}
}

func (p *Peer) Broadcast() {
	var peers []*Peer
	for {
		select {
		case msg := <-p.msg:
			for _, peer := range peers {
				err := peer.Conn.WriteMessage(msg.messageType, msg.content)
				if err != nil {
					//TODO: how are we going to deal with this error
					fmt.Fprintf(os.Stderr, "Broadcast to %s failed: %s\n", peer.Id, err)
					continue
				}
			}
		case peer := <-p.connectedPeers:
			peers = append(peers, peer)
		}
	}
}

func (p *Peer) Connect(peer *Peer) {
	p.connectedPeers <- peer
}
