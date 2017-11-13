package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

type outputKey struct {
	PrivateKey string
	Address    string
}

func NewAllocator() (*outputKey, error) {
	privateKey, err := ecdsa.GenerateKey(secp256k1.S256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	key := fmt.Sprintf("%064x", privateKey.D)
	return &outputKey{PrivateKey: key,
		Address: address.Hex(),
	}, nil
}

func main() {
	output, err := NewAllocator()
	if err != nil {
		log.Fatal(err)
	}
	data, err := json.Marshal(output)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))
}
