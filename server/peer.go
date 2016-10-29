package server

import (
	"fmt"

	"github.com/gorilla/websocket"
)

//go:generate hel --type Conn --output mock_conn_test.go

type Conn interface {
	WriteMessage(msgType int, data []byte) error
	ReadMessage() (int, []byte, error)
}

type Peer struct {
	Id             string
	Conn           Conn
	msgs           chan []byte
	connectedPeers []*Peer
}

func NewPeer(id string, conn Conn) *Peer {
	peer := &Peer{
		Id:   id,
		Conn: conn,
		msgs: make(chan []byte, 1000),
	}
	go peer.Send()
	return peer
}

func (p *Peer) Listen() error {
	for {
		_, b, err := p.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println(string(b))
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

func (p *Peer) Write(msg []byte) {
	p.msgs <- msg
}

func (p *Peer) Send() {
	for {
		select {
		case msg := <-p.msgs:
			p.broadcast(msg)
		}
	}
}

func (p *Peer) broadcast(msg []byte) {
	fmt.Println("connected peers are: ", p.connectedPeers)
	for _, peer := range p.connectedPeers {
		err := peer.Conn.WriteMessage(websocket.BinaryMessage, msg)
		if err != nil {
			//TODO: do something smart with this err
			// perhapes remove peer, retry, etc
			fmt.Println("something bad happened")
		}
	}
}
