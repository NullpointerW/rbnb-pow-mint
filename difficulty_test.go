package main

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestDifficulty(t *testing.T) {
	b, _ := hex.DecodeString("99999999b9c3a331c60cc883cf5c66ebe1dd244a6f2e87363b4b29756b1f17a762")
	_ = MakeDifficulty("0x99999999")
	fmt.Println(match(b))
	fmt.Println(M7Match(b))

}
