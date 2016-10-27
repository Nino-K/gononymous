package server_test

import (
	"net/url"
	"testing"

	"github.com/Nino-K/gononymous/server"
	"github.com/a8m/expect"
	"github.com/gorilla/websocket"
)

func TestSession_emptySessionHandle(t *testing.T) {
	t.Log("It returns an error if given an empty sessionHandle")
	{
		expect := expect.New(t)

		wsConn := &websocket.Conn{}
		peer := server.NewPeer("fakepeerID", wsConn)
		session := server.NewSessionManager()
		err := session.Register("", peer)
		expect(err).Not.To.Be.Nil()

	}
}

func TestSession_successfullyRegistered(t *testing.T) {
	t.Log("It successfully registers a peer when a valid URL passed")
	{
		expect := expect.New(t)

		sessionHandle := testURL().Path
		wsConn := &websocket.Conn{}
		peer := server.NewPeer("fakepeerID", wsConn)
		session := server.NewSessionManager()
		session.Register(sessionHandle, peer)

		expectedpeers := []*server.Peer{peer}
		expect(session.Peers(sessionHandle)).To.Have.Len(1)
		expect(session.Peers(sessionHandle)).To.Equal(expectedpeers)
	}
}

func TestSession_alreayRegisteredpeer(t *testing.T) {
	t.Log("It should not add already registered peer")
	{
		expect := expect.New(t)

		sessionHandle := testURL().Path
		wsConn := &websocket.Conn{}
		peer := server.NewPeer("fakepeerID", wsConn)
		session := server.NewSessionManager()
		session.Register(sessionHandle, peer)
		// re-register
		session.Register(sessionHandle, peer)

		expectedpeers := []*server.Peer{peer}
		expect(session.Peers(sessionHandle)).To.Have.Len(1)
		expect(session.Peers(sessionHandle)).To.Equal(expectedpeers)

	}
}

func TestSession_samePeerWithDifferentConn(t *testing.T) {
	t.Log("It updates peer conn if it changes")
	{
		expect := expect.New(t)

		sessionHandle := testURL().Path
		wsConn := &websocket.Conn{}
		session := server.NewSessionManager()
		oldPeer := server.NewPeer("fakepeerID", wsConn)
		session.Register(sessionHandle, oldPeer)
		// re-register with a different conn
		newConn := &websocket.Conn{}
		newPeer := server.NewPeer("fakepeerID", newConn)
		session.Register(sessionHandle, newPeer)

		expectedpeers := []*server.Peer{newPeer}
		expect(session.Peers(sessionHandle)).To.Have.Len(1)
		expect(session.Peers(sessionHandle)).To.Equal(expectedpeers)
	}
}

func TestSession_multiplePeers(t *testing.T) {
	t.Log("It registers multiple peers with same sessionHandle")
	{
		expect := expect.New(t)

		sessionHandle := testURL().Path
		peerConn1 := &websocket.Conn{}
		peerOne := server.NewPeer("fakepeerID1", peerConn1)
		peerConn2 := &websocket.Conn{}
		peerTwo := server.NewPeer("fakepeerID2", peerConn2)
		session := server.NewSessionManager()
		session.Register(sessionHandle, peerOne)
		session.Register(sessionHandle, peerTwo)

		expectedpeers := []*server.Peer{peerOne, peerTwo}
		expect(session.Peers(sessionHandle)).To.Have.Len(2)
		expect(session.Peers(sessionHandle)).To.Equal(expectedpeers)
	}
}

func TestSession_registerSignalsOtherPeers(t *testing.T) {
	// TODO
	t.Skip("need to do some mocking for this one")
}

func testURL() url.URL {
	return url.URL{
		Scheme: "ws", Host: "testHost", Path: "some_session",
	}
}
