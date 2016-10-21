package server_test

import (
	"net/url"
	"testing"

	"github.com/Nino-K/gononymous/server"
	"github.com/a8m/expect"
	"github.com/gorilla/websocket"
)

func TestSession_emptySessionHandleOrClientId(t *testing.T) {
	t.Log("It returns an error if given")
	{
		expect := expect.New(t)

		wsConn := &websocket.Conn{}
		session := server.NewSessionManager()
		t.Log("an empty sessionHandle")
		{
			err := session.Register("", "fakeClientID", wsConn)
			expect(err).Not.To.Be.Nil()
		}
		t.Log("an empty clientId")
		{
			err := session.Register("testSessionHandle", "", wsConn)
			expect(err).Not.To.Be.Nil()
		}

	}
}

func TestSession_successfullyRegistered(t *testing.T) {
	t.Log("It successfully registers a client when a valid URL passed")
	{
		expect := expect.New(t)

		sessionHandle := testURL().Path
		wsConn := &websocket.Conn{}
		session := server.NewSessionManager()
		session.Register(sessionHandle, "fakeClientID", wsConn)

		expectedClients := []server.Client{server.Client{Id: "fakeClientID", Conn: wsConn}}
		expect(session.Clients(sessionHandle)).To.Have.Len(1)
		expect(session.Clients(sessionHandle)).To.Equal(expectedClients)
	}
}

func TestSession_alreayRegisteredClient(t *testing.T) {
	t.Log("It should not add already registered client")
	{
		expect := expect.New(t)

		sessionHandle := testURL().Path
		wsConn := &websocket.Conn{}
		session := server.NewSessionManager()
		session.Register(sessionHandle, "fakeClientID", wsConn)
		// re-register
		session.Register(sessionHandle, "fakeClientID", wsConn)

		expectedClients := []server.Client{server.Client{Id: "fakeClientID", Conn: wsConn}}
		expect(session.Clients(sessionHandle)).To.Have.Len(1)
		expect(session.Clients(sessionHandle)).To.Equal(expectedClients)

	}
}

func TestSession_sameClientWithDifferentConn(t *testing.T) {
	t.Log("It updates client conn if it changes")
	{
		expect := expect.New(t)

		sessionHandle := testURL().Path
		wsConn := &websocket.Conn{}
		session := server.NewSessionManager()
		session.Register(sessionHandle, "fakeClientID", wsConn)
		// re-register with a different conn
		newConn := &websocket.Conn{}
		session.Register(sessionHandle, "fakeClientID", newConn)

		expectedClients := []server.Client{server.Client{Id: "fakeClientID", Conn: newConn}}
		expect(session.Clients(sessionHandle)).To.Have.Len(1)
		expect(session.Clients(sessionHandle)).To.Equal(expectedClients)
	}
}

func TestSession_multipleClients(t *testing.T) {
	t.Log("It registers multiple clients with same sessionHandle")
	{
		expect := expect.New(t)

		sessionHandle := testURL().Path
		clientConn1 := &websocket.Conn{}
		clientConn2 := &websocket.Conn{}
		session := server.NewSessionManager()
		session.Register(sessionHandle, "fakeClientID1", clientConn1)
		session.Register(sessionHandle, "fakeClientID2", clientConn2)

		expectedClients := []server.Client{
			server.Client{Id: "fakeClientID1", Conn: clientConn1},
			server.Client{Id: "fakeClientID2", Conn: clientConn2},
		}
		expect(session.Clients(sessionHandle)).To.Have.Len(2)
		expect(session.Clients(sessionHandle)).To.Equal(expectedClients)
	}
}

func testURL() url.URL {
	return url.URL{
		Scheme: "ws", Host: "testHost", Path: "some_session",
	}
}
