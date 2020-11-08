package tapdb

import (
	"database/sql"
	"errors"

	taps "github.com/bmheenan/taps"
)

// DBInterface defines which functions can be executued or queried against the database
type DBInterface interface {

	// clear.go

	// ClearStks deletes all stakeholders in `domain`
	ClearStks(domain string) error
	// ClearStkHierLinks deletes all parent/child hierarchy links between stakeholders in `domain`
	ClearStkHierLinks(domain string) error
	// ClearThreads deletes all threads in `domain`
	ClearThreads(domain string) error
	// ClearThreadHierLinks deletes all parent/child hierarchy links in `domain`
	ClearThreadHierLinks(domain string) error
	// ClearThreadStkLinks deletes all relationships between threads and stakeholders in `domain`
	ClearThreadStkLinks(domain string) error
	// ClearThreadStkHierLinks deletes all parent/child hierarchy links between threads for all stakeholders in `domain`
	ClearThreadStkHierLinks(domain string) error

	// stakeholders.go

	// NewStk makes a new stakeholder with the given info
	NewStk(email, domain, name, abbrev, colorf, colorb string, cadence taps.Cadence) error
	// NewStkHierLink links two stakeholders as `parent` and `child` in `domain`
	NewStkHierLink(parent, child, domain string) error
	// GetStk gets the info for the stakeholder with `email`
	GetStk(email string) (*taps.Stakeholder, error)
	// GetStkDes gets all stakeholders that are descendants of the stakeholder with `email` (including itself)
	GetStkDes(email string) (map[string](*taps.Stakeholder), error)

	/*
		NewThread(name, domain, owner, iteration, state string, percentile float64, cost int) (int64, error)
		LinkThreads(parent, child int64, iter string, ord int, domain string) error
		LinkThreadsStakeholder(parent, child int64, stakeholder, domain string) error
		GetThreadOrderBefore(parent int64, iter string, order int) (int, error)
		GetPersonteamOrderBefore(personteam, iter string, order int) (int, error)
		SetThreadCostTotal(id int64, cost int) error
		SetThreadOrderParent(thread, parent int64, order int) error
		SetThreadOrderStakeholder(thread int64, stakeholder string, order int) error

		GetThreadrel(id int64, stakeholder string) (*taps.Threadrel, error)
		GetThreadDescendants(id int64, stakeholder string) (map[int64](*taps.Threadrel), error)
		GetThreadAncestors(id int64, stakeholder string) (map[int64](*taps.Threadrel), error)
		GetChildThreadsSkIter(threads []int64, stakeholder, iteration string) (map[int64](*taps.Threadrel), error)
		GetParentThreadsSkIter(threads []int64, stakeholder, iteration string) (map[int64](*taps.Threadrel), error)
		GetThreadrelsByStakeholderIter(stakeholder, iter string) ([](*taps.Threadrel), error)
		GetThreadrelsByParentIter(parent int64, iter string) ([](*taps.Threadrel), error)

		NewStakeholder(thread int64, stakeholder, domain, iter string, ord int, topLvl bool, cost int) error
		GetStakeholderAncestors(thread int64) (map[string]*taps.Personteam, error)
		GetStakeholderDescendants(thread int64) (map[string]*taps.Personteam, error)
		GetStakeholderOrderBefore(stakeholder, iter string, order int) (int, error)
		SetStakeholderCostTotal(thread int64, stakeholder string, cost int) error
		SetStakeholderTopThread(thread int64, stakeholder string, top bool) error
	*/
}

// ErrNotFound indicates that no matching record was found when querying
var ErrNotFound = errors.New("Not found")

type mysqlDB struct {
	conn *sql.DB
}
