package tapdb

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
)

var _ mysql.Config // keep the linter from cleaning up the import. sql.Open needs it

// Init initialized the database connection, creates the database and tables if needed, then returns an interface
// with all functions that can be executed against the database
// The username, password, and instance name of the database must be provided
func Init(user, pass string) (DBInterface, error) {
	var (
		unixSocket,
		host,
		port string
		dbName = "tapestry"
	)
	if os.Getenv("GAE_INSTANCE") != "" {
		// Running in prod
		unixSocket = "/cloudsql/"
	} else {
		// Running locally. Config for the Cloud SQL proxy
		// https://cloud.google.com/sql/docs/mysql/quickstart-proxy-test
		host = "localhost"
		port = "3306"
	}
	db := &mysqlDB{}
	// First connection has no db specified so it doesn't have to exist; we'll create it with this connection
	crDBconn, errOp0 := sql.Open("mysql", fmtOpenStr(user, pass, host, port, unixSocket, ""))
	if errOp0 != nil {
		return &mysqlDB{}, fmt.Errorf("Could not establish mysql db connection: %v", errOp0)
	}
	if crDBconn.Ping() == driver.ErrBadConn {
		return &mysqlDB{}, errors.New("Could not establish a good connection to db: ping returned bad connection")
	}
	errMkDB := makeDB(crDBconn, dbName)
	if errMkDB != nil {
		return &mysqlDB{}, fmt.Errorf("Could not make database: %v", errMkDB)
	}
	crDBconn.Close()
	// Second persistent connection specifies the DB now that we just created it. Subsequent queries will already be
	// using the db this way
	var errOp1 error
	db.conn, errOp1 = sql.Open("mysql", fmtOpenStr(user, pass, host, port, unixSocket, dbName))
	if errOp1 != nil {
		return &mysqlDB{}, fmt.Errorf("Could not establish mysql db connection: %v", errOp1)
	}
	if db.conn.Ping() == driver.ErrBadConn {
		return &mysqlDB{}, errors.New("Could not establish a good connection to db: ping returned bad connection")
	}
	errMkTbls := db.makeTables()
	if errMkTbls != nil {
		return &mysqlDB{}, fmt.Errorf("Could not make db tables: %v", errMkTbls)
	}
	return db, nil
}

// Format the connection string for both prod and local dev
func fmtOpenStr(user, pass, host, port, unixSocket, dbName string) string {
	var cred string
	if user != "" {
		cred = user
		if pass != "" {
			cred = cred + ":" + pass
		}
		cred = cred + "@"
	}
	if unixSocket != "" {
		return fmt.Sprintf("%sunix(%s)/%s", cred, unixSocket, dbName)
	}
	return fmt.Sprintf("%stcp([%s]:%s)/%s", cred, host, port, dbName)
}
