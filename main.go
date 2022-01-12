package main

import (
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/octavio-luna/small_golang_job"
)

// func storeDBcredentials() (user string, password string, port string, name string) {
// 	scanner := bufio.NewScanner(os.Stdin)
// 	fmt.Println("DB user: ")
// 	scanner.Scan()
// 	user = scanner.Text()

// 	fmt.Printf("%s password: ", user)
// 	scanner.Scan()
// 	password = scanner.Text()

// 	fmt.Println("BD port: ")
// 	scanner.Scan()
// 	port = scanner.Text()

// 	fmt.Println("BD name: ")
// 	scanner.Scan()
// 	name = scanner.Text()

// 	return user, password, port, name
// }

// func searchAPIKey(results *sql.Rows) (memo string, apikey string, secretkey string) {
// 	var tag bitmart.Tag
// 	for results.Next() {
// 		err := results.Scan(&tag.Memo, &tag.Apikey, &tag.Secretkey)
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 	}
// 	return tag.Memo, tag.Apikey, tag.Secretkey
// }

// func storeNewAPIKey() (memo string, apikey string, secretkey string) {
// 	scanner := bufio.NewScanner(os.Stdin)
// 	fmt.Println("Write your memo: ")
// 	scanner.Scan()
// 	memo = scanner.Text()
// 	fmt.Println("Write your APIkey: ")
// 	scanner.Scan()
// 	apikey = scanner.Text()
// 	fmt.Println("Write your SecretKey: ")
// 	scanner.Scan()
// 	secretkey = scanner.Text()
// 	return memo, apikey, secretkey
// }

// func getKey(db *sql.DB) (memo string, apikey string, secretkey string) {
// 	results, err := db.Query("SELECT memo, apikey, secretkey FROM APIconf")
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	memo0, apikey, secretkey := searchAPIKey(results)
// 	fmt.Println("Actual memo and APIkey are  ", memo0, apikey)
// 	fmt.Println("If you wish to change the memo, APIKey or the secretkey press 1: ")
// 	scanner := bufio.NewScanner(os.Stdin)
// 	scanner.Scan()
// 	ch := scanner.Text()

// 	if memo0 == "nil" || ch == "1" {
// 		memo, apikey, secretkey := storeNewAPIKey()
// 		_, err := db.Exec(fmt.Sprintf("UPDATE APIconf set memo = '%s', apikey = '%s', secretkey = '%s' WHERE memo like '%s'", memo, apikey, secretkey, memo0))
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 	}

// 	results, err = db.Query("SELECT memo, apikey, secretkey FROM APIconf")
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	memo, apikey, secretkey = searchAPIKey(results)
// 	return memo, apikey, secretkey
// }

// func connectDB() (db *sql.DB) {
// 	scanner := bufio.NewScanner(os.Stdin)
// 	user, password, port, name := "root", "Octa2003", "3306", "info"
// 	fmt.Println("To change the default values press 1: ")
// 	scanner.Scan()
// 	ch := scanner.Text()
// 	if ch == "1" {
// 		user, password, port, name = storeDBcredentials()
// 	}

// 	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(localhost:%s)/%s", user, password, port, name))
// 	if err != nil {
// 		fmt.Println("error validating sql.open arguments")
// 		panic(err.Error())
// 	}
// 	// defer db.Close()

// 	err = db.Ping()
// 	if err != nil {
// 		fmt.Println("error verifying connection with db.ping")
// 		panic(err.Error())
// 	}
// 	return db
// }

func main() {
	db := bitmart.ConnectDB()
	memo, apikey, secretkey := getKey(db)

	client := bitmart.NewClient(bitmart.Config{
		Url:           "https://api-cloud.bitmart.com", // Ues Https url
		ApiKey:        apikey,
		SecretKey:     secretkey,
		Memo:          memo,
		TimeoutSecond: 10,
		IsPrint:       true,
	})

	_, err, resp := client.GetContractTickersBySymbol("BTCUSDT")
	if err != nil {
		log.Panic(err)
	}

	var ticker bitmart.Ticker
	err = json.Unmarshal([]byte(resp), &ticker)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(ticker.Data.Tickers[0].PriceChangePercent24H)

	fmt.Println("completado")

}
