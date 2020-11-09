package tapdb

import "fmt"

// Make the db tables if they don't already exist
func (db *mysqlDB) makeTables() error {
	_, errStk := db.conn.Exec(`
	CREATE TABLE IF NOT EXISTS stakeholders (
		email   VARCHAR(255) NOT NULL,
		domain  VARCHAR(255) NOT NULL,
		name    VARCHAR(255) NOT NULL,
		abbrev  VARCHAR(63)  NOT NULL,
		colorf  VARCHAR(63)  NOT NULL,
		colorb  VARCHAR(63)  NOT NULL,
		cadence VARCHAR(63)  NOT NULL,
		PRIMARY KEY (email),
		INDEX (domain)
	);`)
	if errStk != nil {
		return fmt.Errorf("Could not create stakeholders table: %v", errStk)
	}

	_, errStkH := db.conn.Exec(`
	CREATE TABLE IF NOT EXISTS stakeholders_hierarchy (
		parent VARCHAR(255) NOT NULL,
		child  VARCHAR(255) NOT NULL,
		domain VARCHAR(255) NOT NULL,
		PRIMARY KEY (parent, child),
		FOREIGN KEY (parent) REFERENCES stakeholders(email),
		FOREIGN KEY (child) REFERENCES stakeholders(email),
		INDEX (parent),
		INDEX (child),
		INDEX (domain)
	);`)
	if errStkH != nil {
		return fmt.Errorf("Could not create stakeholders_heirarchy table: %v", errStkH)
	}

	_, errTh := db.conn.Exec(`
	CREATE TABLE IF NOT EXISTS threads (
		id          INT              NOT NULL AUTO_INCREMENT,
		name        VARCHAR(255)     NOT NULL,
		domain      VARCHAR(255)     NOT NULL,
		description TEXT(65535),
		owner       VARCHAR(255)     NOT NULL,
		iter        VARCHAR(255)     NOT NULL,
		state       VARCHAR(255)     NOT NULL,
		percentile  DOUBLE PRECISION NOT NULL,
		costdir     INT              NOT NULL,
		costtot     INT              NOT NULL,
		PRIMARY KEY (id),
		FOREIGN KEY (owner) REFERENCES stakeholders(email),
		INDEX (owner),
		INDEX (iter),
		INDEX (state)
	);`)
	if errTh != nil {
		return fmt.Errorf("Could not create threads table: %v", errTh)
	}

	_, errThH := db.conn.Exec(`
	CREATE TABLE IF NOT EXISTS threads_hierarchy (
		parent INT          NOT NULL,
		child  INT          NOT NULL,
		domain VARCHAR(255) NOT NULL,
		iter   VARCHAR(255) NOT NULL,
		ord    INT          NOT NULL,
		PRIMARY KEY (parent, child),
		FOREIGN KEY (parent) REFERENCES threads(id),
		FOREIGN KEY (child) REFERENCES threads(id),
		INDEX (parent),
		INDEX (child),
		INDEX (domain),
		INDEX (iter)
	);`)
	if errThH != nil {
		return fmt.Errorf("Could not create threads_hierarchy table: %v", errThH)
	}

	_, errThStk := db.conn.Exec(`
		CREATE TABLE IF NOT EXISTS threads_stakeholders (
			thread INT          NOT NULL,
			stk    VARCHAR(255) NOT NULL,
			domain VARCHAR(255) NOT NULL,
			iter   VARCHAR(255) NOT NULL,
			ord    INT          NOT NULL,
			toplvl BOOL         NOT NULL,
			cost   INT          NOT NULL,
			PRIMARY KEY (thread, stk),
			FOREIGN KEY (thread) REFERENCES threads(id),
			FOREIGN KEY (stk) REFERENCES stakeholders(email),
			INDEX (thread),
			INDEX (stk),
			INDEX (domain),
			INDEX (iter)
		);`)
	if errThStk != nil {
		return fmt.Errorf("Could not create threads_stakeholders table: %v", errThStk)
	}

	_, errThStkH := db.conn.Exec(`
		CREATE TABLE IF NOT EXISTS threads_stakeholders_hierarchy (
			parent INT          NOT NULL,
			child  INT          NOT NULL,
			stk    VARCHAR(255) NOT NULL,
			domain VARCHAR(255) NOT NULL,
			PRIMARY KEY (parent, child, stk),
			FOREIGN KEY (parent) REFERENCES threads(id),
			FOREIGN KEY (child) REFERENCES threads(id),
			FOREIGN KEY (stk) REFERENCES stakeholders(email),
			INDEX (parent),
			INDEX (child),
			INDEX (stk),
			INDEX (domain)
		);`)
	if errThStkH != nil {
		return fmt.Errorf("Could not create threads_stakeholders_hierarchy table: %v", errThStkH)
	}
	return nil
}
