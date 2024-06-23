package main

import (
	"log"
	"net/http"

	"github.com/JuliaKravchenko55/go_final_project/internal/store"
	"github.com/joho/godotenv"

	"github.com/JuliaKravchenko55/go_final_project/internal/config"
	"github.com/JuliaKravchenko55/go_final_project/internal/database"
	"github.com/JuliaKravchenko55/go_final_project/internal/router"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	db := database.Initialize()

	storeDB := store.NewStore(db)

	port := config.GetServerPort()

	r := router.SetupRouter(storeDB)

	log.Printf("Server starting on port %s\n", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
