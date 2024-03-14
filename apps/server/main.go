package main

func main() {
	repositories := NewPostgresRepositories()
	api := NewAPI(repositories)

	api.Run("localhost:6969")
}
