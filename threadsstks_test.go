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

func TestSetAndGetOrdBeforeForStk(t *testing.T) {
	db, stks, ths, errSet := setupWithThreadsStks()
	if errSet != nil {
		t.Errorf("Could not set up test: %v", errSet)
		return
	}
	errSO := db.SetOrdForStk(ths["A"], stks[0], 5)
	if errSO != nil {
		t.Errorf("Could not set order for thread A for stakeholder %v: %v", stks[0], errSO)
		return
	}
	th, errTh := db.GetThread(ths["A"])
	if errTh != nil {
		t.Errorf("Could not get thread B")
		return
	}
	ob, errOB := db.GetOrdBeforeForStk(stks[0], "2020 Oct", th.Stks[stks[0]].Ord)
	if errOB != nil {
		t.Errorf("Could not get order before thread B: %v", errOB)
		return
	}
	if ob != 4 {
		t.Errorf("Expected order before to be 4, got %v", ob)
		return
	}
}

func TestSetCostForStk(t *testing.T) {
	db, stks, ths, errSet := setupWithThreadsStks()
	if errSet != nil {
		t.Errorf("Could not set up test: %v", errSet)
		return
	}
	errC := db.SetCostForStk(ths["B"], stks[2], 50)
	if errC != nil {
		t.Errorf("Could not set cost: %v", errC)
		return
	}
	th, errTh := db.GetThread(ths["B"])
	if errTh != nil {
		t.Errorf("Could not get thread: %v", errTh)
		return
	}
	if th.Stks[stks[2]].Cost != 50 {
		t.Errorf("Expected cost to be 50, got %v", th.Stks[stks[1]].Cost)
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

func TestGetChildrenByParentStkLinks(t *testing.T) {
	db, stks, errSet := setupWithStks()
	if errSet != nil {
		t.Errorf("Could not set up test: %v", errSet)
		return
	}
	ths := map[string]int64{}
	for _, n := range []string{"A", "B", "C"} {
		id, errN := db.NewThread(n, "example.com", stks[0], "2020 Oct", string(taps.NotStarted), 1, 1)
		if errN != nil {
			t.Errorf("Could not insert thread: %v", errN)
			return
		}
		ths[n] = id
	}
	errL1 := db.NewThreadHierLinkForStk(ths["A"], ths["B"], stks[0], "example.com")
	if errL1 != nil {
		t.Errorf("Could not link A and B: %v", errL1)
		return
	}
	errL2 := db.NewThreadHierLinkForStk(ths["A"], ths["C"], stks[0], "example.com")
	if errL2 != nil {
		t.Errorf("Could not link A and C: %v", errL2)
		return
	}
	chs, errG := db.GetChildrenByParentStkLinks(ths["A"], stks[0])
	if errG != nil {
		t.Errorf("Could not get children: %v", errG)
		return
	}
	if len(chs) != 2 {
		t.Errorf("Expected length 2; got %v", len(chs))
		return
	}
}

func TestGetParentsByChildStkLinks(t *testing.T) {
	db, stks, errSet := setupWithStks()
	if errSet != nil {
		t.Errorf("Could not set up test: %v", errSet)
		return
	}
	ths := map[string]int64{}
	for _, n := range []string{"A", "B", "C"} {
		id, errN := db.NewThread(n, "example.com", stks[0], "2020 Oct", string(taps.NotStarted), 1, 1)
		if errN != nil {
			t.Errorf("Could not insert thread: %v", errN)
			return
		}
		ths[n] = id
	}
	errL1 := db.NewThreadHierLinkForStk(ths["B"], ths["C"], stks[0], "example.com")
	if errL1 != nil {
		t.Errorf("Could not link B and C: %v", errL1)
		return
	}
	errL2 := db.NewThreadHierLinkForStk(ths["A"], ths["C"], stks[0], "example.com")
	if errL2 != nil {
		t.Errorf("Could not link A and C: %v", errL2)
		return
	}
	pas, errG := db.GetParentsByChildStkLinks(ths["C"], stks[0])
	if errG != nil {
		t.Errorf("Could not get parents: %v", errG)
		return
	}
	if len(pas) != 2 {
		t.Errorf("Expected length 2; got %v", len(pas))
		return
	}
}
