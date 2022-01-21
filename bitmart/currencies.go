package bitmart

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

type Currency struct {
	Symbol         string
	InitialPrice   float32
	AmountSellable float32
	PriceToSell    []float32
	PriceToBuy     []float32
	DollarsToBuy   float32
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

	//Set the initial price
	_, resp, err := client.GetContractTickersBySymbol(c.Symbol)
	if err != nil {
		log.Panic(err)
	}

	var ticker Ticker
	err = json.Unmarshal([]byte(resp), &ticker)
	if err != nil {
		panic(err.Error())
	}
	price, err := strconv.ParseFloat(ticker.Data.Tickers[0].BestAsk, 32)
	if err != nil {
		panic(err.Error())
	}
	c.InitialPrice = float32(price)

	//Set the sellable percentage of the currency
	for ok := true; ok; ok = (check <= 0 || check > 100) {
		fmt.Println("Type the percentage of the currency you want to sell: ")
		fmt.Scanln(&check)
		if check <= 0 || check > 100 {
			fmt.Print("Error. ")
		}
	}

	amount, err := client.GetAvailableAsset(c.Symbol)
	if err != nil {
		panic(err.Error())
	}
	if amount <= 0 {
		fmt.Println(price, c.PriceToSell, c.InitialPrice)
		fmt.Printf("%s has no available founds \n", c.Symbol)
		var cur Currency
		cur.Symbol = "nil"
		return cur
	} else {
		p := float32(check)
		p = (p / 100)
		amount = amount * p
	}
	c.AmountSellable = amount

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
		p := float32(check)
		p = (p / 100) + 1
		c.PriceToSell = append(c.PriceToSell, (p * c.InitialPrice))
	} else {
		var ch float32
		for ok := true; ok; ok = (ch <= 0) {
			fmt.Println("Type the price increase to sell the asset: ")
			fmt.Scanln(&ch)
			if ch <= 0 {
				fmt.Print("Error. ")
			}
		}
		c.PriceToSell = append(c.PriceToSell, (c.InitialPrice + float32(ch)))
	}

	for ok := true; ok; ok = (check < 1 || check > 5) {
		fmt.Println("In how many steps is the operation going to be be made? (max 5) ")
		fmt.Scanln(&check)
		if check < 1 || check > 5 {
			fmt.Print("Error. ")
		}
	}

	diff := c.PriceToSell[0] - c.InitialPrice
	for x := 1; x < check; x++ {
		c.PriceToSell = append(c.PriceToSell, (c.PriceToSell[x-1] - (diff / 2)))
		diff = diff / 2
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
		p := float32(check)
		p = 1 - (p / 100)
		c.PriceToBuy = append(c.PriceToBuy, (p * c.InitialPrice))
	} else {
		var ch float32
		for ok := true; ok; ok = (ch <= 0) {
			fmt.Println("Type the price decrease to sell the asset: ")
			fmt.Scanln(&ch)
			if ch <= 0 {
				fmt.Print("Error. ")
			}
		}
		c.PriceToBuy = append(c.PriceToBuy, (c.PriceToSell[0] - ch))
	}

	for ok := true; ok; ok = (check < 1 || check > 5) {
		fmt.Println("In how many steps is the operation going to be made? (max 5) ")
		fmt.Scanln(&check)
		if check < 1 || check > 5 {
			fmt.Print("Error. ")
		}
	}

	diff = c.PriceToBuy[0] - c.InitialPrice
	for x := 1; x < check; x++ {
		c.PriceToBuy = append(c.PriceToBuy, (c.PriceToBuy[x-1] - (diff / 2)))
		diff = diff / 2
	}

	c.DollarsToBuy = 0

	return c
}
