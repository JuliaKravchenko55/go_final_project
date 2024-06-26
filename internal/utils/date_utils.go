package utils

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func CalculateNextDate(now time.Time, date string, repeat string) (string, error) {
	startDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", err
	}

	if repeat == "" {
		return "", errors.New("Правило повтора не может быть пустым")
	}

	repeatParts := strings.Split(repeat, " ")

	switch repeatParts[0] {
	case "d":
		if len(repeatParts) != 2 {
			return "", errors.New("Недопустимый формат правила повтора")
		}

		days, err := strconv.Atoi(repeatParts[1])
		if err != nil || days > 400 {
			return "", errors.New("Недопустимое количество дней")
		}

		for {
			startDate = startDate.AddDate(0, 0, days)
			if startDate.After(now) {
				break
			}
		}

	case "y":
		for {
			startDate = startDate.AddDate(1, 0, 0)
			if startDate.After(now) {
				break
			}
		}

	case "w":
		if len(repeatParts) != 2 {
			return "", errors.New("Недопустимый формат правила повтора")
		}

		daysOfWeek := strings.Split(repeatParts[1], ",")
		dayInts := make([]int, len(daysOfWeek))
		for i, day := range daysOfWeek {
			dayInt, err := strconv.Atoi(day)
			if err != nil || dayInt < 1 || dayInt > 7 {
				return "", errors.New("Недопустимый день недели")
			}
			dayInts[i] = dayInt
		}

		for {
			startDate = startDate.AddDate(0, 0, 1)
			for _, day := range dayInts {
				if int(startDate.Weekday()) == day%7 {
					if startDate.After(now) {
						return startDate.Format("20060102"), nil
					}
				}
			}
		}

	case "m":
		if len(repeatParts) < 2 {
			return "", errors.New("Недопустимый формат правила повтора")
		}

		daysOfMonth := strings.Split(repeatParts[1], ",")
		dayInts := make([]int, len(daysOfMonth))
		for i, day := range daysOfMonth {
			dayInt, err := strconv.Atoi(day)
			if err != nil || (dayInt < 1 && dayInt != -1 && dayInt != -2) || dayInt > 31 {
				return "", errors.New("Недопустимый день месяца")
			}
			dayInts[i] = dayInt
		}

		monthsOfYear := []int{}
		if len(repeatParts) > 2 {
			months := strings.Split(repeatParts[2], ",")
			for _, month := range months {
				monthInt, err := strconv.Atoi(month)
				if err != nil || monthInt < 1 || monthInt > 12 {
					return "", errors.New("Недопустимый месяц")
				}
				monthsOfYear = append(monthsOfYear, monthInt)
			}
		}

		for {
			startDate = startDate.AddDate(0, 0, 1)
			day := startDate.Day()
			month := int(startDate.Month())

			validDay := false
			for _, d := range dayInts {
				if d == day || (d == -1 && day == daysInMonth(startDate)) || (d == -2 && day == daysInMonth(startDate)-1) {
					validDay = true
					break
				}
			}

			validMonth := len(monthsOfYear) == 0
			for _, m := range monthsOfYear {
				if m == month {
					validMonth = true
					break
				}
			}

			if validDay && validMonth && startDate.After(now) {
				return startDate.Format("20060102"), nil
			}
		}

	default:
		return "", errors.New("Неподдерживаемый формат правила повтора")
	}

	return startDate.Format("20060102"), nil
}

func daysInMonth(t time.Time) int {
	return time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, t.Location()).Day()
}
