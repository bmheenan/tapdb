package tapdb

import (
	"database/sql"
	"fmt"
)

// Make the database if it doesn't already exist
func makeDB(conn *sql.DB, dbName string) {
	_, err := conn.Exec(fmt.Sprintf(`
		CREATE DATABASE IF NOT EXISTS %v
		DEFAULT CHARACTER SET = 'utf8'
		DEFAULT COLLATE 'utf8_general_ci';`, dbName))
	if err != nil {
		panic(fmt.Sprintf("Could not create database: %v", err))
	}
}
