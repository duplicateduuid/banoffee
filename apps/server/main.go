package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	db := createDb()

	repositories := NewRepositories(db)
	defer repositories.valkey.Close()

	api := NewAPI(repositories)

	runSeeder := flag.Bool("seed", false, "run all the seeders")
	flag.Parse()

	if *runSeeder {
		seed(&repositories)
	} else {
		api.Run("localhost:6969")
	}
}

func createDb() *sqlx.DB {
	db, err := sqlx.Connect("postgres", "user=postgres dbname=banoffee password=5up3r_s3cur3_p4ssw0rd sslmode=disable")

	if err != nil {
		msg := fmt.Sprintf("[ERROR] failed to create database: %s", err)
		panic(msg)
	}

	return db
}
