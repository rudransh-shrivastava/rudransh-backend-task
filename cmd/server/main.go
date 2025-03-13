package main

import "github.com/rudransh-shrivastava/rudransh-backend-task/internal/api"

// Entry point of the application
func main() {
	server := api.NewServer()
	server.Run()
}
