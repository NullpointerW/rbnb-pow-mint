package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

var (
	MintLimit     int
	MintApi       string
	MintCount     atomic.Uint64
	Address       string
	Prefix        string
	Challenge     string
	HexAddress    string
	AddressHB     [20]byte
	HttpClient    *http.Client
	Wallets       []Wal
	GoroutineNums int
)

func init() {
	log.SetFormatter(&log.TextFormatter{TimestampFormat: "15:04:05", FullTimestamp: true})
	// load env
	err := godotenv.Load()
	if err != nil {
		fmt.Println("load .env file failed:", err)
		os.Exit(1)
	}
	Prefix = os.Getenv("diff")
	// disable on gpu mint
	err = MakeDifficulty(Prefix)
	if err != nil {
		fmt.Println("invalid diffPrefix", err)
		os.Exit(1)
	}
	BalanceApi = os.Getenv("balanceUrl")
	MintApi = os.Getenv("mintUrl")
	Challenge = os.Getenv("challenge")
	goroutineNums := os.Getenv("workers")
	if goroutineNums == "" {
		GoroutineNums = runtime.NumCPU() + 1
	} else {
		GoroutineNums, err = strconv.Atoi(goroutineNums)
		if err != nil {
			fmt.Println("invalid workerNums:", err)
			os.Exit(1)
		}
	}
	fmt.Println("goroutine num:", GoroutineNums)
	MintLimit, err = strconv.Atoi(os.Getenv("mintLimit"))
	if err != nil {
		fmt.Println("invalid mintLimit", err)
		os.Exit(1)
	}
	shb, err := hex.DecodeString(Challenge)
	if err != nil {
		fmt.Println("get Challenge HexBytes failed:", err)
		os.Exit(1)
	}
	copy(ChallengeHB[:], shb)

	Wallets = LoadWallets("wallet.json")
	Address = strings.ToLower(Wallets[len(Wallets)-1].Addr)
	HexAddress = "0x" + Address
	HttpClient = (&http.Client{Timeout: 10 * time.Second})
}

func main() {
	minted := uint64(getBalance(HexAddress, false))
	GetAddrHB(Address)
	MintCount.Store(minted)
	HashRateStatistic()
Mint:
	ctx, c := context.WithCancel(context.Background())
	for i := 0; i < GoroutineNums; i++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					NewMakeTx()
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
			if mc >= uint64(MintLimit) {
				c()
				break loop
			}
			fmt.Println("address", Address, "mint:", mc)
		}
	}
	// check actually mint balance
	balance := getBalance(HexAddress, true)
	if balance < MintLimit {
		fmt.Println(HexAddress, "实际mint", balance, "未mint：", MintLimit-balance, "继续mint")
		MintCount.Store(uint64(balance))
		goto Mint
	} else {
		fmt.Println(HexAddress, "mint", balance, "已打完，正在创建新钱包...")
	}
	Wallets[len(Wallets)-1].Balance = balance
	addr, pk := genWallet()
	Address = addr
	HexAddress = "0x" + Address
	GetAddrHB(Address)
	wal := Wal{addr, pk, 0}
	Wallets = append(Wallets, wal)
	go StoreWallets(Wallets)
	MintCount.Store(0)
	goto Mint
}

func sendTX(body string) {
	var data = strings.NewReader(body)
	req, err := http.NewRequest("POST", MintApi, data)
	if err != nil {
		log.Error(err)
		return
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("origin", "https://bnb.reth.cc")
	req.Header.Set("referer", "https://bnb.reth.cc/")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	resp, err := HttpClient.Do(req)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "timeout") {
			MintCount.Add(1)
			log.Info("MINT成功(timeout)")
		} else {
			log.Error(err)
		}
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

func GetAddrHB(addr string) {
	shb, _ := hex.DecodeString(addr)
	copy(AddressHB[:], shb)
}
