package main

import (
	"fmt"

	"github.com/bxcodec/faker/v3"
)

func (s *Seeder) resourceSeed(countInput ...int) {
	count := 100
	if len(countInput) > 0 {
		count = countInput[0]
	}

	// TODO: update golang to +1.22 and use `for i := range count`
	for i := 0; i < count; i++ {
		_, err := s.db.Exec(
			`INSERT INTO "resource" (url, name, image_url, author, description)
			VALUES ($1, $2, $3, $4, $5)`,
			faker.URL(),
			faker.Name(),
			faker.URL(),
			(faker.FirstName() + " " + faker.LastName()),
			faker.Sentence(),
		)

		if err != nil {
			errorMessage := fmt.Sprintf("[ERROR] [Seeder.resources] failed to create resource: %v\n", err)
			panic(errorMessage)
		}
	}
}
