package server

import (
	"errors"
	"io"
	"sync"

	"github.com/Nino-K/gononymous/server/mock"
	"github.com/gorilla/websocket"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Peer", func() {
	Describe("Listen", func() {
		It("places the incoming messages into message chan", func() {
			fakeCon := mock.NewMockConn()
			fakeCon.ReadMessageOutput.Ret0 <- websocket.BinaryMessage
			fakeCon.ReadMessageOutput.Ret1 <- []byte("test content")
			fakeCon.ReadMessageOutput.Ret2 <- nil

			peer := NewPeer("testId", fakeCon)
			go peer.Listen()

			var msg message
			Eventually(peer.msg).Should(Receive(&msg))
			Expect(msg.content).To(Equal([]byte("test content")))
			Expect(msg.messageType).To(Equal(websocket.BinaryMessage))
		})

		Context("error", func() {
			It("returns when close err received from connection", func() {
				fakeCon := mock.NewMockConn()
				fakeCon.ReadMessageOutput.Ret0 <- websocket.CloseMessage
				fakeCon.ReadMessageOutput.Ret1 <- nil

				closeErr := &websocket.CloseError{
					Code: websocket.CloseAbnormalClosure,
					Text: io.ErrUnexpectedEOF.Error(),
				}
				fakeCon.ReadMessageOutput.Ret2 <- closeErr
				peer := NewPeer("testId", fakeCon)
				go peer.Listen()
				Consistently(peer.msg).ShouldNot(Receive())
			})

			It("does not ignore the next message if error is not close error", func() {
				fakeCon := mock.NewMockConn()
				fakeCon.ReadMessageOutput.Ret0 <- websocket.BinaryMessage
				fakeCon.ReadMessageOutput.Ret1 <- []byte("test content")
				fakeCon.ReadMessageOutput.Ret2 <- errors.New("something went wrong")

				peer := NewPeer("testId", fakeCon)
				go peer.Listen()

				var msg message
				Eventually(peer.msg).Should(Receive(&msg))
				Expect(msg.content).To(Equal([]byte("test content")))
				Expect(msg.messageType).To(Equal(websocket.BinaryMessage))
			})

			It("sends a stop sigal when close error received", func() {
				fakeCon := mock.NewMockConn()
				fakeCon.ReadMessageOutput.Ret0 <- websocket.CloseMessage
				fakeCon.ReadMessageOutput.Ret1 <- nil

				closeErr := &websocket.CloseError{
					Code: websocket.CloseAbnormalClosure,
					Text: io.ErrUnexpectedEOF.Error(),
				}
				fakeCon.ReadMessageOutput.Ret2 <- closeErr
				peer := NewPeer("testId", fakeCon)
				go peer.Listen()
				Consistently(peer.msg).ShouldNot(Receive())
				Eventually(peer.stop).Should(Receive())
			})
			It("closes it's own connection when close error received", func() {
				fakeCon := mock.NewMockConn()
				fakeCon.ReadMessageOutput.Ret0 <- websocket.CloseMessage
				fakeCon.ReadMessageOutput.Ret1 <- nil

				closeErr := &websocket.CloseError{
					Code: websocket.CloseAbnormalClosure,
					Text: io.ErrUnexpectedEOF.Error(),
				}
				fakeCon.ReadMessageOutput.Ret2 <- closeErr
				peer := NewPeer("testId", fakeCon)
				go peer.Listen()
				Consistently(peer.msg).ShouldNot(Receive())
				Eventually(peer.stop).Should(Receive())
				Eventually(fakeCon.CloseCalled).Should(Receive())
			})

		})

		Describe("Broadcast", func() {
			It("does not write messages to it's own Conn", func() {
				fakeCon := mock.NewMockConn()
				peer := NewPeer("testId", fakeCon)

				go peer.Broadcast()
				peer.msg <- message{
					messageType: websocket.BinaryMessage,
					content:     []byte("test message"),
				}

				Consistently(fakeCon.WriteMessageInput.Arg0).ShouldNot(Receive())
				Consistently(fakeCon.WriteMessageInput.Arg1).ShouldNot(Receive())
			})

			It("writes messages to all Connected Peers", func() {
				fakeCon := mock.NewMockConn()
				peer := NewPeer("testId", fakeCon)

				go peer.Broadcast()

				fakeCon2 := mock.NewMockConn()
				peer2 := NewPeer("testId2", fakeCon2)
				peer.Connect(peer2)

				peer.msg <- message{
					messageType: websocket.BinaryMessage,
					content:     []byte("test message"),
				}

				var msgType int
				Eventually(fakeCon2.WriteMessageInput.Arg0).Should(Receive(&msgType))
				Expect(msgType).To(Equal(websocket.BinaryMessage))

				var msgContent []byte
				Eventually(fakeCon2.WriteMessageInput.Arg1).Should(Receive(&msgContent))
				Expect(msgContent).To(Equal([]byte("test message")))
			})

			Context("error", func() {
				It("returns an error when stop signal received", func() {
					fakeCon := mock.NewMockConn()
					peer := NewPeer("testId", fakeCon)

					go func() {
						peer.stop <- struct{}{}
					}()
					err := peer.Broadcast()
					Expect(err).To(HaveOccurred())
					expectedErr := &PeerStopErr{PeerId: "testId"}
					Expect(err).To(Equal(expectedErr))
				})

				It("removes the peer when WriteMessage returns an error", func() {
					fakeCon1 := mock.NewMockConn()
					peer1 := NewPeer("testId1", fakeCon1)

					fakeCon2 := mock.NewMockConn()
					peer2 := NewPeer("testId2", fakeCon2)

					fakeCon3 := mock.NewMockConn()
					peer3 := NewPeer("testId3", fakeCon3)

					go peer1.Broadcast()
					peer1.Connect(peer2)
					peer1.Connect(peer3)

					fakeCon2.WriteMessageOutput.Ret0 <- errors.New("something bad happened")
					peer1.msg <- message{
						messageType: websocket.BinaryMessage,
						content:     []byte("test message"),
					}

					Eventually(fakeCon2.WriteMessageCalled).Should(Receive())
					Eventually(fakeCon3.WriteMessageCalled).Should(Receive())
					fakeCon3.WriteMessageOutput.Ret0 <- nil

					peer1.msg <- message{
						messageType: websocket.BinaryMessage,
						content:     []byte("test message2"),
					}

					Eventually(fakeCon2.WriteMessageCalled).ShouldNot(Receive())
					Eventually(fakeCon3.WriteMessageCalled).Should(Receive())
				})

				It("returns PeerStopErr when peer stops listening", func() {
					fakeCon := mock.NewMockConn()
					peer := NewPeer("testId", fakeCon)

					var err error
					var wg sync.WaitGroup
					wg.Add(1)
					go func() {
						defer wg.Done()
						err = peer.Broadcast()
					}()

					peer.stop <- struct{}{}
					wg.Wait()
					Expect(err).NotTo(BeNil())
					Expect(err).To(Equal(&PeerStopErr{PeerId: "testId"}))
				})
			})
			Describe("Connect", func() {
				It("receives a peer", func() {
					fakeCon1 := mock.NewMockConn()
					peer1 := NewPeer("testId1", fakeCon1)
					fakeCon2 := mock.NewMockConn()
					peer2 := NewPeer("testId2", fakeCon2)

					go peer1.Connect(peer2)
					var expectedPeer *Peer
					Eventually(peer1.connectPeers).Should(Receive(&expectedPeer))
					Expect(expectedPeer).To(Equal(peer2))
				})
			})
		})
	})
})
