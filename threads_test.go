package tapdb

import (
	"fmt"
	"testing"
)

func TestNewAndGetThread(t *testing.T) {
	db, pts, errSetup := setupWithPersonteams()
	if errSetup != nil {
		t.Errorf("%v", errSetup)
		return
	}
	id, errIn := db.NewThread("Test thread", "example.com", pts[0], "2020 Oct", "not started", 1, 10)
	if errIn != nil {
		t.Errorf("Could not insert new thread: %v", errIn)
		return
	}
	th, errGet := db.GetThreadrel(id)
	if errGet != nil {
		t.Errorf("Could not get thread: %v", errGet)
		return
	}
	if th.ID != id || th.Iteration != "2020 Oct" || th.CostDirect != 10 {
		t.Errorf("Retrieved thread didn't have the expected data. Got %v", th)
		return
	}
}

func TestThreadLinkingAndDescendants(t *testing.T) {
	db, pts, errSetup := setupWithPersonteams()
	if errSetup != nil {
		t.Errorf("%v", errSetup)
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
		ths[i].id, errNew = db.NewThread(ths[i].name, "example.com", pts[0], "2020 Oct", "not started", 1, ths[i].cost)
		if errNew != nil {
			t.Errorf("Could not insert thread %v: %v", ths[i], errNew)
			return
		}
		if i > 0 {
			errLnk := db.LinkThreads(ths[i-1].id, ths[i].id, 0, "example.com")
			if errLnk != nil {
				t.Errorf("Could not link threads: %v", errLnk)
				return
			}
		}
	}
	des, errDes := db.GetThreadDescendants(ths[1].id)
	if errDes != nil {
		t.Errorf("Could not get thread descendants: %v", errDes)
		return
	}
	if len(des) != 3 {
		t.Errorf("Expected 3 results; got %v", len(des))
	}
	totCost := 0
	for _, d := range des {
		totCost += d.CostDirect
	}
	if totCost != 13 {
		t.Errorf("Expected total cost 13; got %v", totCost)
		return
	}
}

func setupWithPersonteams() (DBInterface, []string, error) {
	db, errSetup := setupForTest()
	if errSetup != nil {
		return nil, nil, fmt.Errorf("Could not set up test: %v", errSetup)
	}
	es := []string{
		"a@example.com",
		"b@example.com",
		"c@example.com",
	}
	for _, e := range es {
		errPT := db.NewPersonteam(e, "example.com", "PT", "PT", "#ffffff", "#000000", "monthly")
		if errPT != nil {
			return nil, nil, fmt.Errorf("Error trying to insert new personteam: %v", errPT)
		}
	}
	return db, es, nil
}
