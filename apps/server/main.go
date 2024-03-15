package main

func main() {
	repositories := NewRepositories()
	api := NewAPI(repositories)

	api.Run("localhost:6969")
}
