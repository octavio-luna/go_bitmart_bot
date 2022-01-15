package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/octavio-luna/go_bitmart_bot/bitmart"

	_ "github.com/go-sql-driver/mysql"
)

/*
• Database integration to store current info in case of restart.
• Status of held currencies should be fetched at regular intervals to ensure that the script is trading with actual available funds.
• There should be options to set trade-able assets, default percentage of an asset to sell, the price or percentage increase to sell
 it at and price or percentage decrease to then buy it again (using all of the funds received from that asset). There should also be
 the option to override these settings per asset.
*/

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
	for ok := true; ok; ok = (check != 1) {
		c := bitmart.CreateCurrency(client)
		currencies = append(currencies, c)

		for ok := true; ok; ok = (check != 1 && check != 2) {
			fmt.Println("Press 1 to finish loading the currencies to check or 2 to load another one: ")
			fmt.Scanln(&check)
		}
	}

	for x := 0; x < len(currencies); x++ {
		_, resp, err := client.GetContractTickersBySymbol(currencies[x].Symbol)
		if err != nil {
			log.Panic(err)
		}

		var ticker bitmart.Ticker
		err = json.Unmarshal([]byte(resp), &ticker)
		if err != nil {
			panic(err.Error())
		}
		price, err := strconv.ParseFloat(ticker.Data.Tickers[0].BestAsk, 32)
		if err != nil {
			panic(err.Error())
		}
		currencies[x].InitialPrice = float32(price)

		/*Sets the price to sell the asset either by adding the price increase to the initial price or by multiplying
		the initial price by the percentage coeficient.
		The resaon of having a variable with this data instead of calcullating it on every iteration
		is to improve the performance
		*/
		if currencies[x].PriceIncreaseToSell > 0 {
			currencies[x].PriceToSell = currencies[x].InitialPrice + currencies[x].PriceIncreaseToSell
		} else {
			p := float32(currencies[x].PercentageIncreaseToSell)
			p = (p / 100) + 1
			currencies[x].PriceToSell = p * currencies[x].InitialPrice
		}

		/*Sets the price to buy the asset back either by subtracting the price decrease to the selling price or by multiplying
		the selling price by the percentage coeficient.
		The resaon of having a variable with this data instead of calcullating it on every iteration
		is to improve the performance
		*/
		if currencies[x].PriceDecreaseToBuyBack > 0 {
			currencies[x].PriceToBuy = currencies[x].PriceToSell - currencies[x].PriceDecreaseToBuyBack
		} else {
			p := float32(currencies[x].PriceDecreaseToBuyBack)
			p = (p / 100) + 1
			currencies[x].PriceToBuy = p * currencies[x].InitialPrice
		}
	}

	for x := 0; x < len(currencies); x++ {
		fmt.Println(currencies[x])
	}

	x := 0
	for {
		time.Sleep(10 * time.Second)
		pos := x % len(currencies)
		price, err := client.GetActualPriceSymbol(currencies[pos].Symbol)
		if err != nil {
			panic(err.Error())
		}
		if price >= currencies[pos].PriceToSell && currencies[pos].PriceToSell > 0 {
			amount, err := client.GetAvailableAsset(currencies[pos].Symbol)
			if err != nil {
				panic(err.Error())
			}
			if amount < 0 {
				fmt.Println(price, currencies[pos].PriceToSell, currencies[pos].InitialPrice)
				fmt.Printf("%s has no available founds \n", currencies[pos].Symbol)
				break
			} else {
				//buy
				currencies[pos].PriceToSell = -1
			}
		}
		if price <= currencies[pos].PriceToBuy && currencies[pos].PriceToSell < 0 {
			amount, err := client.GetAvailableAsset(currencies[pos].Symbol)
			if err != nil {
				panic(err.Error())
			}
			if amount < 0 {
				fmt.Printf("%s has no available founds \n", currencies[pos].Symbol)
				break
			} else {
				//buy
				currencies[pos].PercentSellable = 0
			}
		}

	}

	_, now, err := client.GetSystemTime()
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(now)

	fmt.Println("completado")

}
