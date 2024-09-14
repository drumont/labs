package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"fmt"
)

const (
	host = "host"
	port = 4409
	user = "toto"
	password = "toto"
	dbname = "lab_db"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)

	db, err := sqlx.Open("postgres", psqlInfo)

	if err != nil{
		log.Fatalf("error while connecting to database. err %v", err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Panic")
	}

	log.Println("Successfully connected")

	query := `INSERT INTO users (email, username) VALUES ('aridgiving@ogee.com', 'aridgiving')`

	result, err := db.Exec(query)
	if err != nil {
		log.Fatalf("impossible insert user: %s", err)
	}

	log.Println("inserted data", result)

}
