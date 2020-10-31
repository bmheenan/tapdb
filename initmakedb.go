package tapdb

import (
	"database/sql"
	"fmt"
)

// Make the database if it doesn't already exist
func makeDB(conn *sql.DB, dbName string) error {
	_, errMkDB := conn.Exec(fmt.Sprintf(`
		CREATE DATABASE IF NOT EXISTS %v
		DEFAULT CHARACTER SET = 'utf8'
		DEFAULT COLLATE 'utf8_general_ci';`, dbName))
	if errMkDB != nil {
		return fmt.Errorf("Could not create database: %v", errMkDB)
	}
	return nil
}
