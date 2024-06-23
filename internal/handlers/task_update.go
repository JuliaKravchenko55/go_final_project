package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/JuliaKravchenko55/go_final_project/internal/models"
	"github.com/JuliaKravchenko55/go_final_project/internal/store"
	"github.com/JuliaKravchenko55/go_final_project/internal/utils"
)

func UpdateTask(store *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		if err := store.UpdateTask(&task); err != nil {
			if err.Error() == "task not found" {
				http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
			} else {
				http.Error(w, `{"error":"Не удалось обновить задачу"}`, http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(map[string]string{}); err != nil {
			http.Error(w, `{"error":"Не удалось закодировать ответ"}`, http.StatusInternalServerError)
		}
	}
}
