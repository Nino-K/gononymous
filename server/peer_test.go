package server_test

import (
	"errors"
	"testing"

	"github.com/Nino-K/gononymous/server"
	"github.com/a8m/expect"
	"github.com/gorilla/websocket"
)

func TestPeer_listenReadFromConn(t *testing.T) {
	t.Log("It should read from conn")
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

		messages := peer.Messages()
		data := <-messages
		expect(data).To.Equal([]byte("test stuff"))
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

func TestPeer_signalConnectedPeers(t *testing.T) {
	t.Log("Signal notifies other connected Peers")
	{
		//TODO
		t.Skip("need to figure out how to do this signalling thing")
	}
}

//func TestPeer_sendWaitForAPeeer(t *testing.T) {
//	t.Log("It does not send until a new peer connected")
//	{
//		expect := expect.New(t)
//
//		mockConn := newMockConn()
//		perr := server.NewPeer("testId", mockConn)
//
//		err := peer.Send()
//		expect(err).Not.To.Be.Nil()
//	}
//}
