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

	/*_, errUse := db.conn.Exec(`USE tapestry`)
	if errUse != nil {
		return fmt.Errorf("Could not `USE` database: %v", errUse)
	}*/

	_, errCreatePersonteam := db.conn.Exec(`
		CREATE TABLE IF NOT EXISTS personteams (
			email           VARCHAR(255) NOT NULL,
			domain          VARCHAR(255) NOT NULL,
			name            VARCHAR(255) NOT NULL,
			abbrev          VARCHAR(63)  NOT NULL,
			colorf          VARCHAR(63)  NOT NULL,
			colorb          VARCHAR(63)  NOT NULL,
			haschildren     BOOLEAN      NOT NULL,
			iterationTiming VARCHAR(63)  NOT NULL,
			PRIMARY KEY (email),
			INDEX (domain)
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
			FOREIGN KEY (parent) REFERENCES personteams(email),
			FOREIGN KEY (child) REFERENCES personteams(email),
			INDEX (parent),
			INDEX (child),
			INDEX (domain)
		);`)
	if errCreatePersonteamPC != nil {
		return fmt.Errorf("Could not create personteams_parent_child table: %v", errCreatePersonteamPC)
	}

	_, errCreateThreads := db.conn.Exec(`
		CREATE TABLE IF NOT EXISTS threads (
			id          INT            NOT NULL AUTO_INCREMENT,
			name        VARCHAR(255)   NOT NULL,
			domain      VARCHAR(255)   NOT NULL,
			description TEXT(65535),
			owner       VARCHAR(255)   NOT NULL,
			iteration   VARCHAR(255)   NOT NULL,
			state       VARCHAR(255)   NOT NULL,
			ord         INT            NOT NULL,
			percentile  INT            NOT NULL,
			costdirect  INT            NOT NULL,
			costtotal   INT            NOT NULL,
			PRIMARY KEY (id),
			FOREIGN KEY (owner) REFERENCES personteams(email),
			INDEX (owner),
			INDEX (iteration),
			INDEX (state)
		);`)
	if errCreateThreads != nil {
		return fmt.Errorf("Could not create threads table: %v", errCreateThreads)
	}

	_, errThreadsPC := db.conn.Exec(`
		CREATE TABLE IF NOT EXISTS threads_parent_child (
			parent INT          NOT NULL,
			child  INT          NOT NULL,
			domain VARCHAR(255) NOT NULL,
			PRIMARY KEY (parent, child),
			FOREIGN KEY (parent) REFERENCES threads(id),
			FOREIGN KEY (child) REFERENCES threads(id),
			INDEX (parent),
			INDEX (child),
			INDEX (domain)
		);`)
	if errThreadsPC != nil {
		return fmt.Errorf("Could not create threads parent/child table: %v", errThreadsPC)
	}

	_, errStakeholders := db.conn.Exec(`
		CREATE TABLE IF NOT EXISTS threads_stakeholders (
			thread      INT          NOT NULL,
			stakeholder VARCHAR(255) NOT NULL,
			domain      VARCHAR(255) NOT NULL,
			PRIMARY KEY (thread, stakeholder),
			FOREIGN KEY (thread) REFERENCES threads(id),
			FOREIGN KEY (stakeholder) REFERENCES personteams(email),
			INDEX (thread),
			INDEX (stakeholder),
			INDEX (domain)
		);`)
	if errStakeholders != nil {
		return fmt.Errorf("Could not create thread stakeholders table: %v", errStakeholders)
	}
	return nil
}
