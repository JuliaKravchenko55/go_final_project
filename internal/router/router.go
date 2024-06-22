package router

import (
	"github.com/JuliaKravchenko55/go_final_project/internal/middleware"
	"log"
	"net/http"
	"path/filepath"

	"github.com/JuliaKravchenko55/go_final_project/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func SetupRouter() *chi.Mux {
	r := chi.NewRouter()

	absPath, err := filepath.Abs("./web")
	if err != nil {
		log.Fatal("Failed to get absolute path of web directory: ", err)
	}

	r.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir(absPath))))

	r.Post("/api/signin", handlers.Signin)

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth)
		r.Get("/api/nextdate", handlers.CalculateNextDate)
		r.Route("/api/task", func(r chi.Router) {
			r.Post("/", handlers.CreateTask)
			r.Get("/", handlers.GetTask)
			r.Put("/", handlers.UpdateTask)
			r.Delete("/", handlers.DeleteTask)
		})
		r.Post("/api/task/done", handlers.CompleteTask)
		r.Get("/api/tasks", handlers.ListTasks)
	})

	return r
}
