package tapdb

import (
	"testing"

	taps "github.com/bmheenan/taps"
)

func TestGetThreadsByStkIter(t *testing.T) {
	db, stks, ths, errSet := setupWithThreads()
	if errSet != nil {
		t.Errorf("Could not set up test: %v", errSet)
		return
	}
	errL1 := db.NewThreadStkLink(ths["A"], stks[0], "example.com", "2020 Oct", 50, true, 2)
	if errL1 != nil {
		t.Errorf("Could not link thread and stakeholder: %v", errL1)
		return
	}
	errL2 := db.NewThreadStkLink(ths["AB"], stks[0], "example.com", "2020 Oct", 10, true, 1)
	if errL2 != nil {
		t.Errorf("Could not link thread and stakeholder: %v", errL2)
		return
	}
	errL3 := db.NewThreadStkLink(ths["B"], stks[0], "example.com", "2020 Oct", 1, true, 1)
	if errL3 != nil {
		t.Errorf("Could not link thread and stakeholder: %v", errL3)
		return
	}
	res, errR := db.GetThreadsByStkIter(stks[0], "2020 Oct")
	if errR != nil {
		t.Errorf("Could not get threads by stakeholder and iteration: %v", errR)
		return
	}
	if len(res) != 3 {
		t.Errorf("Expected 3 results, but got %v", len(res))
		return
	}
	if res[0].Name != "B" {
		t.Errorf("First thread was not B")
	}
	if res[1].Name != "AB" {
		t.Errorf("First thread was not AB")
	}
	if res[2].Name != "A" {
		t.Errorf("First thread was not A")
	}
}

func TestGetThreadsByparentIter(t *testing.T) {
	db, _, ths, errSet := setupWithThreads()
	if errSet != nil {
		t.Errorf("Could not set up test: %v", errSet)
		return
	}
	res, errR := db.GetThreadsByParentIter(ths["A"], "2020 Oct")
	if errR != nil {
		t.Errorf("Could not get threads by stakeholder and iteration: %v", errR)
		return
	}
	if len(res) != 2 {
		t.Errorf("Expected 3 results, but got %v", len(res))
		return
	}
	if res[0].Name != "AA" {
		t.Errorf("First thread was not AA")
	}
	if res[1].Name != "AB" {
		t.Errorf("First thread was not AB")
	}
}

func TestGetThreadChildrenByStkIter(t *testing.T) {
	db, stks, ths, errSet := setupWithThreads()
	if errSet != nil {
		t.Errorf("Could not set up test: %v", errSet)
		return
	}
	errL1 := db.NewThreadStkLink(ths["AA"], stks[0], "example.com", "2020 Oct", 1, true, 1)
	if errL1 != nil {
		t.Errorf("Could not link thread and stakeholder: %v", errL1)
		return
	}
	errL2 := db.NewThreadStkLink(ths["AB"], stks[0], "example.com", "2020 Oct", 1, true, 1)
	if errL2 != nil {
		t.Errorf("Could not link thread and stakeholder: %v", errL2)
		return
	}
	errL3 := db.NewThreadStkLink(ths["B"], stks[0], "example.com", "2020 Oct", 1, true, 1)
	if errL3 != nil {
		t.Errorf("Could not link thread and stakeholder: %v", errL3)
		return
	}
	var errN error
	ths["AAA"], errN = db.NewThread("AAA", "example.com", stks[0], "2020 Oct", string(taps.NotStarted), 1, 1)
	if errN != nil {
		t.Errorf("Could not insert new thread: %v", errN)
		return
	}
	errL4 := db.NewThreadHierLink(ths["AA"], ths["AAA"], "2020 Oct", 1, "example.com")
	if errL4 != nil {
		t.Errorf("Could not link AAA to AA: %v", errL4)
		return
	}
	chs, errCh := db.GetThreadChildrenByStkIter([]int64{ths["A"]}, stks[0], "2020 Oct")
	if errCh != nil {
		t.Errorf("Could not get thread children: %v", errCh)
		return
	}
	if len(chs) != 2 {
		t.Errorf("Expected length 2, got %v", len(chs))
		return
	}
	if _, ok := chs[ths["AA"]]; !ok {
		t.Errorf("chs did not have AA")
		return
	}
	if _, ok := chs[ths["AB"]]; !ok {
		t.Errorf("chs did not have AB")
		return
	}
}

func TestGetThreadParentsByStkIter(t *testing.T) {
	db, stks, ths, errSet := setupWithThreads()
	if errSet != nil {
		t.Errorf("Could not set up test: %v", errSet)
		return
	}
	errL1 := db.NewThreadStkLink(ths["A"], stks[0], "example.com", "2020 Oct", 1, true, 1)
	if errL1 != nil {
		t.Errorf("Could not link thread and stakeholder: %v", errL1)
		return
	}
	chs, errCh := db.GetThreadParentsByStkIter([]int64{ths["A"]}, stks[0], "2020 Oct")
	if errCh != nil {
		t.Errorf("Could not get thread children: %v", errCh)
		return
	}
	if len(chs) != 1 {
		t.Errorf("Expected length 2, got %v", len(chs))
		return
	}
	if _, ok := chs[ths["A"]]; !ok {
		t.Errorf("chs did not have AA")
		return
	}
}
