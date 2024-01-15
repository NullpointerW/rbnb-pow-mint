package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync/atomic"
	"time"
)

var BalanceAPI = "https://ec2-18-217-135-255.us-east-2.compute.amazonaws.com/balance?address=%s"

var (
	MintCount  atomic.Uint64
	Address    string
	Prefix     string
	Challenge  string
	HexAddress string
)

func init() {
	Challenge = "72424e4200000000000000000000000000000000000000000000000000000000"
	log.SetFormatter(&log.TextFormatter{TimestampFormat: "15:04:05", FullTimestamp: true})
	fmt.Print("请输入地址：")
	_, err := fmt.Scanln(&Address)
	if err != nil {
		return
	}
	Address = strings.ToLower(strings.TrimPrefix(Address, "0x"))
	HexAddress = "0x" + Address
	fmt.Print("请输入难度：")
	_, err = fmt.Scanln(&Prefix)
	if err != nil {
		return
	}
}

func main() {
	url := fmt.Sprintf(BalanceAPI, HexAddress)
	minted := uint64(0)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("query balance error", err)
	} else {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("read balance error", err)
		} else {
			var rmap map[string]any
			err := json.Unmarshal(b, &rmap)
			if err != nil {
				fmt.Println("unmarshal balance error", err)
			} else {
				minted = uint64((rmap["balance"]).(float64))
			}
		}
	}
	MintCount.Store(minted)
Mint:
	ctx, c := context.WithCancel(context.Background())
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					makeTx()
				}

			}
		}()
	}
	tick := time.NewTicker(3 * time.Second)
loop:
	for {
		select {
		case <-tick.C:
			mc := MintCount.Load()
			if mc >= 4900 {
				c()
				break loop
			}
			fmt.Println("address", Address, "mint:", mc)
		}
	}
	addr, pk := genWallet()
	Address = addr
	HexAddress = "0x" + Address
	go WritePk(addr, pk)
	MintCount.Store(0)
	goto Mint
}

func WritePk(addr, pk string) {
	file, err := os.OpenFile("wal.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("open file error", err)
		return
	}
	defer file.Close()
	w := bufio.NewWriter(file)
	d := fmt.Sprintf("adrr:%s|pk:%s\n", addr, pk)
	_, err = w.Write([]byte(d))
	if err != nil {
		fmt.Println("write file error", err)
		return
	}
	_ = w.Flush()
}

func sendTX(body string) {
	client := &http.Client{}
	var data = strings.NewReader(body)
	req, err := http.NewRequest("POST", "https://ec2-18-217-135-255.us-east-2.compute.amazonaws.com/validate", data)
	if err != nil {
		log.Error(err)
		return
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("origin", "https://bnb.reth.cc")
	req.Header.Set("referer", "https://bnb.reth.cc/")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error(err)
			return
		}
	}(resp.Body)

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return
	}
	bodyString := string(bodyText)
	containsValidateSuccess := strings.Contains(bodyString, "validate success!")
	if containsValidateSuccess {
		log.Info("MINT成功")
		MintCount.Add(1)
	} else {
		log.WithFields(log.Fields{"错误": err}).Error("MINT错误")
	}

}

func makeTx() {
	randomValue := make([]byte, 32)
	_, err := rand.Read(randomValue)
	if err != nil {
		log.Error(err)
		return
	}

	potentialSolution := hex.EncodeToString(randomValue)
	//fmt.Println("hex address", Address)
	address64 := fmt.Sprintf("%064s", strings.ToLower(Address))
	dataTemps := fmt.Sprintf(`%s%s%s`, potentialSolution, Challenge, address64)

	dataBytes, err := hex.DecodeString(dataTemps)
	if err != nil {
		fmt.Println("oops!")
		log.Error(err)
		return
	}

	hashedSolutionBytes := crypto.Keccak256(dataBytes)
	hashedSolution := fmt.Sprintf("0x%s", hex.EncodeToString(hashedSolutionBytes))

	if strings.HasPrefix(hashedSolution, Prefix) {
		log.WithFields(log.Fields{"Solution": hashedSolution}).Info("找到新ID")
		body := fmt.Sprintf(`{"solution": "0x%s", "challenge": "0x%s", "address": "%s", "difficulty": "%s", "tick": "%s"}`, potentialSolution, Challenge, strings.ToLower(HexAddress), Prefix, "rBNB")
		sendTX(body)
	}
}
