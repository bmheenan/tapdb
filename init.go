package tapdb

import (
	"database/sql"
	"fmt"
)

type tapdb struct {
	db *sql.DB
}

func (t *tapdb) Connect() {
	var (
		socketDir = "/cloudsql"
		err       error
	)
	user, pass, conn, dbname := getCredentials()
	dbURI := fmt.Sprintf("user=%s password=%s database=%s host=%s/%s", user, pass, dbname, socketDir, conn)
	t.db, err = sql.Open("pgx", dbURI)
	if err != nil {
		panic(fmt.Errorf("Could not open SQL db: %v", err))
	}
}

var ensureTableExistsQuery = `
CREATE TABLE IF NOT EXISTS personteams (
	email VARCHAR(255),
	name VARCHAR(255),
	abbrev VARCHAR(63),
	colorf VARCHAR(63),
	colorb VARCHAR(63),
)
`

func (t *tapdb) EnsureDBTables() {

}
