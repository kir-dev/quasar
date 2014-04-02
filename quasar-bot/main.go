package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
)

func main() {
	fmt.Println("Hello from quasar-bot!")

	c := config{}
	// TODO: make config path configurable via command line args
	configData, err := ioutil.ReadFile("config/config.json")
	if err != nil {
		log.Fatalln("Could not read config file:", err)
	}

	if err = json.Unmarshal(configData, &c); err != nil {
		log.Fatalln("Malformed config:", err)
	}

	db, err := sql.Open("postgres", fmt.Sprintf("dbname=%s user=%s password=%s", c.Database, c.DbUser, c.DbPass))
	if err != nil {
		log.Fatalln("Could not connect to database:", err)
	}

	var count int
	db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	fmt.Printf("There are %d user(s) in the database.\n", count)
}
