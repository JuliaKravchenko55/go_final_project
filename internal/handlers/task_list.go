package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/JuliaKravchenko55/go_final_project/internal/database"
	"github.com/JuliaKravchenko55/go_final_project/internal/models"
)

func ListTasks(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	var rows *sql.Rows
	var err error

	if search == "" {
		rows, err = database.DB.Query(`SELECT * FROM scheduler ORDER BY date ASC LIMIT 50`)
	} else {
		if _, dateErr := time.Parse("02.01.2006", search); dateErr == nil {
			searchDate := time.Now().Format("20060102")
			rows, err = database.DB.Query(`SELECT * FROM scheduler WHERE date = ? ORDER BY date ASC LIMIT 50`, searchDate)
		} else {
			rows, err = database.DB.Query(`SELECT * FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date ASC LIMIT 50`, "%"+search+"%", "%"+search+"%")
		}
	}

	if err != nil {
		http.Error(w, `{"error":"Не удалось получить задачи"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	tasks := make([]models.Task, 0)
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			http.Error(w, `{"error":"Не удалось сканировать задачу"}`, http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, `{"error":"Не удалось прочитать задачи"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(map[string][]models.Task{"tasks": tasks}); err != nil {
		http.Error(w, `{"error":"Не удалось закодировать задачи"}`, http.StatusInternalServerError)
	}
}
