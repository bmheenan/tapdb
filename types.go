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

	// GetStkAns gets all stakeholders that are ancestors of the stakeholder with `email` (including itself)
	GetStkAns(email string) (map[string](*taps.Stakeholder), error)

	// threads.go

	// NewThread makes a new thread with the given info. It returns the thread's new id
	NewThread(name, domain, owner, iter, state string, percentile float64, cost int) (int64, error)

	// NewThreadHierLink makes `parent` a parent of `child`. `child` will show up in `iter` (which should be the
	// child's base iteration expresed in the cadence of `parent`'s owner)
	NewThreadHierLink(parent, child int64, iter string, ord int, domain string) error

	// GetOrdBeforeForParent returns the highest order of any thread under `parent` in `iter`, thats lower than
	// `order`
	GetOrdBeforeForParent(parent int64, iter string, ord int) (int, error)

	// SetThreadOrderForParent sets `thread`'s order under `parent` to `order`
	SetOrdForParent(thread, parent int64, ord int) error

	// SetThreadCosttot sets `thread`'s total cost (including descendants) to `cost`
	SetCostTot(thread int64, cost int) error

	// threadsget.go

	// GetThread returns the Thread with id matching `thread`
	GetThread(thread int64) (*taps.Thread, error)

	// GetThreadsByStkIter returns all threads that have `stk` as a stakeholder in `iter`, ordered by ord
	GetThreadsByStkIter(stk, iter string) ([](*taps.Thread), error)

	// GetThreadsByParentIter returns all threads that are children of `parent` in `iter` ordered by ord
	GetThreadsByParentIter(parent int64, iter string) ([](*taps.Thread), error)

	// GetThreadDes gets all descendant threads of `thread` (including itself)
	GetThreadDes(thread int64) (map[int64](*taps.Thread), error)

	// GetThreadAns gets all ancestor threads of `thread` (including itself)
	GetThreadAns(thread int64) (map[int64](*taps.Thread), error)

	// GetThreadChildrenByStkIter returns the smallest map of threads that contains all descendants of `threads` which
	// have `stk` as their stakeholder in `iter`.
	GetThreadChildrenByStkIter(threads []int64, stk, iter string) (map[int64](*taps.Thread), error)

	// GetThreadParentsByStkIter returns the smallest map of threads that contains all ancestors of `threads` which
	// have `stk` as their stakeholder in `iter`.
	GetThreadParentsByStkIter(threads []int64, stk, iter string) (map[int64](*taps.Thread), error)

	// GetThreadrowsByStkIter returns a hierarchical, ordered array of Threadrows where `stk` is a stakeholder in `iter`
	GetThreadrowsByStkIter(stk, iter string) ([](*taps.Threadrow), error)

	// GetThreadrowsByParentIter returns a hierarchical, ordered array of Threadrows that are descendants of `parent`
	// in `iter`
	GetThreadrowsByParentIter(parent int64, iter string) ([](*taps.Threadrow), error)

	// threadsstks.go

	// NewThreadStkLink makes `stk` a stakeholder of `thread`, with `thread` showing in `iter` in order `ord`, costing
	// `cost` for this stakeholder and all subteams + teammembers. It's `toplvl` if the stakeholder has no ancestor
	// threads in the same iteration
	NewThreadStkLink(thread int64, stk, domain, iter string, ord int, toplvl bool, cost int) error

	// NewThreadHierLinkForStk makes `parent` a parent of `child` in `stk`'s context. `child` should be a descendant of
	// `parent` and `stk` should be a stakeholder of both of them
	NewThreadHierLinkForStk(parent, child int64, stk, domain string) error

	// GetOrdBeforeForStk returns the highest order of any thread with `stk` as a stakeholder in `iter`, thats lower
	// than `order`
	GetOrdBeforeForStk(stk, iter string, ord int) (int, error)

	// SetThreadOrderForStk sets `thread`'s order under `stk` to `order`
	SetOrdForStk(thread int64, stk string, ord int) error

	// SetCostForStk sets the cost of `thread` (just the peices owned by `stk`) to `cost`
	SetCostForStk(thread int64, stk string, cost int) error

	// SetTopForStk sets if `thread` is a top-level thread for `stk`-if it has no ancestors in the same iteration where
	// `stk` is a stakeholder
	SetTopForStk(thread int64, stk string, top bool) error
}

// ErrNotFound indicates that no matching record was found when querying
var ErrNotFound = errors.New("Not found")

type mysqlDB struct {
	conn *sql.DB
}
