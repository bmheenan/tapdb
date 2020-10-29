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
	IterationsByPersonteam(string) ([]string, error)
	NewThread(*tapstruct.Threaddetail, []*tapstruct.Threadrow, []*tapstruct.Threadrow) (int64, error)
	GetThreadrowsByPersonteamPlan(string, []string) ([]tapstruct.Threadrow, error)
	NewStakeholder(int64, *tapstruct.Personteam) error
	MoveThread(*tapstruct.Threadrow, BeforeAfter, *tapstruct.Threadrow) error
}

// MySQL implementation of DBInterface
type mySQLDB struct {
	conn  *sql.DB
	stmts map[string](*sql.Stmt)
}

var _ mysql.Config // keep the linter from cleaning up the import. sql.Open needs it

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
		cv.dbName = "tapestry"
	}
	db := &mySQLDB{}
	var errO error
	db.conn, errO = sql.Open("mysql", cv.formatName())
	if errO != nil {
		return &mySQLDB{}, fmt.Errorf("Could not establish mysql db connection: %v", errO)
	}
	if db.conn.Ping() == driver.ErrBadConn {
		return &mySQLDB{}, errors.New("Could not establish a good connection to db: ping returned bad connection")
	}
	errMk := db.makeTables()
	if errMk != nil {
		return &mySQLDB{}, fmt.Errorf("Could not make db tables: %v", errMk)
	}
	db.stmts = make(map[string](*sql.Stmt))
	initFuncs := []func() error{
		// Every init for each exported function must be added here
		db.initGetPersonteam,
		db.initNewPersonteam,
		db.initClearDomain,
		db.initIterationsByPersonteam,
		db.initGetThreadrowsByPersonteamPlan,
		db.initNewThread,
		db.initNewStakeholder,
		db.initThreadMove,
		db.initCalibrateOrdPct,
	}
	for _, f := range initFuncs {
		initErr := f()
		if initErr != nil {
			return &mySQLDB{}, fmt.Errorf("Could not initialize function: %v", initErr)
		}
	}
	return db, nil
}
