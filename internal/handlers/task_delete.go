package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/JuliaKravchenko55/go_final_project/internal/database"
)

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error":"ID задачи обязателен"}`, http.StatusBadRequest)
		return
	}

	res, err := database.DB.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		http.Error(w, `{"error":"Не удалось удалить задачу"}`, http.StatusInternalServerError)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, `{"error":"Не удалось получить количество удаленных строк"}`, http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(map[string]string{}); err != nil {
		http.Error(w, `{"error":"Не удалось закодировать ответ"}`, http.StatusInternalServerError)
	}
}
