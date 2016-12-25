package server

import (
	"fmt"
	"os"

	"github.com/gorilla/websocket"
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
	Id           string
	Conn         Conn
	msg          chan message
	connectPeers chan *Peer
}

func NewPeer(id string, conn Conn) *Peer {
	return &Peer{
		Id:           id,
		Conn:         conn,
		msg:          make(chan message),
		connectPeers: make(chan *Peer),
	}
}

func (p *Peer) Listen() {
	for {
		mt, content, err := p.Conn.ReadMessage()
		if e, ok := err.(*websocket.CloseError); ok {
			if e.Code == websocket.CloseAbnormalClosure {
				fmt.Fprintf(os.Stderr, "Listen is stopped: %s\n", err)
				return
			}
		}
		p.msg <- message{messageType: mt, content: content}
	}
}

func (p *Peer) Broadcast() error {
	var peers []*Peer
	for {
		select {
		case msg := <-p.msg:
			for _, peer := range peers {
				err := peer.Conn.WriteMessage(msg.messageType, msg.content)
				if err != nil {
					return err
				}
			}
		case peer := <-p.connectPeers:
			peers = append(peers, peer)
		}
	}
	return nil
}

func (p *Peer) Connect(peer *Peer) {
	p.connectPeers <- peer
}
