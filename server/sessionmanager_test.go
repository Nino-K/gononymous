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
})
