package main

import (
	"flag"
	"log"

	"github.com/jmoiron/sqlx"
)

func main() {
	db, err := sqlx.Connect("postgres", "user=postgres dbname=banoffee password=5up3r_s3cur3_p4ssw0rd sslmode=disable")

	if err != nil {
		log.Fatalln(err)
	}

	repositories := NewRepositories(db)
	api := NewAPI(repositories)

	runSeeder := flag.Bool("seed", false, "run all the seeders")
	flag.Parse()

	if *runSeeder {
		seed(db)
	} else {
		api.Run("localhost:6969")
	}
}
