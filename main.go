package main

import (
	"flag"
	"fmt"
	"net"
)

func main() {
	nodeNum := flag.Int("node", 1, "Node number (unique identifier for the node)")
	bootstrap := flag.String("bootstrap", "", "Bootstrap node address (e.g., 127.0.0.1:8001)")
	fanout := flag.Int("fanout", 2, "Number of random peers to gossip to")
	flag.Parse()

	basePort := 8000
	address := fmt.Sprintf("127.0.0.1:%d", basePort+*nodeNum)
	nodeID := fmt.Sprintf("node%d", *nodeNum)

	node := Node{
		ID: nodeID,
		Address: Peer{
			Address:     address,
			NetworkType: "tcp",
		},
		Peers: make(map[string]net.Addr),
	}

	go node.Start()

	if *bootstrap != "" {
		go node.JoinNetwork(*bootstrap)
		go node.Gossip(*fanout)
	}

	select {}
}
