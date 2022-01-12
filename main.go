package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/bitmartexchange/bitmart-go-sdk-api"
	_ "github.com/go-sql-driver/mysql"
)

//Esto deberia reevaluarse
func cargar() (user string, password string, port string, name string) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Ingrese el usuario de la BD: ")
	scanner.Scan()
	user = scanner.Text()

	fmt.Printf("Ingrese el password de %s: ", user)
	scanner.Scan()
	password = scanner.Text()

	fmt.Println("Ingrese el puerto de la BD: ")
	scanner.Scan()
	port = scanner.Text()

	fmt.Println("Ingrese el nombre de la BD: ")
	scanner.Scan()
	name = scanner.Text()

	return user, password, port, name
}

func searchAPIKey(results *sql.Rows) (memo string, apikey string, secretkey string) {
	var tag Tag
	for results.Next() {
		// var tag Tag
		// scan the result into our tag composite object
		err := results.Scan(&tag.memo, &tag.apikey, &tag.secretkey)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		// and then print out the tag's Name attribute
	}
	return tag.memo, tag.apikey, tag.secretkey
}

func storeNewAPIKey() (memo string, apikey string, secretkey string) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Ingrese el memo: ")
	scanner.Scan()
	memo = scanner.Text()
	fmt.Println("Ingrese el APIkey: ")
	scanner.Scan()
	apikey = scanner.Text()
	fmt.Println("Ingrese el SecretKey: ")
	scanner.Scan()
	secretkey = scanner.Text()
	return memo, apikey, secretkey
}

func getKey(db *sql.DB) (memo string, apikey string, secretkey string) {
	results, err := db.Query("SELECT memo, apikey, secretkey FROM APIconf")
	if err != nil {
		panic(err.Error())
	}

	memo0, apikey, secretkey := searchAPIKey(results)
	fmt.Println("El memo y la apikey actuales son ", memo0, apikey)
	fmt.Println("si desea cambiar el memo, la APIKey o la secret key ingrese 1: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	ch := scanner.Text()

	if memo0 == "nil" || ch == "1" {
		memo, apikey, secretkey := storeNewAPIKey()
		_, err := db.Exec(fmt.Sprintf("UPDATE APIconf set memo = '%s', apikey = '%s', secretkey = '%s' WHERE memo like '%s'", memo, apikey, secretkey, memo0))
		if err != nil {
			panic(err.Error())
		}
	}

	results, err = db.Query("SELECT memo, apikey, secretkey FROM APIconf")
	if err != nil {
		panic(err.Error())
	}

	memo, apikey, secretkey = searchAPIKey(results)
	return memo, apikey, secretkey
}

func connectDB() (db *sql.DB) {
	scanner := bufio.NewScanner(os.Stdin)
	user, password, port, name := "root", "Octa2003", "3306", "info"
	fmt.Println("Si desea cambiar los valores por defecto ingrese 1: ")
	scanner.Scan()
	ch := scanner.Text()
	if ch == "1" {
		user, password, port, name = cargar()
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(localhost:%s)/%s", user, password, port, name))
	if err != nil {
		fmt.Println("error validating sql.open arguments")
		panic(err.Error())
	}
	// defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println("error verifying connection with db.ping")
		panic(err.Error())
	}
	return db
}

type Tag struct {
	memo      string `json:"memo"`
	apikey    string `json:"apikey"`
	secretkey string `json:"secretkey"`
}

type Ticker struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Trace   string `json:"trace"`
	Data    struct {
		Tickers []struct {
			ContractSymbol        string `json:"contract_symbol"`
			LastPrice             string `json:"last_price"`
			IndexPrice            string `json:"index_price"`
			LastFundingRate       string `json:"last_funding_rate"`
			PriceChangePercent24H string `json:"price_change_percent_24h"`
			Volume24H             string `json:"volume_24h"`
			URL                   string `json:"url"`
			HighPrice             string `json:"high_price"`
			LowPrice              string `json:"low_price"`
			LegalCoinPrice        string `json:"legal_coin_price"`
		} `json:"tickers"`
	} `json:"data"`
}

func main() {
	db := connectDB()
	memo, apikey, secretkey := getKey(db)

	client := bitmart.NewClient(bitmart.Config{
		Url:           "https://api-cloud.bitmart.com", // Ues Https url
		ApiKey:        apikey,
		SecretKey:     secretkey,
		Memo:          memo,
		TimeoutSecond: 10,
		IsPrint:       true,
	})

	var ac, err, resp = client.GetContractTickersBySymbol("BTCUSDT")
	if err != nil {
		log.Panic(err)
	}
	//  else {
	// 	bitmart.PrintResponse(ac)
	// }

	var ticker Ticker
	err = json.Unmarshal([]byte(resp), &ticker)

	if err != nil {
		panic(err.Error())
	}

	fmt.Println(resp)
	fmt.Println(ticker.Data.Tickers[0].PriceChangePercent24H)

	fmt.Println("completado")

}
