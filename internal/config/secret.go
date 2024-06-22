package config

import (
	"fmt"
	"os"
)

func GetSecret() string {
	secret := os.Getenv("TODO_SECRET")
	if secret == "" {
		fmt.Println("TODO_SECRET not set, using default secret_key")
		return "secret_key"
	}
	return secret
}

func GetAppPassword() string {
	return os.Getenv("TODO_PASSWORD")
}
