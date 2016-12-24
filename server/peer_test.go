package server

import (
	"errors"

	"github.com/gorilla/websocket"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Peer", func() {
	Describe("Listen", func() {
		It("places the incoming messages into message chan", func() {
			fakeCon := newMockConn()
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
			It("does not process the message when error received", func() {
				fakeCon := newMockConn()
				fakeCon.ReadMessageOutput.Ret0 <- websocket.BinaryMessage
				fakeCon.ReadMessageOutput.Ret1 <- []byte("test content")
				fakeCon.ReadMessageOutput.Ret2 <- errors.New("something went wrong")

				peer := NewPeer("testId", fakeCon)
				go peer.Listen()

				Consistently(peer.msg).ShouldNot(Receive())
			})

			It("does not ignore the next message when error received", func() {
				fakeCon := newMockConn()
				fakeCon.ReadMessageOutput.Ret0 <- websocket.BinaryMessage
				fakeCon.ReadMessageOutput.Ret1 <- []byte("test content")
				fakeCon.ReadMessageOutput.Ret2 <- errors.New("something went wrong")

				peer := NewPeer("testId", fakeCon)
				go peer.Listen()

				Consistently(peer.msg).ShouldNot(Receive())

				fakeCon.ReadMessageOutput.Ret0 <- websocket.BinaryMessage
				fakeCon.ReadMessageOutput.Ret1 <- []byte("test content")
				fakeCon.ReadMessageOutput.Ret2 <- nil
				var msg message
				Eventually(peer.msg).Should(Receive(&msg))
				Expect(msg.content).To(Equal([]byte("test content")))
				Expect(msg.messageType).To(Equal(websocket.BinaryMessage))
			})

		})

		Describe("Broadcast", func() {

			It("does not write messages to it's own Conn", func() {
				fakeCon := newMockConn()
				peer := NewPeer("testId", fakeCon)

				go peer.Broadcast()

				peer.msg <- message{
					messageType: websocket.BinaryMessage,
					content:     []byte("test message"),
				}

				Eventually(fakeCon.WriteMessageInput.Arg0).ShouldNot(Receive())
				Eventually(fakeCon.WriteMessageInput.Arg1).ShouldNot(Receive())
			})

			It("writes messages to all Connected Peers", func() {
				fakeCon := newMockConn()
				peer := NewPeer("testId", fakeCon)

				go peer.Broadcast()

				fakeCon2 := newMockConn()
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
				It("does not stop reading from msg chan if WriteMessage returns error", func() {
					fakeCon1 := newMockConn()
					peer1 := NewPeer("testId1", fakeCon1)

					fakeCon2 := newMockConn()
					peer2 := NewPeer("testId2", fakeCon2)

					go peer1.Broadcast()
					peer1.Connect(peer2)

					fakeCon2.WriteMessageOutput.Ret0 <- errors.New("something went wrong while writing")

					peer1.msg <- message{
						messageType: websocket.BinaryMessage,
						content:     []byte("test message one"),
					}

					Expect(<-fakeCon2.WriteMessageCalled).To(BeTrue())

					//drain the 1st message
					<-fakeCon2.WriteMessageInput.Arg0
					<-fakeCon2.WriteMessageInput.Arg1

					go func() {
						peer1.msg <- message{
							messageType: websocket.BinaryMessage,
							content:     []byte("test message"),
						}
					}()

					var msgType int
					Eventually(fakeCon2.WriteMessageInput.Arg0).Should(Receive(&msgType))
					Expect(msgType).To(Equal(websocket.BinaryMessage))

					var msgContent []byte
					Eventually(fakeCon2.WriteMessageInput.Arg1).Should(Receive(&msgContent))
					Expect(msgContent).To(Equal([]byte("test message")))
				})
			})
		})
	})
})
