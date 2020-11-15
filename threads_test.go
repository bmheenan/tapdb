package tapdb

import (
	"fmt"
	"math"
	"testing"

	"github.com/bmheenan/taps"
)

func TestNewAndGetThread(t *testing.T) {
	db, stks, errSetup := setupWithStks()
	if errSetup != nil {
		t.Errorf("%v", fmt.Errorf("Could not set up test: %v", errSetup))
		return
	}
	id, errIn := db.NewThread("Test thread", "example.com", stks[0], "2020 Oct", "not started", 1, 10)
	if errIn != nil {
		t.Errorf("Could not insert new thread: %v", errIn)
		return
	}
	th, errGet := db.GetThread(id)
	if errGet != nil {
		t.Errorf("Could not get thread: %v", errGet)
		return
	}
	if th.ID != id || th.Iter != "2020 Oct" || th.CostDir != 10 || th.Owner.Email != stks[0] {
		t.Errorf("Retrieved thread didn't have the expected data. Got %v", th)
		return
	}
}

func TestThreadHierDes(t *testing.T) {
	db, stks, errSetup := setupWithStks()
	if errSetup != nil {
		t.Errorf("%v", fmt.Errorf("Could not setup test: %v", errSetup))
		return
	}
	ths := []struct {
		name string
		cost int
		id   int64
	}{
		{
			name: "A",
			cost: 5,
		},
		{
			name: "B",
			cost: 10,
		},
		{
			name: "C",
			cost: 1,
		},
		{
			name: "D",
			cost: 2,
		},
	}
	for i := 0; i < len(ths); i++ {
		var errNew error
		ths[i].id, errNew = db.NewThread(ths[i].name, "example.com", stks[0], "2020 Oct", "not started", 1, ths[i].cost)
		if errNew != nil {
			t.Errorf("Could not insert thread %v: %v", ths[i], errNew)
			return
		}
		if i > 0 {
			errLnk := db.NewThreadHierLink(ths[i-1].id, ths[i].id, "2020 Oct", 0, "example.com")
			if errLnk != nil {
				t.Errorf("Could not link threads: %v", errLnk)
				return
			}
		}
	}
	des, errDes := db.GetThreadDes(ths[1].id)
	if errDes != nil {
		t.Errorf("Could not get thread descendants: %v", errDes)
		return
	}
	if len(des) != 3 {
		t.Errorf("Expected 3 results; got %v", len(des))
	}
	totCost := 0
	for _, d := range des {
		totCost += d.CostDir
	}
	if totCost != 13 {
		t.Errorf("Expected total cost 13; got %v", totCost)
		return
	}
}

func TestThreadHierAns(t *testing.T) {
	db, stks, errSetup := setupWithStks()
	if errSetup != nil {
		t.Errorf("%v", fmt.Errorf("Could not setup test: %v", errSetup))
		return
	}
	ths := []struct {
		name string
		cost int
		id   int64
	}{
		{
			name: "A",
			cost: 5,
		},
		{
			name: "B",
			cost: 10,
		},
		{
			name: "C",
			cost: 1,
		},
		{
			name: "D",
			cost: 2,
		},
	}
	for i := 0; i < len(ths); i++ {
		var errNew error
		ths[i].id, errNew = db.NewThread(ths[i].name, "example.com", stks[0], "2020 Oct", "not started", 1, ths[i].cost)
		if errNew != nil {
			t.Errorf("Could not insert thread %v: %v", ths[i], errNew)
			return
		}
		if i > 0 {
			errLnk := db.NewThreadHierLink(ths[i-1].id, ths[i].id, "2020 Oct", 0, "example.com")
			if errLnk != nil {
				t.Errorf("Could not link threads: %v", errLnk)
				return
			}
		}
	}
	ans, errAns := db.GetThreadAns(ths[1].id)
	if errAns != nil {
		t.Errorf("Could not get thread ancestors: %v", errAns)
		return
	}
	if len(ans) != 2 {
		t.Errorf("Expected 2 results; got %v", len(ans))
	}
	totCost := 0
	for _, d := range ans {
		totCost += d.CostDir
	}
	if totCost != 15 {
		t.Errorf("Expected total cost 15; got %v", totCost)
		return
	}
}

func TestParentOrd(t *testing.T) {
	db, stks, errSetup := setupWithStks()
	if errSetup != nil {
		t.Errorf("%v", fmt.Errorf("Could not setup test: %v", errSetup))
		return
	}
	pid, errP := db.NewThread("P", "example.com", stks[0], "2020 Oct", "not started", 1, 1)
	if errP != nil {
		t.Errorf("Could not insert parent thread: %v", errP)
		return
	}
	ths := map[string](*struct {
		ord int
		id  int64
	}){
		"a": {
			ord: 5,
		},
		"b": {
			ord: 10,
		},
		"c": {
			ord: 100,
		},
	}
	for n, th := range ths {
		var errNew error
		ths[n].id, errNew = db.NewThread(n, "example.com", stks[0], "2020 Oct", "not started", 1, 1)
		if errNew != nil {
			t.Errorf("Could not insert thread %v: %v", n, errNew)
			return
		}
		errLnk := db.NewThreadHierLink(pid, th.id, "2020 Oct", th.ord, "example.com")
		if errLnk != nil {
			t.Errorf("Could not make %v a child of P: %v", n, errLnk)
			return
		}
	}
	ordBef, errOB := db.GetOrdBeforeForParent(pid, "2020 Oct", 10)
	if errOB != nil {
		t.Errorf("Could not get order before 10: %v", errOB)
		return
	}
	if ordBef != 5 {
		t.Errorf("Expected order to be 5, got %v", ordBef)
		return
	}
	errStOrd := db.SetOrdForParent(ths["a"].id, pid, 101)
	if errStOrd != nil {
		t.Errorf("Could not set order of thread: %v", errStOrd)
		return
	}
	ordBef, errOB = db.GetOrdBeforeForParent(pid, "2020 Oct", math.MaxInt32)
	if errOB != nil {
		t.Errorf("Could not get order before 10: %v", errOB)
		return
	}
	if ordBef != 101 {
		t.Errorf("Expected order to be 101, got %v", ordBef)
		return
	}
}

func TestSetCostTot(t *testing.T) {
	db, stks, errSetup := setupWithStks()
	if errSetup != nil {
		t.Errorf("%v", fmt.Errorf("Could not setup test: %v", errSetup))
		return
	}
	thID, errN := db.NewThread("Thread", "example.com", stks[0], "2020 Oct", string(taps.NotStarted), 1, 1)
	if errN != nil {
		t.Errorf("Could not create new thread: %v", errN)
		return
	}
	errC := db.SetCostTot(thID, 10)
	if errC != nil {
		t.Errorf("Could not set total cost of thread: %v", errC)
		return
	}
	th, errTh := db.GetThread(thID)
	if errTh != nil {
		t.Errorf("Could not get thread: %v", errTh)
		return
	}
	if th.CostTot != 10 {
		t.Errorf("Expected CostTot to be 10, but it was %v", th.CostTot)
		return
	}
}

func setupWithStks() (DBInterface, []string, error) {
	db, errSetup := setupEmptyDB()
	if errSetup != nil {
		return nil, nil, errSetup
	}
	es := []string{
		"a@example.com",
		"b@example.com",
		"c@example.com",
	}
	for _, e := range es {
		errStk := db.NewStk(e, "example.com", "Stakeholder "+e, "STK", "#ffffff", "#000000", "monthly")
		if errStk != nil {
			return nil, nil, fmt.Errorf("Error trying to insert new stakeholder: %v", errStk)
		}
	}
	return db, es, nil
}
