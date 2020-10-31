package tapdb

import (
	"database/sql"
	"errors"

	taps "github.com/bmheenan/taps"
)

// ErrNotFound indicates that no matching record was found when querying
var ErrNotFound = errors.New("Not found")

// ErrBadArgs indicates that the arguments given to the function were bad
var ErrBadArgs = errors.New("Bad arguments")

// DBInterface defines which functions can be executued or queried against the database
type DBInterface interface {
	ClearPersonteams(domain string) error
	ClearPersonteamsPC(domain string) error
	ClearThreads(domain string) error
	ClearThreadsPC(domain string) error
	ClearStakeholders(domain string) error

	NewPersonteam(email, domain, name, abbrev, colorf, colorb string, itertiming taps.IterTiming) error
	NewPersonteamPC(parent, child, domain string) error
	GetPersonteam(email string) (*taps.Personteam, error)
	GetPersonteamDescendants(email string) (map[string](*taps.Personteam), error)
}

type mysqlDB struct {
	conn *sql.DB
}
