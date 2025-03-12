package main

import (
	"github.com/rudransh-shrivastava/rudransh-backend-task/internal/server"
)

func main() {
	server := server.NewServer()
	server.Run()
}
