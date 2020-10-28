package tapdb

import (
	"fmt"
	"testing"

	"github.com/bmheenan/tapstruct"
)

func TestSingleThreadInsertAndGetThreadrowByPT(t *testing.T) {
	db, pts, errSetup := setupForNewThread()
	if errSetup != nil {
		t.Errorf("Setup failed: %v", errSetup)
		return
	}
	_, err := db.NewThread(&tapstruct.Threaddetail{
		Domain:     "example.com",
		Name:       "Example thread",
		State:      tapstruct.NotStarted,
		Iteration:  "2020 Q4",
		CostDirect: 10,
		CostTotal:  10,
		Owner:      pts[0],
		Order:      0,
		Percentile: 10,
	}, []int64{}, []int64{})
	if err != nil {
		t.Errorf("NewThread returned error: %v", err)
		return
	}
	results, errGet := db.GetThreadrowsByPersonteamPlan(pts[0].Email, []string{"2020 Q4"})
	if errGet != nil {
		t.Errorf("GetThreadrowsByPersonteamPlan returned error: %v", errGet)
		return
	}
	if len(results) != 1 {
		t.Errorf("Expected 1 threadrow, but got %v", len(results))
		return
	}
	if results[0].Name != "Example thread" ||
		results[0].Owner.Email != pts[0].Email {
		t.Errorf("Threadrow didn't have the expected data")
		return
	}
}

func TestThreadNewAndGetTree(t *testing.T) {
	db, pts, errSetup := setupForNewThread()
	if errSetup != nil {
		t.Errorf("Setup failed: %v", errSetup)
		return
	}
	id1, err1 := db.NewThread(&tapstruct.Threaddetail{
		Domain:     "example.com",
		Name:       "Thread 1",
		State:      tapstruct.NotStarted,
		Iteration:  "2020 Oct",
		CostDirect: 10,
		CostTotal:  10,
		Owner:      pts[0],
		Order:      0,
		Percentile: 10,
	}, []int64{}, []int64{})
	if err1 != nil {
		t.Errorf("NewThread returned error: %v", err1)
	}
	id11, err11 := db.NewThread(&tapstruct.Threaddetail{
		Domain:     "example.com",
		Name:       "Thread 1.1",
		State:      tapstruct.NotStarted,
		Iteration:  "2020 Oct",
		CostDirect: 10,
		CostTotal:  10,
		Owner:      pts[0],
		Order:      0,
		Percentile: 10,
	}, []int64{id1}, []int64{})
	if err11 != nil {
		t.Errorf("NewThread returned error: %v", err11)
	}
	_, err12 := db.NewThread(&tapstruct.Threaddetail{
		Domain:     "example.com",
		Name:       "Thread 1.2",
		State:      tapstruct.NotStarted,
		Iteration:  "2020 Oct",
		CostDirect: 10,
		CostTotal:  10,
		Owner:      pts[0],
		Order:      0,
		Percentile: 10,
	}, []int64{id1}, []int64{})
	if err12 != nil {
		t.Errorf("NewThread returned error: %v", err12)
	}
	_, err111 := db.NewThread(&tapstruct.Threaddetail{
		Domain:     "example.com",
		Name:       "Thread 1.1.1",
		State:      tapstruct.NotStarted,
		Iteration:  "2020 Oct",
		CostDirect: 10,
		CostTotal:  10,
		Owner:      pts[0],
		Order:      0,
		Percentile: 10,
	}, []int64{id11}, []int64{})
	if err111 != nil {
		t.Errorf("NewThread returned error: %v", err111)
	}
	_, err2 := db.NewThread(&tapstruct.Threaddetail{
		Domain:     "example.com",
		Name:       "Thread 2",
		State:      tapstruct.NotStarted,
		Iteration:  "2020 Oct",
		CostDirect: 10,
		CostTotal:  10,
		Owner:      pts[0],
		Order:      0,
		Percentile: 10,
	}, []int64{}, []int64{})
	if err2 != nil {
		t.Errorf("NewThread returned error: %v", err2)
	}
	res, errGet := db.GetThreadrowsByPersonteamPlan(pts[0].Email, []string{"2020 Oct"})
	if errGet != nil {
		t.Errorf("GetThreadrowsByPersonteamPlan returned error: %v", errGet)
		return
	}
	if len(res) != 2 ||
		res[0].Name != "Thread 1" ||
		res[1].Name != "Thread 2" ||
		len(res[0].Children) != 2 ||
		res[0].Children[0].Name != "Thread 1.1" ||
		res[0].Children[1].Name != "Thread 1.2" ||
		len(res[0].Children[0].Children) != 1 ||
		res[0].Children[0].Children[0].Name != "Thread 1.1.1" {
		t.Errorf("Tree returned from GetThreadrowsByPersonteamPlan was wrong")
		return
	}
}

