package tapdb

import "fmt"

func (m *mySQLDB) makeTables() error {
	_, errCreateDB := m.conn.Exec(`
		CREATE DATABASE IF NOT EXISTS tapestry
		DEFAULT CHARACTER SET = 'utf8'
		DEFAULT COLLATE 'utf8_general_ci';`)
	if errCreateDB != nil {
		return fmt.Errorf("Could not create database: %v", errCreateDB)
	}

	_, errUse := m.conn.Exec(`USE tapestry`)
	if errUse != nil {
		return fmt.Errorf("Could not `USE` database: %v", errUse)
	}

	_, errCreatePersonteam := m.conn.Exec(`
		CREATE TABLE IF NOT EXISTS personteams (
			email  VARCHAR(255) NOT NULL,
			name   VARCHAR(255) NOT NULL,
			abbrev VARCHAR(63) NOT NULL,
			colorf VARCHAR(63) NOT NULL,
			colorb VARCHAR(63) NOT NULL,
			PRIMARY KEY (email)
		);`)
	if errCreatePersonteam != nil {
		return fmt.Errorf("Could not create personteams table: %v", errCreatePersonteam)
	}
	return nil
}
