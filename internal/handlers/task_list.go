package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/JuliaKravchenko55/go_final_project/internal/models"
	"github.com/JuliaKravchenko55/go_final_project/internal/store"
)

func ListTasks(store *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		search := r.URL.Query().Get("search")

		tasks, err := store.ListTasks(search)
		if err != nil {
			fmt.Printf("Error getting task: %v", err)
			http.Error(w, `{"error":"Не удалось получить задачи"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(map[string][]models.Task{"tasks": tasks}); err != nil {
			http.Error(w, `{"error":"Не удалось закодировать задачи"}`, http.StatusInternalServerError)
		}
	}
}
