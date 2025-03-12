package main

import "github.com/rudransh-shrivastava/rudransh-backend-task/internal/api"

func main() {
	server := api.NewServer()
	server.Run()
}
