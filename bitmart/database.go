package bitmart

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Tag struct {
	Memo      string `json:"memo"`
	Apikey    string `json:"apikey"`
	Secretkey string `json:"secretkey"`
}

//Replace 'user', 'pasword' and 'port' for your own database credentials and ports
func ConnectDB() (db *sql.DB) {
	scanner := bufio.NewScanner(os.Stdin)
	user, password, port, name := "root", "Octa2003", "3306", "info"
	fmt.Println("To change the default values press 1: ")
	scanner.Scan()
	ch := scanner.Text()
	if ch == "1" {
		user, password, port, name = StoreDBcredentials()
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(localhost:%s)/%s", user, password, port, name))
	if err != nil {
		fmt.Println("error validating sql.open arguments")
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("error verifying connection with db.ping")
		panic(err.Error())
	}
	return db
}

//looks for your bitmart credentials on the database and returns the variables to the main file
func SearchAPIKey(results *sql.Rows) (memo string, apikey string, secretkey string) {
	var tag Tag
	for results.Next() {
		err := results.Scan(&tag.Memo, &tag.Apikey, &tag.Secretkey)
		if err != nil {
			panic(err.Error())
		}
	}
	return tag.Memo, tag.Apikey, tag.Secretkey
}

func StoreDBcredentials() (user string, password string, port string, name string) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("DB user: ")
	scanner.Scan()
	user = scanner.Text()

	fmt.Printf("%s password: ", user)
	scanner.Scan()
	password = scanner.Text()

	fmt.Println("BD port: ")
	scanner.Scan()
	port = scanner.Text()

	fmt.Println("BD name: ")
	scanner.Scan()
	name = scanner.Text()

	return user, password, port, name
}

func StoreNewAPIKey() (memo string, apikey string, secretkey string) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Write your memo: ")
	scanner.Scan()
	memo = scanner.Text()
	fmt.Println("Write your APIkey: ")
	scanner.Scan()
	apikey = scanner.Text()
	fmt.Println("Write your SecretKey: ")
	scanner.Scan()
	secretkey = scanner.Text()
	return memo, apikey, secretkey
}

//Looks for the actually stored values
func GetKey(db *sql.DB) (memo string, apikey string, secretkey string) {
	results, err := db.Query("SELECT memo, apikey, secretkey FROM APIconf")
	if err != nil {
		panic(err.Error())
	}

	memo0, apikey, secretkey := SearchAPIKey(results)
	fmt.Printf("Actual memo and APIkey are %s %s \n", memo0, apikey)
	fmt.Println("If you wish to change the memo, APIKey or the secretkey press 1: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	ch := scanner.Text()

	if memo0 == "nil" || ch == "1" {
		memo, apikey, secretkey := StoreNewAPIKey()
		_, err := db.Exec(fmt.Sprintf("UPDATE APIconf set memo = '%s', apikey = '%s', secretkey = '%s' WHERE memo like '%s'", memo, apikey, secretkey, memo0))
		if err != nil {
			panic(err.Error())
		}
	}

	results, err = db.Query("SELECT memo, apikey, secretkey FROM APIconf")
	if err != nil {
		panic(err.Error())
	}

	memo, apikey, secretkey = SearchAPIKey(results)
	return memo, apikey, secretkey
}

func InsertConsult(db *sql.DB, symbol string, time string, price float32, op string) {
	_, err := db.Exec(fmt.Sprintf("Insert into consults (symbol, moment, price, op) values (%s, %s, %f, %s);", symbol, time, price, op))
	if err != nil {
		panic(err.Error())
	}
}
