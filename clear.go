package tapdb

import (
	"fmt"
)

func (db *mysqlDB) ClearStks(domain string) {
	_, err := db.conn.Exec(fmt.Sprintf(`
	DELETE FROM stakeholders
	WHERE       domain = '%v'
	;`, domain))
	if err != nil {
		panic(err)
	}
}

func (db *mysqlDB) ClearStkHierLinks(domain string) {
	_, err := db.conn.Exec(fmt.Sprintf(`
	DELETE FROM stakeholders_hierarchy
	WHERE       domain = '%v'
	;`, domain))
	if err != nil {
		panic(err)
	}
}

func (db *mysqlDB) ClearThreads(domain string) {
	_, err := db.conn.Exec(fmt.Sprintf(`
	DELETE FROM threads
	WHERE       domain = '%v'
	;`, domain))
	if err != nil {
		panic(err)
	}
}

func (db *mysqlDB) ClearThreadHierLinks(domain string) {
	_, err := db.conn.Exec(fmt.Sprintf(`
	DELETE FROM threads_hierarchy
	WHERE       domain = '%v'
	;`, domain))
	if err != nil {
		panic(err)
	}
}

func (db *mysqlDB) ClearThreadStkLinks(domain string) {
	_, err := db.conn.Exec(fmt.Sprintf(`
	DELETE FROM threads_stakeholders
	WHERE       domain = '%v'
	;`, domain))
	if err != nil {
		panic(err)
	}
}
