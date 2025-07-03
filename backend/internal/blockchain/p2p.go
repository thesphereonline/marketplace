package blockchain

import (
	"bufio"
	"fmt"
	"net"
)

type Peer struct {
	Conn net.Conn
	IP   string
}

type Node struct {
	Peers []Peer
}

func (n *Node) Start(address string) {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	fmt.Println("✅ Node listening on", address)

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		peer := Peer{Conn: conn, IP: conn.RemoteAddr().String()}
		n.Peers = append(n.Peers, peer)

		go n.handleConnection(peer)
	}
}

func (n *Node) handleConnection(peer Peer) {
	scanner := bufio.NewScanner(peer.Conn)
	for scanner.Scan() {
		msg := scanner.Text()
		fmt.Printf("📨 [%s] %s\n", peer.IP, msg)
		// Broadcast message to all peers
		n.Gossip(msg, peer)
	}
}

func (n *Node) Gossip(msg string, origin Peer) {
	for _, p := range n.Peers {
		if p.IP != origin.IP {
			fmt.Fprintf(p.Conn, "%s\n", msg)
		}
	}
}
