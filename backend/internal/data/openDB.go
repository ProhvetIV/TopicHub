package datahandler

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var Database *sql.DB
var databasePath = "../internal/data/database.sql"
var migrationsPath = "../internal/data/migrations/sqlite3/"

func OpenDb() {
	var err error
	Database, err = sql.Open("sqlite3", "../internal/data/data.db")
	if err != nil {
		log.Fatal("Failed to open the database:", err)
	}

	// Verify the connection to the database
	if err := Database.Ping(); err != nil {
		log.Fatal("Failed to load the database: ", err)
	}

	// Read and execute the sql schema file
	if err := executeSchema(databasePath); err != nil {
		log.Fatal("Failed to execute schema: ", err)
	}

	// Read and execute the migration files
	if err := executeMigrations(migrationsPath); err != nil {
		log.Fatal("Failed to execute migrations: ", err)
	}
}

// Read and execute the sql schema file
func executeSchema(dbPath string) error {
	sqlTable, err := os.ReadFile(dbPath)
	if err != nil {
		fmt.Println("error reading db")
		return err
	}

	_, err = Database.Exec(string(sqlTable))
	if err != nil {
		fmt.Println("error with execution")
		return err
	}
	return err
}

// Read and execute the migration files
func executeMigrations(dirPath string) error {
	// Open the directory
	f, err := os.Open(dirPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Read all the files
	files, err := f.ReadDir(-1)
	if err != nil {
		return err
	}

	// Iterate through the migrations and execute the _up.sql files
	for _, file := range files {
		if strings.Contains(file.Name(), ".up.sql") {
			migration, err := os.ReadFile(dirPath + "/" + file.Name())
			if err != nil {
				return err
			}

			_, err = Database.Exec(string(migration))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
