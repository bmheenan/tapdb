package tapdb

import (
	"fmt"
)

func (db *mysqlDB) ClearStks(domain string) error {
	_, err := db.conn.Exec(fmt.Sprintf(`
	DELETE FROM stakeholders
	WHERE       domain = '%v'
	;`, domain))
	return err
}

func (db *mysqlDB) ClearStkHierLinks(domain string) error {
	_, err := db.conn.Exec(fmt.Sprintf(`
	DELETE FROM stakeholders_hierarchy
	WHERE       domain = '%v'
	;`, domain))
	return err
}

func (db *mysqlDB) ClearThreads(domain string) error {
	_, err := db.conn.Exec(fmt.Sprintf(`
	DELETE FROM threads
	WHERE       domain = '%v'
	;`, domain))
	return err
}

func (db *mysqlDB) ClearThreadHierLinks(domain string) error {
	_, err := db.conn.Exec(fmt.Sprintf(`
	DELETE FROM threads_hierarchy
	WHERE       domain = '%v'
	;`, domain))
	return err
}

func (db *mysqlDB) ClearThreadStkLinks(domain string) error {
	_, err := db.conn.Exec(fmt.Sprintf(`
	DELETE FROM threads_stakeholders
	WHERE       domain = '%v'
	;`, domain))
	return err
}

func (db *mysqlDB) ClearThreadStkHierLinks(domain string) error {
	_, err := db.conn.Exec(fmt.Sprintf(`
	DELETE FROM threads_stakeholders_hierarchy
	WHERE       domain = '%v'
	;`, domain))
	return err
}
