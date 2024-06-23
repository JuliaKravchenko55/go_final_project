package store

import (
	"database/sql"
	"errors"
	"time"

	"github.com/JuliaKravchenko55/go_final_project/internal/models"
	"github.com/JuliaKravchenko55/go_final_project/internal/utils"
)

type Store struct {
	DB *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{DB: db}
}

func (s *Store) GetTaskByID(id string) (*models.Task, error) {
	row := s.DB.QueryRow(`SELECT * FROM scheduler WHERE id = ?`, id)
	var task models.Task
	if err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("task not found")
		}
		return nil, err
	}
	return &task, nil
}

func (s *Store) ListTasks(search string) ([]models.Task, error) {
	var rows *sql.Rows
	var err error

	if search == "" {
		rows, err = s.DB.Query(`SELECT * FROM scheduler ORDER BY date ASC LIMIT 50`)
	} else {
		if _, dateErr := time.Parse("02.01.2006", search); dateErr == nil {
			searchDate := time.Now().Format("20060102")
			rows, err = s.DB.Query(`SELECT * FROM scheduler WHERE date = ? ORDER BY date ASC LIMIT 50`, searchDate)
		} else {
			rows, err = s.DB.Query(`SELECT * FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date ASC LIMIT 50`, "%"+search+"%", "%"+search+"%")
		}
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]models.Task, 0)
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *Store) DeleteTaskByID(id string) error {
	_, err := s.DB.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	return err
}

func (s *Store) UpdateTaskDate(id string, nextDate string) error {
	_, err := s.DB.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`, nextDate, id)
	return err
}

func (s *Store) CreateTask(task *models.Task) (int64, error) {
	res, err := s.DB.Exec(
		`INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`,
		task.Date, task.Title, task.Comment, task.Repeat,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (s *Store) UpdateTask(task *models.Task) error {
	res, err := s.DB.Exec(
		`UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`,
		task.Date, task.Title, task.Comment, task.Repeat, task.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("task not found")
	}

	return nil
}

func (s *Store) CompleteTask(id string) error {
	task, err := s.GetTaskByID(id)
	if err != nil {
		return err
	}

	if task.Repeat == "" {
		if err := s.DeleteTaskByID(id); err != nil {
			return err
		}
	} else {
		nextDate, err := utils.CalculateNextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			return err
		}

		if err := s.UpdateTaskDate(id, nextDate); err != nil {
			return err
		}
	}

	return nil
}
