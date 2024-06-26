package router

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"

	"github.com/JuliaKravchenko55/go_final_project/internal/handlers"
	"github.com/JuliaKravchenko55/go_final_project/internal/middleware"
	"github.com/JuliaKravchenko55/go_final_project/internal/store"
)

func SetupRouter(store *store.Store) *chi.Mux {
	r := chi.NewRouter()

	absPath, err := filepath.Abs("./web")
	if err != nil {
		log.Fatal("Failed to get absolute path of web directory: ", err)
	}

	r.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir(absPath))))

	r.Post("/api/signin", handlers.Signin)

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth)
		r.Get("/api/nextdate", handlers.CalculateNextDate())
		r.Route("/api/task", func(r chi.Router) {
			r.Post("/", handlers.CreateTask(store))
			r.Get("/", handlers.GetTask(store))
			r.Put("/", handlers.UpdateTask(store))
			r.Delete("/", handlers.DeleteTask(store))
		})
		r.Post("/api/task/done", handlers.CompleteTask(store))
		r.Get("/api/tasks", handlers.ListTasks(store))
	})

	return r
}
