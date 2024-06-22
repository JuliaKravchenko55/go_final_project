package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/JuliaKravchenko55/go_final_project/internal/models"
	"net/http"
	"time"

	"github.com/JuliaKravchenko55/go_final_project/internal/database"
	"github.com/JuliaKravchenko55/go_final_project/internal/utils"
)

func CompleteTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Метод не разрешен"}`, http.StatusMethodNotAllowed)
		return
	}

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

	if task.Repeat == "" {
		if _, err := database.DB.Exec(`DELETE FROM scheduler WHERE id = ?`, id); err != nil {
			http.Error(w, `{"error":"Не удалось удалить задачу"}`, http.StatusInternalServerError)
			return
		}
	} else {
		nextDate, err := utils.CalculateNextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			http.Error(w, `{"error":"Не удалось рассчитать следующую дату"}`, http.StatusInternalServerError)
			return
		}

		if _, err := database.DB.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`, nextDate, id); err != nil {
			http.Error(w, `{"error":"Не удалось обновить дату задачи"}`, http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(map[string]string{}); err != nil {
		http.Error(w, `{"error":"Не удалось закодировать ответ"}`, http.StatusInternalServerError)
	}
}
