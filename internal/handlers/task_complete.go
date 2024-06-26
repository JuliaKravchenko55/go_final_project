package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/JuliaKravchenko55/go_final_project/internal/store"
)

func CompleteTask(store *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, `{"error":"ID задачи обязателен"}`, http.StatusBadRequest)
			return
		}

		if err := store.CompleteTask(id); err != nil {
			if err.Error() == "task not found" {
				http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
			} else {
				http.Error(w, `{"error":"Не удалось выполнить задачу"}`, http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(map[string]string{}); err != nil {
			http.Error(w, `{"error":"Не удалось закодировать ответ"}`, http.StatusInternalServerError)
		}
	}
}
