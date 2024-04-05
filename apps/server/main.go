package main

import (
	"flag"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func main() {
	db := createDb()

	repositories := NewRepositories(db)
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
