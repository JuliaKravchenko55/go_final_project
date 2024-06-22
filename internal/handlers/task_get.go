package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/JuliaKravchenko55/go_final_project/internal/database"
	"github.com/JuliaKravchenko55/go_final_project/internal/models"
)

func GetTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error":"ID задачи обязателен"}`, http.StatusBadRequest)
		return
	}

	row := database.DB.QueryRow(`SELECT * FROM scheduler WHERE id = ?`, id)
	var task models.Task
	if err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error":"Не удалось получить задачу"}`, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, `{"error":"Не удалось закодировать ответ"}`, http.StatusInternalServerError)
	}
}
