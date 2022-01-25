package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/octavio-luna/go_bitmart_bot/bitmart"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	/*
		If you're not going to use the DB connection, you can replace the "apikey", "secretkey" and "memo" in the NewClient config, and
		comment the DB connection commands
	*/

	db := bitmart.ConnectDB()
	memo, apikey, secretkey := bitmart.GetKey(db)

	client := bitmart.NewClient(bitmart.Config{
		Url:           "https://api-cloud.bitmart.com", // Ues Https url
		ApiKey:        apikey,
		SecretKey:     secretkey,
		Memo:          memo,
		TimeoutSecond: 10,
		IsPrint:       true,
	})

	var currencies []bitmart.Currency

	var check int

	for ok := true; ok; ok = (check != 1 && check != 2) {
		fmt.Println("Type 1 to use the current configuration or 2 to set your own: ")
		fmt.Scanln(&check)
		if check != 1 && check != 2 {
			fmt.Print("Error. ")
		}
	}

	if check == 2 {
		for ok := true; ok; ok = (check != 1) {
			c := bitmart.CreateCurrency(client)
			if c.Symbol == "nil" {
				//Do not append, DEBUG function
			} else {
				currencies = append(currencies, c)
			}

			for ok := true; ok; ok = (check != 1 && check != 2) {
				fmt.Println("Press 1 to finish loading the currencies to check or 2 to load another one: ")
				fmt.Scanln(&check)
			}
		}
		bitmart.SetCurrencies(bitmart.OpenFile(), currencies)
	} else {
		currencies = bitmart.GetCurrencies(bitmart.OpenFile())
	}

	x := 0
	if len(currencies) > 0 {
		for {
			time.Sleep(10 * time.Second)
			pos := x % len(currencies)
			price, err := client.GetActualPriceSymbol(currencies[pos].Symbol)
			if err != nil {
				panic(err.Error())
			}
			if len(currencies[pos].PriceToSell) > 0 {
				if price >= currencies[pos].PriceToSell[len(currencies[pos].PriceToSell)-1] {
					amount, err := client.GetAvailableAsset(currencies[pos].Symbol)
					if err != nil {
						panic(err.Error())
					}
					if amount < 0 || amount < currencies[pos].AmountSellable {
						fmt.Println(price, currencies[pos].PriceToSell, currencies[pos].InitialPrice)
						fmt.Printf("%s has no available founds \n", currencies[pos].Symbol)
						break
					} else {
						size := currencies[pos].AmountSellable / float32(len(currencies[pos].PriceToSell))

						_, resp, err := client.PostSpotSubmitMarketSellOrder(bitmart.MarketSellOrder{Symbol: fmt.Sprintf("%s_USDT", currencies[pos].Symbol), Size: fmt.Sprintf("%f", size)})
						if err != nil {
							panic(err.Error())
						}
						var order bitmart.Order
						err = json.Unmarshal([]byte(resp), &order)
						if err != nil {
							panic(err.Error())
						}
						fmt.Println("Order id: ", order.Data.OrderID)
						currencies[pos].AmountSellable = currencies[pos].AmountSellable - size
						currencies[pos].PriceToSell = currencies[pos].PriceToSell[:len(currencies[pos].PriceToSell)-2]
						currencies[pos].DollarsToBuy += price * size

						_, now, err := client.GetSystemTime()
						if err != nil {
							panic(err.Error())
						}
						bitmart.InsertConsult(db, currencies[pos].Symbol, now, price, "sell")
					}
				}
			} else if len(currencies[pos].PriceToBuy) > 0 {
				if price <= currencies[pos].PriceToBuy[len(currencies[pos].PriceToBuy)-1] {
					amount, err := client.GetAvailableAsset("USDT")
					if err != nil {
						panic(err.Error())
					}
					if amount <= 0 || amount < currencies[pos].DollarsToBuy/float32(len(currencies[pos].PriceToSell)) {
						fmt.Printf("%s has no available founds \n", currencies[pos].Symbol)
						break
					} else {
						value := currencies[pos].DollarsToBuy / float32(len(currencies[pos].PriceToSell))
						_, resp, err := client.PostSpotSubmitMarketBuyOrder(bitmart.MarketBuyOrder{Symbol: fmt.Sprintf("%s_USDT", currencies[pos].Symbol), Notional: strconv.FormatFloat(float64(value), 'E', -1, 32)})
						if err != nil {
							panic(err.Error())
						}
						var order bitmart.Order
						err = json.Unmarshal([]byte(resp), &order)
						if err != nil {
							panic(err.Error())
						}
						fmt.Println("Order id: ", order.Data.OrderID)
						currencies[pos].PriceToBuy = currencies[pos].PriceToBuy[:len(currencies[pos].PriceToBuy)-2]
						currencies[pos].DollarsToBuy -= value

						_, now, err := client.GetSystemTime()
						if err != nil {
							panic(err.Error())
						}
						bitmart.InsertConsult(db, currencies[pos].Symbol, now, price, "sell")
					}
				}
			}
			_, now, err := client.GetSystemTime()
			if err != nil {
				panic(err.Error())
			}
			bitmart.InsertConsult(db, currencies[pos].Symbol, now, price, "getprice")
			x++
		}
	}

	fmt.Println("completado")

}
