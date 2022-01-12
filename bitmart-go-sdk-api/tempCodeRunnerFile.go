package main

import (
	"fmt"
	"sync"

	"github.com/bitmartexchange/bitmart-go-sdk-api"
)

func OnMessage(message string) {
	fmt.Println("------------------------>")
}

func main() {
	var wg sync.WaitGroup
	wg.Add(3)

	ws := bitmart.NewWS(bitmart.Config{
		WsUrl:         "wss://ws-manager-compress.bitmart.com/?protocol=1.1",
		ApiKey:        "560f59dfee72957f93b95c00f38f6c03b2b0a0",
		SecretKey:     "7c42ebe0245b6d410343a55cdab54ca4d2a72c2b62704c4b075e5a888da8723f",
		Memo:          "qwertyuiop",
		TimeoutSecond: 10,
		IsPrint:       true,
	})
	_ = ws.Connection(OnMessage)

	channels := []string{
		// public channel
		bitmart.CreateChannel("WS_PUBLIC_SPOT_TICKER", "BTC_USDT"),
		// private channel
		bitmart.CreateChannel("WS_USER_SPOT_ORDER", "BTC_USDT"),
	}
	ws.SubscribeWithLogin(channels)

	// Just test, Please do not use in production.
	wg.Wait()
}
