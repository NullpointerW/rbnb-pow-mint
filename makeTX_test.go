package main

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
)

func BenchmarkTxNew(b *testing.B) {
	Prefix = "0x999999"
	Address = "xxx"
	Address = strings.ToLower(strings.TrimPrefix(Address, "0x"))
	HexAddress = "0x" + Address
	GetAddrHB(Address)
	//for i := 0; i < b.N; i++ {
	//	NewMakeTx()
	//}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// do something
			NewMakeTx()
		}
	})
}
func TestTxNew(t *testing.T) {
	Prefix = "0x999999"
	Address = "xxx"
	Address = strings.ToLower(strings.TrimPrefix(Address, "0x"))
	HexAddress = "0x" + Address
	GetAddrHB(Address)
	s := hex.EncodeToString(AddressHB[:])
	fmt.Println("hbaddr", s)
	//for i := 0; i < 10; i++ {
	NewMakeTx()
	//}
}

func BenchmarkTxOld(b *testing.B) {
	Prefix = "0x999999"
	Address = "xxx"
	Address = strings.ToLower(strings.TrimPrefix(Address, "0x"))
	HexAddress = "0x" + Address
	//for i := 0; i < b.N; i++ {
	//	OldMakeTx()
	//}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// do something
			NewMakeTx()
		}
	})
}
