package database

import (
	"database/sql"
	"github.com/JuliaKravchenko55/go_final_project/internal/config"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Initialize() {
	dbFilePath := config.GetDBFilePath()

	var err error
	DB, err = sql.Open("sqlite", dbFilePath)
	if err != nil {
		log.Fatal(err)
	}

	createTableSQL := `
    CREATE TABLE IF NOT EXISTS scheduler (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        date TEXT NOT NULL,
        title TEXT NOT NULL,
        comment TEXT,
        repeat TEXT CHECK(length(repeat) <= 128)
    );`

	if _, err = DB.Exec(createTableSQL); err != nil {
		log.Fatal(err)
	}

	createIndexSQL := `CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);`
	if _, err = DB.Exec(createIndexSQL); err != nil {
		log.Fatal(err)
	}
}
