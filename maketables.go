package tapdb

import "fmt"

func (db *mySQLDB) makeTables() error {
	_, errCreateDB := db.conn.Exec(`
		CREATE DATABASE IF NOT EXISTS tapestry
		DEFAULT CHARACTER SET = 'utf8'
		DEFAULT COLLATE 'utf8_general_ci';`)
	if errCreateDB != nil {
		return fmt.Errorf("Could not create database: %v", errCreateDB)
	}

	_, errUse := db.conn.Exec(`USE tapestry`)
	if errUse != nil {
		return fmt.Errorf("Could not `USE` database: %v", errUse)
	}

	_, errCreatePersonteam := db.conn.Exec(`
		CREATE TABLE IF NOT EXISTS personteams (
			email       VARCHAR(255) NOT NULL,
			domain      VARCHAR(255) NOT NULL,
			name        VARCHAR(255) NOT NULL,
			abbrev      VARCHAR(63)  NOT NULL,
			colorf      VARCHAR(63)  NOT NULL,
			colorb      VARCHAR(63)  NOT NULL,
			haschildren BOOLEAN      NOT NULL,
			PRIMARY KEY (email, domain),
			INDEX(email, domain)
		);`)
	if errCreatePersonteam != nil {
		return fmt.Errorf("Could not create personteams table: %v", errCreatePersonteam)
	}

	_, errCreatePersonteamPC := db.conn.Exec(`
		CREATE TABLE IF NOT EXISTS personteams_parent_child (
			parent VARCHAR(255) NOT NULL,
			child  VARCHAR(255) NOT NULL,
			domain VARCHAR(255) NOT NULL,
			PRIMARY KEY (parent, child),
			INDEX(parent),
			INDEX(child)
		);`)
	if errCreatePersonteamPC != nil {
		return fmt.Errorf("Could not create personteams_parent_child table: %v", errCreatePersonteamPC)
	}
	return nil
}
