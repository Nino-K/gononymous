package server

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sessionmanager", func() {
	Describe("peerExist", func() {
		It("returns true if finds a peer by an ID", func() {
			peers := []*Peer{
				NewPeer("id1", nil),
				NewPeer("id2", nil),
				NewPeer("id3", nil),
			}
			Expect(peerExist(peers, "id2")).To(BeTrue())
		})
	})
	Describe("run", func() {
		It("connects the newly added peer to all existing peers", func() {
			sm := NewSessionManager()

			session1 := Session{
				Id:   "session1",
				Peer: NewPeer("client1", nil),
			}
			sm.Register(session1)

			session2 := Session{
				Id:   "session1",
				Peer: NewPeer("client1", nil),
			}
			sm.Register(session2)

		})
	})
})
