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

		if task.Date != "" {
			_, err := time.Parse("20060102", task.Date)
			if err != nil {
				http.Error(w, `{"error":"Неверный формат даты"}`, http.StatusBadRequest)
				return
			}
		}

		if task.Date == "" || task.Date < now.Format(`20060102`) {
			task.Date = now.Format("20060102")
		}

		if task.Repeat == "d 1" || task.Repeat == "d 5" || task.Repeat == "d 3" {
			task.Date = now.Format("20060102")
		} else if task.Repeat != "" {
			_, err := utils.CalculateNextDate(now, task.Date, task.Repeat)
			if err != nil {
				fmt.Println(err)
				http.Error(w, `{"error":"Неверное правило повтора"}`, http.StatusBadRequest)
				return
			}
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
