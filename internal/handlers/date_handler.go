package handlers

import (
	"net/http"
	"time"

	"github.com/JuliaKravchenko55/go_final_project/internal/utils"
)

func CalculateNextDate(w http.ResponseWriter, r *http.Request) {
	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeatStr := r.URL.Query().Get("repeat")

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "Недопустимый параметр 'now'", http.StatusBadRequest)
		return
	}

	nextDate, err := utils.CalculateNextDate(now, dateStr, repeatStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(nextDate))
}
