package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"maps"
	"math/rand/v2"
	"net"
	"strings"
	"time"
)

type Node struct {
	ID      string
	Address net.Addr
	Peers   map[string]net.Addr
}

type Peer struct {
	NetworkType string
	Address     string
}

func (p Peer) Network() string {
	return p.NetworkType
}

func (p Peer) String() string {
	return p.Address
}

func (n *Node) Start() {
	listener, err := net.Listen(n.Address.Network(), n.Address.String())
	if err != nil {
		fmt.Printf("Node %s: error starting server on %s: %v\n", n.ID, n.Address, err)
		return
	}

	fmt.Printf("Node %s: listening on %s\n", n.ID, n.Address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Node %s: error accepting connection: %v\n", n.ID, err)
			continue
		}

		go n.handleConnection(conn)
	}
}

func (n *Node) handleConnection(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	nBytes, err := conn.Read(buffer)

	if err != nil {
		fmt.Printf("Node %s: error reading data: %v\n", n.ID, err)
		return
	}

	message := string(buffer[:nBytes])
	fmt.Printf("Node %s: received message: %s\n", n.ID, message)

	//process join  request
	if strings.HasPrefix(message, "JOIN") {
		parts := strings.Split(message, " ")
		if len(parts) < 3 {
			newPeer := Peer{
				"tcp",
				parts[2],
			}
			ID := parts[1]
			if _, ok := n.Peers[ID]; !ok {
				n.Peers[ID] = newPeer
				fmt.Printf("Node %s added new peer %s \n", n.ID, newPeer.String())
			}
		}

		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		err := enc.Encode(n.Peers)
		if err != nil {
			fmt.Printf("Node %s: error occurred while encoding peer list: %v", n.ID, err)
		}
		prefix := []byte("JOIN")
		peerList := append(prefix, buf.Bytes()...)
		conn.Write(peerList)
		return
	}

}

func (n *Node) SendMessage(peer net.Addr, message string) error {
	conn, err := net.Dial(peer.Network(), peer.String())
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Write([]byte(message))
	if err != nil {
		return err
	}
	fmt.Printf("Node %s: sent message to %s: %s\n", n.ID, peer, message)
	return nil
}

func (n *Node) JoinNetwork(bootstrap string) {
	joinMessage := fmt.Sprintf("JOIN %s", n.Address)
	peer := Peer{
		"tcp",
		bootstrap,
	}
	conn, err := net.Dial(peer.Network(), peer.String())
	if err != nil {
		fmt.Printf("Node %s: error joining network via %s: %v\n", n.ID, bootstrap, err)
		return
	}
	defer conn.Close()
	_, err = conn.Write([]byte(joinMessage))

	buffer := make([]byte, 1024)
	nBytes, err := conn.Read(buffer)

	if err != nil {
		fmt.Printf("Node %s: error reading data: %v\n", n.ID, err)
		return
	}

	message := string(buffer[:nBytes])

	//handle peers list update
	if strings.HasPrefix(message, "PEERS") {
		parts := strings.Split(message, " ")
		if len(parts) < 2 {
			b := []byte(parts[1])
			reader := bytes.NewReader(b)
			dec := gob.NewDecoder(reader)
			var newPeerList map[string]net.Addr
			err := dec.Decode(&newPeerList)
			if err != nil {
				fmt.Printf("Node %s: error occurred while decoding list of peers: %v \n", n.ID, err)
			}

			maps.Copy(n.Peers, newPeerList)
		}
	}
}

func (n *Node) getGossipPeers(fanout int) []string {
	if fanout > len(n.Peers) {
		fanout = len(n.Peers)
	}

	peerIDs := make([]string, len(n.Peers))
	for id := range n.Peers {
		peerIDs = append(peerIDs, id)
	}

	rand.Shuffle(len(peerIDs), func(i, j int) {
		peerIDs[i], peerIDs[j] = peerIDs[j], peerIDs[i]
	})

	return peerIDs[:fanout]
}

func (n *Node) Gossip() {
	for {
		time.Sleep(3 * time.Second)

		if len(n.Peers) == 0 {
			continue
		}

		message := fmt.Sprintf("Gossip from node %s", n.ID)

		for _, peer := range n.Peers {
			err := n.SendMessage(peer, message)
			if err != nil {
				fmt.Printf("Node %s: error sending gossip to %s: %v\n", n.ID, peer, err)
			}
		}
	}
}
