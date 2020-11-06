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
	errSk := db.NewStakeholder(id, pts[0], "example.com", "2020 Oct", 0, true, 10)
	if errSk != nil {
		t.Errorf("Could not insert owner as stakeholder: %v", errSk)
		return
	}
	th, errGet := db.GetThreadrel(id, pts[0])
	if errGet != nil {
		t.Errorf("Could not get thread: %v", errGet)
		return
	}
	if th.ID != id || th.Iteration != "2020 Oct" || th.CostDirect != 10 || !th.StakeholderMatch {
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
		errSk := db.NewStakeholder(ths[i].id, pts[0], "example.com", "2020 Oct", 0, true, ths[i].cost)
		if errSk != nil {
			t.Errorf("Could not insert owner as stakeholder: %v", errSk)
			return
		}
		if i > 0 {
			errLnk := db.LinkThreads(ths[i-1].id, ths[i].id, "2020 Oct", 0, "example.com")
			if errLnk != nil {
				t.Errorf("Could not link threads: %v", errLnk)
				return
			}
		}
	}
	des, errDes := db.GetThreadDescendants(ths[1].id, pts[0])
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
	for _, th := range des {
		if !th.StakeholderMatch {
			t.Errorf("Expected StakholderMatch to be true, but it was not for %v", th.ID)
			return
		}
	}
}

func TestThreadLinkingAndAncestors(t *testing.T) {
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
		errSk := db.NewStakeholder(ths[i].id, pts[0], "example.com", "2020 Oct", 0, true, ths[i].cost)
		if errSk != nil {
			t.Errorf("Could not insert owner as stakeholder: %v", errSk)
			return
		}
		if i > 0 {
			errLnk := db.LinkThreads(ths[i-1].id, ths[i].id, "2020 Oct", 0, "example.com")
			if errLnk != nil {
				t.Errorf("Could not link threads: %v", errLnk)
				return
			}
		}
	}
	anc, errAnc := db.GetThreadAncestors(ths[1].id, pts[0])
	if errAnc != nil {
		t.Errorf("Could not get thread ancestors: %v", errAnc)
		return
	}
	if len(anc) != 2 {
		t.Errorf("Expected 2 results; got %v", len(anc))
	}
	totCost := 0
	for _, d := range anc {
		totCost += d.CostDirect
	}
	if totCost != 15 {
		t.Errorf("Expected total cost 15; got %v", totCost)
		return
	}
	for _, th := range anc {
		if !th.StakeholderMatch {
			t.Errorf("Expected StakholderMatch to be true, but it was not for %v", th.ID)
			return
		}
	}
}

func TestGetChildThreadsSkIter(t *testing.T) {
	db, _, errSetup := setupWithPersonteams()
	if errSetup != nil {
		t.Errorf("%v", errSetup)
		return
	}
	ths := map[string](*struct {
		sks     []string
		parents []string
		cost    int
		id      int64
	}){
		"A": {
			sks:  []string{"a@example.com"},
			cost: 5,
		},
		"B": {
			sks:     []string{"a@example.com", "b@example.com"},
			parents: []string{"A"},
			cost:    10,
		},
		"C": {
			sks:     []string{"b@example.com", "c@example.com"},
			parents: []string{"B"},
			cost:    1,
		},
		"D": {
			sks:     []string{"a@example.com", "c@example.com"},
			parents: []string{"B"},
			cost:    2,
		},
	}
	for n, th := range ths {
		var errNew error
		ths[n].id, errNew = db.NewThread(n, "example.com", ths[n].sks[0], "2020 Oct", "not started", 1, ths[n].cost)
		if errNew != nil {
			t.Errorf("Could not insert thread %v: %v", n, errNew)
			return
		}
		for _, sk := range th.sks {
			errSk := db.NewStakeholder(ths[n].id, sk, "example.com", "2020 Oct", 0, true, ths[n].cost)
			if errSk != nil {
				t.Errorf("Could not insert owner as stakeholder: %v", errSk)
				return
			}
		}
		for _, p := range th.parents {
			errLnk := db.LinkThreads(ths[p].id, ths[n].id, "2020 Oct", 0, "example.com")
			if errLnk != nil {
				t.Errorf("Could not link threads: %v", errLnk)
				return
			}
		}
	}
	res0, err0 := db.GetChildThreadsSkIter([]int64{ths["B"].id}, "a@example.com", "2020 Oct")
	if err0 != nil {
		t.Errorf("Could not get child threads by stakeholder and iteration: %v", err0)
		return
	}
	if len(res0) != 1 {
		t.Errorf("Expected res0 to have 1 result but got %v", len(res0))
		return
	}
	res1, err1 := db.GetChildThreadsSkIter([]int64{ths["A"].id}, "c@example.com", "2020 Oct")
	if err1 != nil {
		t.Errorf("Could not get child threads by stakeholder and iteration: %v", err1)
		return
	}
	if len(res1) != 2 {
		t.Errorf("Expected res0 to have 2 results but got %v", len(res1))
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
