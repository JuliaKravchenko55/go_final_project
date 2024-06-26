package config

import (
	"log"
	"os"
	"path/filepath"
)

func GetDBFilePath() string {
	dbFilePath := os.Getenv("TODO_DBFILE")
	if dbFilePath == "" {
		currentDir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		dbFilePath = filepath.Join(currentDir, "scheduler.db")
	}
	return dbFilePath
}
