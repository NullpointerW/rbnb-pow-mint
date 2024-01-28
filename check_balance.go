package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var BalanceClient = (&http.Client{Timeout: 20 * time.Second})
var BalanceApi string

type ApiResponse struct {
	Address string `json:"address"`
	Balance int    `json:"balance"`
}

func getBalance(address string, requireBalance bool) int {
	url := BalanceApi + address
	for {
		time.Sleep(time.Millisecond * 500)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Create balance request failed:", err)
			continue
		}

		resp, err := BalanceClient.Do(req)
		if err != nil {
			fmt.Println("Request for balance failed, retrying:", err)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			fmt.Println("Failed to read balance response body:", err)
			continue
		}

		var response ApiResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			fmt.Println("Failed to parse balance JSON:", err)
			continue
		}
		fmt.Println(string(body))
		if response.Address == "" || requireBalance && response.Balance == 0 {
			fmt.Println("Unexpected balance response,retrying")
			continue
		}
		return response.Balance
	}
}
