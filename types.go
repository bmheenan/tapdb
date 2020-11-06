package tapdb

import (
	"database/sql"
	"errors"

	taps "github.com/bmheenan/taps"
)

// DBInterface defines which functions can be executued or queried against the database
type DBInterface interface {
	ClearPersonteams(domain string) error
	ClearPersonteamsPC(domain string) error
	ClearThreads(domain string) error
	ClearThreadsPC(domain string) error
	ClearStakeholders(domain string) error
	ClearThreadsStakeholdersPC(domain string) error

	NewPersonteam(email, domain, name, abbrev, colorf, colorb string, itertiming taps.IterTiming) error
	LinkPersonteams(parent, child, domain string) error
	GetPersonteam(email string) (*taps.Personteam, error)
	GetPersonteamDescendants(email string) (map[string](*taps.Personteam), error)

	NewThread(name, domain, owner, iteration, state string, percentile float64, cost int) (int64, error)
	LinkThreads(parent, child int64, iter string, ord int, domain string) error
	LinkThreadsStakeholder(parent, child int64, stakeholder, domain string) error
	GetThreadOrderBefore(parent int64, iter string, order int) (int, error)

	GetThreadrel(id int64, stakeholder string) (*taps.Threadrel, error)
	GetThreadDescendants(id int64, stakeholder string) (map[int64](*taps.Threadrel), error)
	GetThreadAncestors(id int64, stakeholder string) (map[int64](*taps.Threadrel), error)
	GetChildThreadsSkIter(threads []int64, stakeholder, iteration string) (map[int64](*taps.Threadrel), error)
	GetParentThreadsSkIter(threads []int64, stakeholder, iteration string) (map[int64](*taps.Threadrel), error)

	NewStakeholder(thread int64, stakeholder, domain, iter string, ord int, topLvl bool, cost int) error
	GetStakeholderAncestors(thread int64) (map[string]*taps.Personteam, error)
	GetStakeholderDescendants(thread int64) (map[string]*taps.Personteam, error)
	GetStakeholderOrderBefore(stakeholder, iter string, order int) (int, error)
}

// ErrNotFound indicates that no matching record was found when querying
var ErrNotFound = errors.New("Not found")

// ErrBadArgs indicates that the arguments given to the function were bad
var ErrBadArgs = errors.New("Bad arguments")

type mysqlDB struct {
	conn *sql.DB
}
