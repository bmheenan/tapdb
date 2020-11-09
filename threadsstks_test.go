package tapdb

import (
	"fmt"
	"testing"

	"github.com/bmheenan/taps"
)

func TestThStkLinkAndGetHasStks(t *testing.T) {
	db, stks, ths, errSet := setupWithThreads()
	if errSet != nil {
		t.Errorf("Could not set up test: %v", errSet)
	}
	errL1 := db.NewThreadStkLink(ths["A"], stks[0], "example.com", "2020 Oct", 1, true, 1)
	if errL1 != nil {
		t.Errorf("Could not link thread and stakeholder: %v", errL1)
		return
	}
	errL2 := db.NewThreadStkLink(ths["AA"], stks[0], "example.com", "2020 Oct", 1, true, 1)
	if errL2 != nil {
		t.Errorf("Could not link thread and stakeholder: %v", errL2)
		return
	}
	errL3 := db.NewThreadStkLink(ths["AB"], stks[0], "example.com", "2020 Oct", 1, true, 1)
	if errL3 != nil {
		t.Errorf("Could not link thread and stakeholder: %v", errL3)
		return
	}
	errL4 := db.NewThreadStkLink(ths["B"], stks[0], "example.com", "2020 Oct", 1, true, 1)
	if errL4 != nil {
		t.Errorf("Could not link thread and stakeholder: %v", errL4)
		return
	}
	res, errR := db.GetThreadDes(ths["A"])
	if errR != nil {
		t.Errorf("Could not get thread descendants: %v", errR)
		return
	}
	if len(res) != 3 {
		t.Errorf("Expected length to be 3, but it was %v", len(res))
		return
	}
	if _, ok := res[ths["A"]]; !ok {
		t.Errorf("A was not in returned results")
		return
	}
	if _, ok := res[ths["A"]].Stks[stks[0]]; !ok {
		t.Errorf("Stks[0] was not a stakeholder of A")
		return
	}
	if _, ok := res[ths["AA"]]; !ok {
		t.Errorf("AA was not in returned results")
		return
	}
	if _, ok := res[ths["AA"]].Stks[stks[0]]; !ok {
		t.Errorf("Stks[0] was not a stakeholder of AA")
		return
	}
	if _, ok := res[ths["AB"]]; !ok {
		t.Errorf("AB was not in returned results")
		return
	}
	if _, ok := res[ths["AB"]].Stks[stks[0]]; !ok {
		t.Errorf("Stks[0] was not a stakeholder of AB")
		return
	}
}

func setupWithThreads() (DBInterface, []string, map[string](int64), error) {
	db, stks, errSet := setupWithStks()
	if errSet != nil {
		return nil, nil, nil, errSet
	}
	ths := map[string](int64){}
	for _, t := range []struct {
		n  string
		ps []string
		o  int
	}{
		{
			n:  "A",
			ps: []string{},
		},
		{
			n:  "AA",
			ps: []string{"A"},
			o:  1,
		},
		{
			n:  "AB",
			ps: []string{"A"},
			o:  2,
		},
		{
			n:  "B",
			ps: []string{},
		},
	} {
		id, errN := db.NewThread(t.n, "example.com", stks[0], "2020 Oct", string(taps.NotStarted), 1, 1)
		if errN != nil {
			return nil, nil, nil, fmt.Errorf("Could not insert thread: %v", errN)
		}
		ths[t.n] = id
		for _, p := range t.ps {
			errL := db.NewThreadHierLink(ths[p], ths[t.n], "2020 Oct", t.o, "example.com")
			if errL != nil {
				return nil, nil, nil, fmt.Errorf("Could not link parent %v with child %v: %v", p, t.n, errL)
			}
		}
	}
	return db, stks, ths, nil
}
