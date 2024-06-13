// wallet.go
package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"log"
)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  []byte
}

func NewWallet() *Wallet {
	private, public := newKeyPair()
	return &Wallet{private, public}
}

func newKeyPair() (*ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	pub := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return private, pub
}

func (w *Wallet) GetAddress() string {
	pubKeyHash := hashPublicKey(w.PublicKey)
	address := append(pubKeyHash, checksum(pubKeyHash)...)
	return hex.EncodeToString(address)
}

func hashPublicKey(publicKey []byte) []byte {
	pubSHA256 := sha256.Sum256(publicKey)
	pubRIPEMD160, err := ripemd160Hash(pubSHA256[:])
	if err != nil {
		log.Panic(err)
	}
	return pubRIPEMD160
}

func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])
	return secondSHA[:4]
}
