package tapdb

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"os"

	"github.com/bmheenan/tapstruct"
	"github.com/go-sql-driver/mysql"
)

// ErrNotFound is wrapped and returned in cases where a query has no matches
var ErrNotFound = errors.New("Not found")

// DBInterface defines which actions can be taken against the database
type DBInterface interface {
	NewPersonteam(*tapstruct.Personteam, string) error
	GetPersonteam(string, int) (*tapstruct.Personteam, error)
	ClearDomain(string) error
}

// MySQL implementation of DBInterface
type mySQLDB struct {
	conn  *sql.DB
	stmts map[string](*sql.Stmt)
}

var _ mysql.Config // keep the compiler from cleaning up the import. sql.Open needs it

// InitDB initializes a db connection and returns a DBInterface with the available methods
func InitDB() (DBInterface, error) {
	var cv = &connVars{}
	cv.user, cv.pass, cv.dbName = getCredentials()
	if os.Getenv("GAE_INSTANCE") != "" {
		// Running in prod
		cv.unixSocket = "/cloudsql/"
	} else {
		// Running locally. Config for the Cloud SQL proxy
		// https://cloud.google.com/sql/docs/mysql/quickstart-proxy-test
		cv.host = "localhost"
		cv.port = "3306"
		cv.dbName = ""
	}
	db := &mySQLDB{}
	var err1 error
	db.conn, err1 = sql.Open("mysql", cv.formatName())
	if err1 != nil {
		return &mySQLDB{}, fmt.Errorf("Could not establish mysql db connection: %v", err1)
	}
	if db.conn.Ping() == driver.ErrBadConn {
		return &mySQLDB{}, errors.New("Could not establish a good connection to db: ping returned bad connection")
	}
	err2 := db.makeTables()
	if err2 != nil {
		return &mySQLDB{}, fmt.Errorf("Could not make db tables: %v", err2)
	}
	db.stmts = make(map[string](*sql.Stmt))
	initFuncs := []struct {
		key string
		f   func() error
	}{
		// Every init for each exported function must be added here
		{keyGetPersonteam, db.initGetPersonteam},
		{keyNewPersonteam, db.initNewPersonteam},
		{keyClearDomainPT, db.initClearDomain},
	}
	for _, v := range initFuncs {
		initErr := v.f()
		if initErr != nil {
			return &mySQLDB{}, fmt.Errorf("Could not initialize function %s: %v", v.key, initErr)
		}
	}
	return db, nil
}
