package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/JuliaKravchenko55/go_final_project/internal/database"
	"github.com/JuliaKravchenko55/go_final_project/internal/models"
	"github.com/JuliaKravchenko55/go_final_project/internal/utils"
)

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, `{"error":"Неверный формат запроса"}`, http.StatusBadRequest)
		return
	}

	if task.ID == 0 {
		http.Error(w, `{"error":"ID задачи обязателен"}`, http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		http.Error(w, `{"error":"Заголовок задачи обязателен"}`, http.StatusBadRequest)
		return
	}

	if task.Date != "" {
		if _, err := time.Parse("20060102", task.Date); err != nil {
			http.Error(w, `{"error":"Неверный формат даты"}`, http.StatusBadRequest)
			return
		}
	}

	if task.Date == "" || task.Date < time.Now().Format("20060102") {
		task.Date = time.Now().Format("20060102")
	}

	if task.Repeat != "" {
		date, err := utils.CalculateNextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			http.Error(w, `{"error":"Неверное правило повтора"}`, http.StatusBadRequest)
			return
		}
		task.Date = date
	}

	res, err := database.DB.Exec(
		`UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`,
		task.Date, task.Title, task.Comment, task.Repeat, task.ID,
	)
	if err != nil {
		http.Error(w, `{"error":"Не удалось обновить задачу"}`, http.StatusInternalServerError)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, `{"error":"Не удалось получить количество обновленных строк"}`, http.StatusInternalServerError)
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
