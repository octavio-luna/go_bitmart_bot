package bitmart

import (
	"encoding/json"
	"fmt"
)

type Currency struct {
	Symbol                      string
	InitialPrice                float32
	PercentSellable             int
	PriceIncreaseToSell         float32
	PercentageIncreaseToSell    int
	PriceToSell                 float32
	PercentageDecreaseToBuyBack int
	PriceDecreaseToBuyBack      float32
	PriceToBuy                  float32
}

func CreateCurrency(client *CloudClient) Currency {
	var check int

	//Set the currency symbol
	var symbol string
	for ok := true; ok; ok = (check < 1) {
		fmt.Println("Type the Symbol of the currency (e.g: ETH, BTC, SHIB): ")
		fmt.Scanln(&symbol)
		_, resp, err := client.GetContractTickersBySymbol(symbol)
		if err != nil {
			panic(err.Error())
		}
		var r Ticker
		err = json.Unmarshal([]byte(resp), &r)
		if err != nil {
			panic(err.Error())
		}
		check = len(r.Data.Tickers)
		if check < 1 {
			fmt.Println("Error. Please try again")
		}
	}
	var c Currency
	c.Symbol = symbol

	//Set the sellable percentage of the currency
	for ok := true; ok; ok = (check <= 0 || check > 100) {
		fmt.Println("Type the percentage of the currency you want to sell: ")
		fmt.Scanln(&check)
		if check <= 0 || check > 100 {
			fmt.Print("Error. ")
		}
	}
	c.PercentSellable = check

	//Set to sell by either a price or a percentage increase
	for ok := true; ok; ok = (check != 1 && check != 2) {
		fmt.Println("Type 1 to sell the asset by percentage increase or 2 to set it by price increase")
		fmt.Scanln(&check)
	}

	if check == 1 {
		for ok := true; ok; ok = (check <= 0 || check > 100) {
			fmt.Println("Type the percentage increase to sell the asset: ")
			fmt.Scanln(&check)
			if check <= 0 || check > 100 {
				fmt.Print("Error. ")
			}
		}
		c.PercentageIncreaseToSell = check
		c.PriceIncreaseToSell = -1
	} else {
		var ch float32
		for ok := true; ok; ok = (ch <= 0) {
			fmt.Println("Type the price increase to sell the asset: ")
			fmt.Scanln(&ch)
			if ch <= 0 {
				fmt.Print("Error. ")
			}
		}
		c.PriceIncreaseToSell = ch
		c.PercentageIncreaseToSell = -1
	}

	//Set to buy by either a price or a percentage increase
	for ok := true; ok; ok = (check != 1 && check != 2) {
		fmt.Println("Type 1 to sell buy back asset by percentage decrease or 2 to set it by price decrease")
		fmt.Scanln(&check)
	}

	if check == 1 {
		for ok := true; ok; ok = (check <= 0 || check > 100) {
			fmt.Println("Type the percentage decrease to buy the asset back: ")
			fmt.Scanln(&check)
			if check <= 0 || check > 100 {
				fmt.Print("Error. ")
			}
		}
		c.PercentageDecreaseToBuyBack = check
		c.PriceDecreaseToBuyBack = -1
	} else {
		var ch float32
		for ok := true; ok; ok = (ch <= 0) {
			fmt.Println("Type the price decrease to sell the asset: ")
			fmt.Scanln(&ch)
			if ch <= 0 {
				fmt.Print("Error. ")
			}
		}
		c.PriceDecreaseToBuyBack = ch
		c.PercentageDecreaseToBuyBack = -1
	}

	return c
}
