// main.go
package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	// created two nodes here, maybe
	node1 := NewNode("localhost:3000")
	node2 := NewNode("localhost:3001")

	// peers
	node1.AddPeer(node2.Address)
	node2.AddPeer(node1.Address)

	// Start servers
	go node1.StartServer()
	go node2.StartServer()

	
	time.Sleep(2 * time.Second)

	
	wallet1 := NewWallet()
	wallet2 := NewWallet()

	fmt.Printf("Wallet 1 Address: %s\n", wallet1.GetAddress())
	fmt.Printf("Wallet 2 Address: %s\n", wallet2.GetAddress())

	
	tx := NewTransaction(wallet1.GetAddress(), wallet2.GetAddress(), 10, wallet1.PrivateKey)
	if !tx.Verify() {
		log.Println("Transaction verification failed!")
		return
	}

	
	message := fmt.Sprintf("Transaction: %s -> %s : %d", tx.Sender, tx.Receiver, tx.Amount)

	
	node1.Broadcast(message)

	
	select {}
}
