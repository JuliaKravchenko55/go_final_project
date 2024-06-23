package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/JuliaKravchenko55/go_final_project/internal/models"
	"github.com/JuliaKravchenko55/go_final_project/internal/store"
	"github.com/JuliaKravchenko55/go_final_project/internal/utils"
)

func CreateTask(store *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task models.Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, `{"error":"Неверный формат запроса"}`, http.StatusBadRequest)
			return
		}

		if task.Title == "" {
			http.Error(w, `{"error":"Заголовок задачи обязателен"}`, http.StatusBadRequest)
			return
		}

		now := time.Now()

		if task.Date == "" || task.Date == "today" {
			task.Date = now.Format("20060102")
		} else {
			parsedDate, err := time.Parse("20060102", task.Date)
			if err != nil {
				http.Error(w, `{"error":"Неверный формат даты"}`, http.StatusBadRequest)
				return
			}
			if parsedDate.Before(now.Truncate(24 * time.Hour)) {
				http.Error(w, `{"error":"Дата не может быть в прошлом"}`, http.StatusBadRequest)
				return
			}
		}

		if task.Repeat != "" {
			date, err := utils.CalculateNextDate(now, task.Date, task.Repeat)
			if err != nil {
				fmt.Println(err)
				http.Error(w, `{"error":"Неверное правило повтора"}`, http.StatusBadRequest)
				return
			}
			task.Date = date
		}

		id, err := store.CreateTask(&task)
		if err != nil {
			http.Error(w, `{"error":"Не удалось создать задачу"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(map[string]string{"id": fmt.Sprint(id)}); err != nil {
			http.Error(w, `{"error":"Не удалось закодировать ответ"}`, http.StatusInternalServerError)
		}
	}
}
