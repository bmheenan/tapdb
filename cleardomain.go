package tapdb

import (
	"fmt"
)

// ClearPersonteams deletes all personteams of the matching domain
func (db *mysqlDB) ClearPersonteams(domain string) error {
	_, err := db.conn.Exec(fmt.Sprintf(`
	DELETE FROM personteams
	WHERE       domain = '%v'
	;`, domain))
	return err
}

// ClearPersonteamsPC deletes all personteams_parent_child relationships for the matching domain
func (db *mysqlDB) ClearPersonteamsPC(domain string) error {
	_, err := db.conn.Exec(fmt.Sprintf(`
	DELETE FROM personteams_parent_child
	WHERE       domain = '%v'
	;`, domain))
	return err
}

// ClearThreads deletes all threads for the matching domain
func (db *mysqlDB) ClearThreads(domain string) error {
	_, err := db.conn.Exec(fmt.Sprintf(`
	DELETE FROM threads
	WHERE       domain = '%v'
	;`, domain))
	return err
}

// ClearThreadsPC deletes all threads_parent_child relationships for the matching domain
func (db *mysqlDB) ClearThreadsPC(domain string) error {
	_, err := db.conn.Exec(fmt.Sprintf(`
	DELETE FROM threads_parent_child
	WHERE       domain = '%v'
	;`, domain))
	return err
}

// ClearStakeholders deletes all threads_stakeholders relationships for the matching domain
func (db *mysqlDB) ClearStakeholders(domain string) error {
	_, err := db.conn.Exec(fmt.Sprintf(`
	DELETE FROM threads_stakeholders
	WHERE       domain = '%v'
	;`, domain))
	return err
}

// ClearThreadsStakholdersPC deletes all threads_stakeholders_parent_child relationiships for the matching domain
func (db *mysqlDB) ClearThreadsStakeholdersPC(domain string) error {
	_, err := db.conn.Exec(fmt.Sprintf(`
	DELETE FROM threads_stakeholders_parent_child
	WHERE       domain = '%v'
	;`, domain))
	return err
}
