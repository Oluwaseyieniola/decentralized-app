// transaction.go
package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
)

type Transaction struct {
	Sender    string
	Receiver  string
	Amount    int
	Signature []byte
}

func NewTransaction(sender, receiver string, amount int, privateKey *ecdsa.PrivateKey) *Transaction {
	tx := &Transaction{Sender: sender, Receiver: receiver, Amount: amount}
	tx.Sign(privateKey)
	return tx
}

func (tx *Transaction) Sign(privateKey *ecdsa.PrivateKey) {
	data := []byte(tx.Sender + tx.Receiver + fmt.Sprintf("%d", tx.Amount))
	hash := sha256.Sum256(data)

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		log.Panic(err)
	}
	signature := append(r.Bytes(), s.Bytes()...)
	tx.Signature = signature
}

func (tx *Transaction) Verify() bool {
	curve := elliptic.P256()

	r := big.Int{}
	s := big.Int{}
	sigLen := len(tx.Signature)
	r.SetBytes(tx.Signature[:(sigLen / 2)])
	s.SetBytes(tx.Signature[(sigLen / 2):])

	data := []byte(tx.Sender + tx.Receiver + fmt.Sprintf("%d", tx.Amount))
	hash := sha256.Sum256(data)

	pubKeyBytes, _ := hex.DecodeString(tx.Sender)
	x := big.Int{}
	y := big.Int{}
	keyLen := len(pubKeyBytes)
	x.SetBytes(pubKeyBytes[:(keyLen / 2)])
	y.SetBytes(pubKeyBytes[(keyLen / 2):])

	rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
	return ecdsa.Verify(&rawPubKey, hash[:], &r, &s)
}
