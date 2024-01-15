package main

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
)

func genWallet() (address, pk string) {
	// 生成一个新的私钥
	privateKey, _ := crypto.GenerateKey()
	// 获取私钥的字节形式
	privateKeyBytes := crypto.FromECDSA(privateKey)

	// 将私钥字节转换为十六进制字符串
	privateKeyHex := fmt.Sprintf("0x%x", privateKeyBytes)
	fmt.Println("Private Key:", privateKeyHex)

	// 从私钥生成公钥
	publicKey := privateKey.Public()

	// 从公钥生成公钥的字节形式
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)

	// 将公钥字节转换为十六进制字符串
	publicKeyHex := fmt.Sprintf("0x%x", publicKeyBytes[1:]) // 跳过ECDSA公钥前的0x04
	fmt.Println("Public Key:", publicKeyHex)

	// 从公钥生成以太坊地址
	address = crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println("Address:", address)
	return address[2:], privateKeyHex
}
