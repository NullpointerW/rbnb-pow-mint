package main

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/crypto"
	"sync"
)

var ChallengeHB [32]byte
var pLen = 32

var Pool = pool{Pool: sync.Pool{New: func() any {
	return NewCalculateHash()
}}}

type pool struct{ sync.Pool }

func (p *pool) Get() *CalculatedHash  { return p.Pool.Get().(*CalculatedHash) }
func (p *pool) Put(h *CalculatedHash) { p.Pool.Put(h) }

type CalculatedHash [96]byte

func NewCalculateHash() *CalculatedHash {
	h := new(CalculatedHash)
	h.m32(ChallengeHB)
	return h
}

func (ch *CalculatedHash) P32() error {
	_, err := rand.Read(ch[:pLen])
	return err
}

func (ch *CalculatedHash) m32(src [32]byte) { copy(ch[pLen:pLen*2], src[:]) }

func (ch *CalculatedHash) S20LeftPaddingZero12(src [20]byte) { copy(ch[96-20:], src[:]) }

func (ch *CalculatedHash) GetP32() []byte {
	return ch[:pLen]
}

func (ch *CalculatedHash) Keccak256() []byte { return crypto.Keccak256(ch[:]) }

func (ch *CalculatedHash) ToHexString() string {
	return hex.EncodeToString(ch[:])
}
