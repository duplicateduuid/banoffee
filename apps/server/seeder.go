package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/bxcodec/faker/v3"
)

func seed(r *Repositories) {
	users := seedUsers(r, 5)
	resources := seedResources(r, 25)
	seedUserResources(r, users, resources)
}

func seedUsers(r *Repositories, count int) []User {
	user, err := NewUser("user@user.com", "user", "password", nil, nil, nil)

	if err != nil {
		log.Fatalln(err)
	}

	user, _ = r.userRepository.CreateUser(user)

	fmt.Printf("[INFO] [Seeder.resources] Default user created %v\n", user)

	users := []User{}

	for i := 0; i < count; i++ {
		user, err = NewUser(faker.Email(), faker.Username(), faker.Password(), nil, nil, nil)
		user, err = r.userRepository.CreateUser(user)

		users = append(users, *user)
	}

	fmt.Printf("[INFO] [Seeder.resources] Users created: %d\n", len(users))

	return users
}

func seedResources(r *Repositories, count int) []Resource {
	resources := []Resource{}

	for i := 0; i < count; i++ {
		imageUrl := faker.URL()
		authorName := fmt.Sprintf("%s %s", faker.Name(), faker.Name())
		description := faker.Sentence()

		resource := NewResource(
			faker.URL(),
			faker.Name(),
			&imageUrl,
			&authorName,
			&description,
		)

		resource, err := r.resourceRepository.CreateResource(resource)
		resources = append(resources, *resource)

		if err != nil {
			log.Fatalln(err)
		}
	}

	fmt.Printf("[INFO] [Seeder.resources] Resources created: %d\n", len(resources))

	return resources
}

func generateRating() string {
	ratings := []string{"one", "two", "three", "four", "five"}
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := rand.Intn(len(ratings))

	return ratings[index]
}

func seedUserResources(r *Repositories, users []User, resources []Resource) {
	for userIndex := range users {
		for resourceIndex := range resources {
			user := users[userIndex]
			resource := resources[resourceIndex]

			var status *string
			var rating *string

			if resourceIndex < 1*(len(resources)/4) {
				status = new(string)
				rating = new(string)
				*status = "ongoing"
				*rating = generateRating()
			} else if resourceIndex < 2*(len(resources)/4) {
				status = new(string)
				rating = new(string)
				*status = "completed"
				*rating = generateRating()
			} else if resourceIndex < 3*(len(resources)/4) {
				status = new(string)
				rating = new(string)
				*status = "bookmarked"
				*rating = generateRating()
			} else {
				continue
			}

			err := r.userRepository.CreateUserResource(&user, resource.Id.String(), status, rating, nil)

			if err != nil {
				fmt.Printf("[ERROR] [Seeder.resources] Failed to link user(id=%s) with resource(id=%s): %s\n", user.Id, resource.Id, err)
				log.Fatalln(err)
			}
		}
	}

	fmt.Printf("[INFO] [Seeder.resources] Linked every user to every resource\n")
}