func TestDontGetOtherOwners(t *testing.T) {
	db, pts, errSetup := setupForNewThread()
	if errSetup != nil {
		t.Errorf("Setup failed: %v", errSetup)
		return
	}
	_, err1 := db.NewThread(&tapstruct.Threaddetail{
		Domain:     "example.com",
		Name:       "Thread 1",
		State:      tapstruct.NotStarted,
		Iteration:  "2020 Oct",
		CostDirect: 10,
		CostTotal:  10,
		Owner:      pts[0],
		Order:      0,
		Percentile: 10,
	}, []int64{}, []int64{})
	if err1 != nil {
		t.Errorf("NewThread returned error: %v", err1)
	}
	_, err2 := db.NewThread(&tapstruct.Threaddetail{
		Domain:     "example.com",
		Name:       "Thread 2",
		State:      tapstruct.NotStarted,
		Iteration:  "2020 Oct",
		CostDirect: 10,
		CostTotal:  10,
		Owner:      pts[1],
		Order:      0,
		Percentile: 10,
	}, []int64{}, []int64{})
	if err2 != nil {
		t.Errorf("NewThread returned error: %v", err2)
	}
	results, errGet := db.GetThreadrowsByPersonteamPlan(pts[0].Email, []string{"2020 Oct"})
	if errGet != nil {
		t.Errorf("GetThreadrowsByPersonteamPlan returned error: %v", errGet)
		return
	}
	if len(results) != 1 {
		t.Errorf("Expected 1 result, but got %v", len(results))
		return
	}
	if results[0].Name != "Thread 1" {
		t.Errorf("The wrong thread was returned: %v", results[0].Name)
		return
	}
}

func TestStakeholders(t *testing.T) {
	db, pts, errSetup := setupForNewThread()
	if errSetup != nil {
		t.Errorf("Setup failed: %v", errSetup)
		return
	}
	id1, err1 := db.NewThread(&tapstruct.Threaddetail{
		Domain:     "example.com",
		Name:       "Thread 1",
		State:      tapstruct.NotStarted,
		Iteration:  "2020 Oct",
		CostDirect: 10,
		CostTotal:  10,
		Owner:      pts[0],
		Order:      0,
		Percentile: 10,
	}, []int64{}, []int64{})
	if err1 != nil {
		t.Errorf("NewThread returned error: %v", err1)
	}
	_, err2 := db.NewThread(&tapstruct.Threaddetail{
		Domain:     "example.com",
		Name:       "Thread 2",
		State:      tapstruct.NotStarted,
		Iteration:  "2020 Q4",
		CostDirect: 10,
		CostTotal:  10,
		Owner:      pts[0],
		Order:      0,
		Percentile: 10,
	}, []int64{}, []int64{})
	if err2 != nil {
		t.Errorf("NewThread returned error: %v", err2)
	}
	errSt1 := db.NewStakeholder(id1, &pts[1])
	if errSt1 != nil {
		t.Errorf("Could not add stakeholder: %v", errSt1)
	}
	// Setting a stakeholder multiple times should succeed
	errSt2 := db.NewStakeholder(id1, &pts[1])
	if errSt2 != nil {
		t.Errorf("Could not add stakeholder: %v", errSt2)
	}
}

func setupForNewThread() (DBInterface, []tapstruct.Personteam, error) {
	db, errInit := InitDB()
	if errInit != nil {
		return nil, []tapstruct.Personteam{}, fmt.Errorf("Init returned error: %v", errInit)
	}
	errClear := db.ClearDomain("example.com")
	if errClear != nil {
		return nil, []tapstruct.Personteam{}, fmt.Errorf("Clear domain returned error: %v", errClear)
	}
	pt1 := tapstruct.Personteam{
		Email:      "brandon@example.com",
		Domain:     "example.com",
		Name:       "Brandon",
		Abbrev:     "BR",
		ColorF:     "#ffffff",
		ColorB:     "#1f57cf",
		IterTiming: tapstruct.Monthly,
	}
	pt2 := tapstruct.Personteam{
		Email:      "brenda@example.com",
		Domain:     "example.com",
		Name:       "Brenda",
		Abbrev:     "BR",
		ColorF:     "#ffffff",
		ColorB:     "#1f57cf",
		IterTiming: tapstruct.Quarterly,
	}
	errPT := db.NewPersonteam(&pt1, "")
	if errPT != nil {
		return nil, []tapstruct.Personteam{}, fmt.Errorf("NewPersonteam returned an error: %v", errPT)
	}
	errPT = db.NewPersonteam(&pt2, "")
	if errPT != nil {
		return nil, []tapstruct.Personteam{}, fmt.Errorf("NewPersonteam returned an error: %v", errPT)
	}
	return db, []tapstruct.Personteam{pt1, pt2}, nil
}
