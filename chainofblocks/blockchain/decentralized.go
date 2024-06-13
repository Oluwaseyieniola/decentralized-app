
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type Block struct {
	Timestamp     int64
	Data          string
	PrevBlockHash string
	Hash          string
}

type Blockchain struct {
	Blocks []*Block
}

func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
}

func NewBlock(data string, prevBlockHash string) *Block {
	block := &Block{Timestamp: time.Now().Unix(), Data: data, PrevBlockHash: prevBlockHash}
	block.Hash = calculateHash(*block)
	return block
}

func NewBlockchain() *Blockchain {
	genesisBlock := &Block{Timestamp: time.Now().Unix(), Data: "Genesis Block"}
	genesisBlock.Hash = calculateHash(*genesisBlock)
	return &Blockchain{Blocks: []*Block{genesisBlock}}
}

func calculateHash(block Block) string {
	record := string(block.Timestamp) + block.Data + block.PrevBlockHash
	hash := sha256.New()
	hash.Write([]byte(record))
	hashed := hash.Sum(nil)
	return hex.EncodeToString(hashed)
}
