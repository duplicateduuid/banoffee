package main

import (
	"fmt"
	"log"

	"github.com/bxcodec/faker/v3"
	"github.com/jmoiron/sqlx"
)

func seed(db *sqlx.DB) {
	seedUsers(db)
	seedResources(100, db)
}

func seedUsers(db *sqlx.DB) {
	user, err := NewUser("user@user.com", "user", "password", nil, nil, nil)

	if err != nil {
		log.Fatalln(err)
	}

	_, err = db.Exec(
		`INSERT INTO "user" (email, password, username)
		VALUES ($1, $2, $3)`,
		user.Email,
		user.EncryptedPassword,
		user.Username,
	)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("[INFO] [Seeder.resources] User created: %s", user.Email)
}

func seedResources(count int, db *sqlx.DB) {
	for i := 0; i < count; i++ {
		_, err := db.Exec(
			`INSERT INTO "resource" (url, name, image_url, author, description)
			VALUES ($1, $2, $3, $4, $5)`,
			faker.URL(),
			faker.Name(),
			faker.URL(),
			fmt.Sprintf("%s %s", faker.Name(), faker.Name()),
			faker.Sentence(),
		)

		if err != nil {
			log.Fatalln(err)
		}
	}

	fmt.Printf("[INFO] [Seeder.resources] Resources created: %d", count)
}
