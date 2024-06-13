
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type Node struct {
	Address    string
	Blockchain *Blockchain
	Peers      []string
}

func NewNode(address string) *Node {
	blockchain := NewBlockchain()
	return &Node{Address: address, Blockchain: blockchain, Peers: make([]string, 0)}
}

func (n *Node) AddPeer(peerAddress string) {
	n.Peers = append(n.Peers, peerAddress)
}

func (n *Node) Broadcast(data string) {
	for _, peer := range n.Peers {
		fmt.Printf("Sending to %s: %s\n", peer, data)
		n.SendMessage(peer, data)
	}
}

func (n *Node) Receive(data string) {
	fmt.Printf("Node %s received: %s\n", n.Address, data)
	n.Blockchain.AddBlock(data)
}

func (n *Node) StartServer() {
	ln, err := net.Listen("tcp", n.Address)
	if err != nil {
		log.Fatalf("Error starting server: %s", err.Error())
	}
	defer ln.Close()
	log.Printf("Node running on %s", n.Address)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %s", err.Error())
			continue
		}
		go n.handleConnection(conn)
	}
}

func (n *Node) handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		message = strings.TrimSpace(message)
		n.Receive(message)
	}
}

func (n *Node) SendMessage(peerAddress, message string) {
	conn, err := net.Dial("tcp", peerAddress)
	if err != nil {
		log.Printf("Error connecting to peer %s: %s", peerAddress, err.Error())
		return
	}
	defer conn.Close()
	fmt.Fprintf(conn, "%s\n", message)
}
