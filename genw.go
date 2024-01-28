package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"os"
)

type Wal struct {
	Addr    string `json:"address"`
	PK      string `json:"privateKey"`
	Balance int    `json:"balance"`
}

func LoadWallets(fp string) []Wal {
	fb, err := os.ReadFile(fp)
	var wallets []Wal
	if os.IsNotExist(err) {
		wal := Wal{}
		wal.Addr, wal.PK = genWallet()
		wallets = append(wallets, wal)
		raw, _ := json.MarshalIndent(wallets, "", "    ")
		f, _ := os.Create("wallet.json")
		_, _ = f.Write(raw)
		_ = f.Close()
		return wallets
	} else {
		_ = json.Unmarshal(fb, &wallets)
		return wallets
	}
}

func StoreWallets(wallets []Wal) {
	f, _ := os.Create("wallet.json")
	raw, _ := json.MarshalIndent(wallets, "", "    ")
	_, err := f.Write(raw)
	if err != nil {
		fmt.Println("write file error", err)
		return
	}

}

func genWallet() (address, pk string) {
	// generate a new private key
	privateKey, _ := crypto.GenerateKey()
	// get the byte form of the private key
	privateKeyBytes := crypto.FromECDSA(privateKey)

	// Get hexString form bytes of the private key
	privateKeyHex := fmt.Sprintf("0x%x", privateKeyBytes)
	fmt.Println("Private Key:", privateKeyHex)

	// get public key from private key
	publicKey := privateKey.Public()

	// generate the byte form of the public key from the public key
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)

	// convert public key bytes to hex string
	publicKeyHex := fmt.Sprintf("0x%x", publicKeyBytes[1:]) // skip the 0x04 before the ECDSA public key
	fmt.Println("Public Key:", publicKeyHex)

	// Generate Ethereum address from public key
	address = crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println("Address:", address)
	return address[2:], privateKeyHex
}
