// package main

// import (
// 	"log"

// 	"github.com/bitmartexchange/bitmart-go-sdk-api"
// )

// func main() {

// 	client := bitmart.NewClient(bitmart.Config{
// 		Url:           "https://api-cloud.bitmart.com", // Ues Https url
// 		ApiKey:        "ec6b5dacf5d7ed3ea0d2934c05add27d830103ec",
// 		SecretKey:     "39dc92125f1d223e2eb34fce251eb806a2e538c6e553e2c50a97abf7fa33addd",
// 		Memo:          "11239242438",
// 		TimeoutSecond: 10,
// 		IsPrint:       true,
// 	})

// 	var ac, err = client.PostSpotSubmitLimitBuyOrder(bitmart.LimitBuyOrder{Symbol: "BTC_USDT", Size: "8800", Price: "0.01"})
// 	if err != nil {
// 		log.Panic(err)
// 	} else {
// 		bitmart.PrintResponse(ac)
// 	}

// }

// #### WebSocket Example

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
		ApiKey:        "Your API KEY",
		SecretKey:     "Your Secret KEY",
		Memo:          "Your Memo",
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
