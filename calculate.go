package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

var (
	GlobalCount atomic.Uint64
	LastCount   uint64 = 0
)

// difficulty is a simple structure for handling mining matching prefixes
//
//		uint32[0]: The uint32 form of the prefixed hexadecimal character
//		uint32[1]: Number of byte arrays prefixed with hexadecimal characters
//	    uint32[2]: Parity of the number of hexadecimal prefix strings, odd numbers are 1, even are 0
//
// This structure is designed to support comparisons of up to 4 bytes (8 hexadecimal characters).
type difficulty [3]uint32

func (d difficulty) GetUint32() uint32 { return d[0] }
func (d difficulty) GetNumCmp() int    { return int(d[1]) }
func (d difficulty) Odd() bool         { return d[2] == 1 }

var Difficulty difficulty

func MakeDifficulty(s string) error {
	parsed, err := strconv.ParseUint(s, 0, 32)
	if err != nil {
		return fmt.Errorf("parse prefix failed: %w", err)
	}
	Difficulty[0] = uint32(parsed)
	strLen := len(strings.TrimPrefix(s, "0x"))
	Difficulty[1] = uint32(strLen / 2)
	if strLen%2 != 0 {
		Difficulty[1] += 1
		Difficulty[2] = 1
	}
	return nil
}

func HashRateStatistic() {
	interval := 5 * time.Second
	timer := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-timer.C:
				count := GlobalCount.Load() - LastCount
				LastCount = GlobalCount.Load()
				fmt.Printf("Hash count %d, Hashrate: %dH/s\n", count, count/5)
			}
		}
	}()
}

func NewMakeTx() {
	// improve
	cHash := Pool.Get()
	defer Pool.Put(cHash)
	err := cHash.P32()
	if err != nil {
		log.Error(err)
		return
	}
	cHash.S20LeftPaddingZero12(AddressHB)
	GlobalCount.Add(1)
	if match(cHash.Keccak256()) {
		potentialSolution := hex.EncodeToString(cHash.GetP32())
		log.WithFields(log.Fields{"Solution": Prefix + "..."}).Info("找到新ID")
		body := fmt.Sprintf(`{"solution": "0x%s", "challenge": "0x%s", "address": "%s", "difficulty": "%s", "tick": "%s"}`, potentialSolution, Challenge, strings.ToLower(HexAddress), Prefix, "rBNB")
		sendTX(body)
	}
}

func OldMakeTx() {
	randomValue := make([]byte, 32)
	_, err := rand.Read(randomValue)
	if err != nil {
		log.Error(err)
		return
	}
	potentialSolution := hex.EncodeToString(randomValue)
	address64 := fmt.Sprintf("%064s", Address)
	dataTemps := fmt.Sprintf(`%s%s%s`, potentialSolution, Challenge, address64)
	dataBytes, err := hex.DecodeString(dataTemps)
	if err != nil {
		fmt.Println("oops!")
		log.Error(err)
		return
	}
	hashedSolutionBytes := crypto.Keccak256(dataBytes)
	hashedSolution := fmt.Sprintf("0x%s", hex.EncodeToString(hashedSolutionBytes))
	GlobalCount.Add(1)
	if strings.HasPrefix(hashedSolution, Prefix) {
		log.WithFields(log.Fields{"Solution": hashedSolution}).Info("找到新ID")
		body := fmt.Sprintf(`{"solution": "0x%s", "challenge": "0x%s", "address": "%s", "difficulty": "%s", "tick": "%s"}`, potentialSolution, Challenge, strings.ToLower(HexAddress), Prefix, "rBNB")
		sendTX(body)
	}
}

func match(k256h []byte) bool {
	var segment uint32
	offsetBits := (Difficulty.GetNumCmp() - 1) * 8
	for i := 0; i < Difficulty.GetNumCmp(); i++ {
		segment |= uint32(k256h[i]) << offsetBits
		offsetBits -= 8
	}
	if Difficulty.Odd() {
		segment >>= 4
	}
	return Difficulty.GetUint32() == segment
}

func M6Match(k256h []byte) bool {
	return uint32(10066329) == uint32(k256h[0])<<16|uint32(k256h[1])<<8|uint32(k256h[2])
}

func M7Match(k256h []byte) bool {
	return uint32(161061273) == (uint32(k256h[0])<<24|uint32(k256h[1])<<16|uint32(k256h[2])<<8|uint32(k256h[3]))>>4
}

func CUDAKeccak256() {
	address64 := fmt.Sprintf("%064s", Address)
	input := fmt.Sprintf(`%s%s`, Challenge, address64)
	cmd := exec.Command("CUDAKeccak256/hash32768.exe", input)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error running hash.exe:", err)
		return
	}
	potentialSolution := strings.TrimSpace(string(output))
	body := fmt.Sprintf(`{"solution": "0x%s", "challenge": "0x%s", "address": "%s", "difficulty": "%s", "tick": "%s"}`, potentialSolution, Challenge, strings.ToLower(HexAddress), Prefix, "rBNB")
	sendTX(body)
	GlobalCount.Add(1)
	//fmt.Println(body)
}
