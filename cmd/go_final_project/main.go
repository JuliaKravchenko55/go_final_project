package main

import (
	"github.com/JuliaKravchenko55/go_final_project/internal/config"
	"github.com/JuliaKravchenko55/go_final_project/internal/database"
	"github.com/JuliaKravchenko55/go_final_project/internal/router"
	"log"
	"net/http"
)

func main() {
	database.Initialize()

	port := config.GetServerPort()

	r := router.SetupRouter()

	log.Printf("Server starting on port %s\n", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
