package tapdb

import (
	"fmt"
	"math"
	"testing"

	"github.com/bmheenan/taps"
)

func TestNewAndGetThread(t *testing.T) {
	db, stks := setupWithStks()
	id := db.NewThread("Test thread", "example.com", stks[0], "2020 Oct", "not started", 1, 10)
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
	db, stks := setupWithStks()
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
		ths[i].id = db.NewThread(ths[i].name, "example.com", stks[0], "2020 Oct", "not started", 1, ths[i].cost)
		if i > 0 {
			db.NewThreadHierLink(ths[i-1].id, ths[i].id, "2020 Oct", 0, "example.com")
		}
	}
	des := db.GetThreadDes(ths[1].id)
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

func TestThreadHierLinkMultipleTimes(t *testing.T) {
	db, stks := setupWithStks()
	ida := db.NewThread("A", "example.com", stks[0], "2020-12 Dec", string(taps.NotStarted), 1, 1)
	idb := db.NewThread("B", "example.com", stks[0], "2020-12 Dec", string(taps.NotStarted), 1, 1)
	for i := 0; i < 2; i++ {
		db.NewThreadHierLink(ida, idb, "2020-12 Dec", 0, "example.com")
	}
}

func TestThreadHierAns(t *testing.T) {
	db, stks := setupWithStks()
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
		ths[i].id = db.NewThread(ths[i].name, "example.com", stks[0], "2020 Oct", "not started", 1, ths[i].cost)
		if i > 0 {
			db.NewThreadHierLink(ths[i-1].id, ths[i].id, "2020 Oct", 0, "example.com")
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
	db, stks := setupWithStks()
	pid := db.NewThread("P", "example.com", stks[0], "2020 Oct", "not started", 1, 1)
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
		ths[n].id = db.NewThread(n, "example.com", stks[0], "2020 Oct", "not started", 1, 1)
		db.NewThreadHierLink(pid, th.id, "2020 Oct", th.ord, "example.com")
	}
	ordBef := db.GetOrdBeforeForParent(pid, "2020 Oct", 10)
	if ordBef != 5 {
		t.Errorf("Expected order to be 5, got %v", ordBef)
		return
	}
	db.SetOrdForParent(ths["a"].id, pid, 101)
	ordBef = db.GetOrdBeforeForParent(pid, "2020 Oct", math.MaxInt32)
	if ordBef != 101 {
		t.Errorf("Expected order to be 101, got %v", ordBef)
		return
	}
}

func TestSetCostTot(t *testing.T) {
	db, stks := setupWithStks()
	thID := db.NewThread("Thread", "example.com", stks[0], "2020 Oct", string(taps.NotStarted), 1, 1)
	db.SetCostTot(thID, 10)
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

func setupWithStks() (DBInterface, []string) {
	db, err := setupEmptyDB()
	if err != nil {
		panic(err)
	}
	es := []string{
		"a@example.com",
		"b@example.com",
		"c@example.com",
	}
	for _, e := range es {
		err := db.NewStk(e, "example.com", "Stakeholder "+e, "STK", "#ffffff", "#000000", "monthly")
		if err != nil {
			panic(fmt.Sprintf("Error trying to insert new stakeholder: %v", err))
		}
	}
	return db, es
}
