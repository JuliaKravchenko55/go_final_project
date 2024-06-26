package config

import (
	"fmt"
	"os"
)

func GetServerPort() string {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		fmt.Println("TODO_PORT not set, using default port 7540")
		return "7540"
	}
	return port
}
