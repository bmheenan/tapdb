package tapdb

import "fmt"

// Make the db tables if they don't already exist
func (db *mysqlDB) makeTables() error {
	_, errPT := db.conn.Exec(`
		CREATE TABLE IF NOT EXISTS personteams (
			email      VARCHAR(255) NOT NULL,
			domain     VARCHAR(255) NOT NULL,
			name       VARCHAR(255) NOT NULL,
			abbrev     VARCHAR(63)  NOT NULL,
			colorf     VARCHAR(63)  NOT NULL,
			colorb     VARCHAR(63)  NOT NULL,
			itertiming VARCHAR(63)  NOT NULL,
			PRIMARY KEY (email),
			INDEX (domain)
		);`)
	if errPT != nil {
		return fmt.Errorf("Could not create personteams table: %v", errPT)
	}

	_, errPTPC := db.conn.Exec(`
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
	if errPTPC != nil {
		return fmt.Errorf("Could not create personteams_parent_child table: %v", errPTPC)
	}

	_, errTh := db.conn.Exec(`
		CREATE TABLE IF NOT EXISTS threads (
			id          INT              NOT NULL AUTO_INCREMENT,
			name        VARCHAR(255)     NOT NULL,
			domain      VARCHAR(255)     NOT NULL,
			description TEXT(65535),
			owner       VARCHAR(255)     NOT NULL,
			iteration   VARCHAR(255)     NOT NULL,
			state       VARCHAR(255)     NOT NULL,
			percentile  DOUBLE PRECISION NOT NULL,
			costdirect  INT              NOT NULL,
			costtotal   INT              NOT NULL,
			PRIMARY KEY (id),
			FOREIGN KEY (owner) REFERENCES personteams(email),
			INDEX (owner),
			INDEX (iteration),
			INDEX (state)
		);`)
	if errTh != nil {
		return fmt.Errorf("Could not create threads table: %v", errTh)
	}

	_, errThPC := db.conn.Exec(`
		CREATE TABLE IF NOT EXISTS threads_parent_child (
			parent    INT          NOT NULL,
			child     INT          NOT NULL,
			domain    VARCHAR(255) NOT NULL,
			iteration VARCHAR(255) NOT NULL,
			ord       INT          NOT NULL,
			PRIMARY KEY (parent, child),
			FOREIGN KEY (parent) REFERENCES threads(id),
			FOREIGN KEY (child) REFERENCES threads(id),
			INDEX (parent),
			INDEX (child),
			INDEX (domain)
		);`)
	if errThPC != nil {
		return fmt.Errorf("Could not create threads parent/child table: %v", errThPC)
	}

	_, errSk := db.conn.Exec(`
		CREATE TABLE IF NOT EXISTS threads_stakeholders (
			thread      INT          NOT NULL,
			stakeholder VARCHAR(255) NOT NULL,
			domain      VARCHAR(255) NOT NULL,
			iteration   VARCHAR(255) NOT NULL,
			ord         INT          NOT NULL,
			toplevel    BOOL         NOT NULL,
			costctx     INT          NOT NULL,
			PRIMARY KEY (thread, stakeholder),
			FOREIGN KEY (thread) REFERENCES threads(id),
			FOREIGN KEY (stakeholder) REFERENCES personteams(email),
			INDEX (thread),
			INDEX (stakeholder),
			INDEX (domain)
		);`)
	if errSk != nil {
		return fmt.Errorf("Could not create thread stakeholders table: %v", errSk)
	}

	_, errSkPC := db.conn.Exec(`
		CREATE TABLE IF NOT EXISTS threads_stakeholders_parent_child (
			parent      INT          NOT NULL,
			child       INT          NOT NULL,
			stakeholder VARCHAR(255) NOT NULL,
			domain      VARCHAR(255) NOT NULL,
			PRIMARY KEY (parent, child, stakeholder),
			FOREIGN KEY (parent) REFERENCES threads(id),
			FOREIGN KEY (child) REFERENCES threads(id),
			FOREIGN KEY (stakeholder) REFERENCES personteams(email),
			INDEX (parent),
			INDEX (child),
			INDEX (stakeholder),
			INDEX (domain)
		);`)
	if errSkPC != nil {
		return fmt.Errorf("Could not create thread stakeholders parent child table: %v", errSkPC)
	}
	return nil
}
