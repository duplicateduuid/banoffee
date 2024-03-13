package main

type Repositories struct {
}

type UserRepository interface {
	CreateUser(*User) error
}
