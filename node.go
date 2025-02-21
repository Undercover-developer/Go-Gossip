package main

import (
	"fmt"
	"net"
	"time"
)

type Node struct {
	ID      string
	Address net.Addr
	Peers   map[string]net.Addr
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
