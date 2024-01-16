package main

import (
	"fmt"
	"testing"
)

func TestGetBalance(t *testing.T) {
	ba := getBalance("0x4A91555066Fe178b29DF226625e87aeBF42b1371")
	fmt.Println(ba)
}
