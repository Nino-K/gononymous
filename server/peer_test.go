package server_test

import (
	"errors"
	"testing"

	"github.com/Nino-K/gononymous/server"
	"github.com/a8m/expect"
	"github.com/gorilla/websocket"
)

func TestPeer_listenReadFromConn(t *testing.T) {
	t.Log("It reads from conn")
	{
		expect := expect.New(t)

		mockConn := newMockConn()
		mockConn.ReadMessageOutput.Ret0 <- websocket.BinaryMessage
		mockConn.ReadMessageOutput.Ret1 <- []byte("test stuff")
		mockConn.ReadMessageOutput.Ret2 <- nil

		peer := server.NewPeer("testId", mockConn)

		go func() {
			err := peer.Listen()
			expect(err).To.Be.Nil()
		}()

		called := <-mockConn.ReadMessageCalled
		expect(called).To.Equal(true)
	}
}

func TestPeer_listenErrorFromConn(t *testing.T) {
	t.Log("It returns immediatly if conn read errors")
	{
		expect := expect.New(t)

		mockConn := newMockConn()
		mockConn.ReadMessageOutput.Ret0 <- 0
		mockConn.ReadMessageOutput.Ret1 <- nil
		mockConn.ReadMessageOutput.Ret2 <- errors.New("something went wrong")

		peer := server.NewPeer("testId", mockConn)

		err := peer.Listen()
		expect(err).Not.To.Be.Nil()
	}
}

func TestPeer_write(t *testing.T) {
	t.Log("It should insert the messages into messages chan")
	{
		expect := expect.New(t)

		mockConn := newMockConn()
		peer := server.NewPeer("testId", mockConn)
		peer.Write([]byte("some test stuff"))
		expect(<-peer.Messages()).To.Equal([]byte("some test stuff"))
	}
}

func TestPeer_send(t *testing.T) {
	t.Log("It sends reads from msgs and writes to all connected Peers")
	{
		expect := expect.New(t)

		mockConnOne := newMockConn()
		mockConnTwo := newMockConn()

		peerOne := server.NewPeer("testId", mockConnOne)
		peerTwo := server.NewPeer("testId2", mockConnTwo)

		peerOne.Connect(peerTwo)
		peerOne.Write([]byte("some test stuff"))

		expect(<-mockConnTwo.WriteMessageCalled).To.Equal(true)
		expect(<-mockConnTwo.WriteMessageInput.MsgType).To.Equal(websocket.BinaryMessage)
		expect(<-mockConnTwo.WriteMessageInput.Data).To.Equal([]byte("some test stuff"))
	}
}

func TestPeer_writePeer(t *testing.T) {
	t.Log("It calls WriteMessage of conn")
	{
		expect := expect.New(t)

		mockConn := newMockConn()
		mockConn.WriteMessageOutput.Ret0 <- nil
		peer := server.NewPeer("testId", mockConn)
		peer.WritePeer(websocket.BinaryMessage, []byte("test Stuff"))

		expect(<-mockConn.WriteMessageCalled).To.Equal(true)
		expect(<-mockConn.WriteMessageInput.Data).To.Equal([]byte("test Stuff"))
		expect(<-mockConn.WriteMessageInput.MsgType).To.Equal(websocket.BinaryMessage)
	}
}
