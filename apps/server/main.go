package main

func main() {
	repositories := NewPostgresRepositories()
	api := NewAPI("localhost:6969", *repositories)

	api.Run()
}
