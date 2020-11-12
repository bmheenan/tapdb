package tapdb

import (
	"fmt"
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

func TestGetThreadrowsByStkIter(t *testing.T) {
	db, stks, _, errSet := setupWithThreadsStks()
	if errSet != nil {
		t.Errorf("Could not set up test: %v", errSet)
		return
	}
	res, errG := db.GetThreadrowsByStkIter(stks[0], "2020 Oct")
	if errG != nil {
		t.Errorf("Could not get threadrows for stakeholder %v and iteration 2020 Oct: %v", stks[0], errG)
		return
	}
	if len(res) != 2 {
		t.Errorf("Expected 2 threads but got %v", len(res))
		return
	}
	if len(res[0].Children) != 1 {
		t.Errorf("Expected 1 child thread of A but got %v", len(res[0].Children))
		return
	}
	if len(res[1].Children) != 0 {
		t.Errorf("Expected 0 child threads of B but got %v", len(res[1].Children))
		return
	}
	if res[0].Name != "A" {
		t.Errorf("A was not the first thread in the results")
		return
	}
	if res[0].Children[0].Name != "AA" {
		t.Errorf("AA was not the second thread in the results")
		return
	}
	if res[1].Name != "B" {
		t.Errorf("B was not the second thread in the results")
		return
	}
}

func TestGetThreadrowsByParentIter(t *testing.T) {
	db, _, ths, errSet := setupWithThreadsStks()
	if errSet != nil {
		t.Errorf("Could not set up test: %v", errSet)
		return
	}
	res, errG := db.GetThreadrowsByParentIter(ths["A"], "2020 Oct")
	if errG != nil {
		t.Errorf("Could not get threadrows for thread A and iteration 2020 Oct: %v", errG)
		return
	}
	if len(res) != 2 {
		t.Errorf("Expected 2 threads but got %v", len(res))
		return
	}
	if len(res[0].Children) != 0 {
		t.Errorf("Expected 0 child threads of AA but got %v", len(res[0].Children))
		return
	}
	if len(res[1].Children) != 0 {
		t.Errorf("Expected 0 child threads of AB but got %v", len(res[1].Children))
		return
	}
	if res[0].Name != "AA" {
		t.Errorf("AA was not the first thread in the results")
		return
	}
	if res[1].Name != "AB" {
		t.Errorf("AB was not the second thread in the results")
		return
	}
}

func setupWithThreadsStks() (DBInterface, []string, map[string](int64), error) {
	db, stks, errSet := setupWithStks()
	if errSet != nil {
		return nil, nil, nil, errSet
	}
	errLStk := db.NewStkHierLink(stks[0], stks[1], "example.com")
	if errLStk != nil {
		return nil, nil, nil, fmt.Errorf("Could not link stakeholders in hierarchy: %v", errLStk)
	}
	ths := map[string](int64){}
	for _, t := range []struct {
		name         string
		parents      []string
		parentsByStk map[string]([]string)
		cost         int
		ord          int
		owner        string
		stks         []string
		topFor       map[string](bool)
	}{
		{
			name:  "A",
			owner: stks[0],
			stks:  []string{stks[0]},
			ord:   3,
			topFor: map[string](bool){
				stks[0]: true,
			},
		},
		{
			name:    "AA",
			parents: []string{"A"},
			owner:   stks[1],
			stks:    []string{stks[0], stks[1]},
			ord:     1,
			topFor: map[string](bool){
				stks[1]: true,
			},
			parentsByStk: map[string]([]string){
				stks[0]: []string{"A"},
			},
		},
		{
			name:    "AB",
			parents: []string{"A"},
			owner:   stks[1],
			stks:    []string{stks[1]},
			topFor: map[string](bool){
				stks[1]: true,
			},
			ord: 2,
		},
		{
			name:  "B",
			ord:   4,
			owner: stks[2],
			stks:  []string{stks[0], stks[2]},
			topFor: map[string](bool){
				stks[0]: true,
				stks[2]: true,
			},
		},
	} {
		id, errN := db.NewThread(t.name, "example.com", t.owner, "2020 Oct", string(taps.NotStarted), 1, 1)
		if errN != nil {
			return nil, nil, nil, fmt.Errorf("Could not insert thread: %v", errN)
		}
		ths[t.name] = id
		for _, p := range t.parents {
			errL := db.NewThreadHierLink(ths[p], ths[t.name], "2020 Oct", t.ord, "example.com")
			if errL != nil {
				return nil, nil, nil, fmt.Errorf("Could not link parent %v with child %v: %v", p, t.name, errL)
			}
		}
		for _, s := range t.stks {
			errS := db.NewThreadStkLink(ths[t.name], s, "example.com", "2020 Oct", t.ord, t.topFor[s], t.cost)
			if errS != nil {
				return nil, nil, nil, fmt.Errorf("Could not add stakeholder %v to thread %v: %v", s, t.name, errS)
			}
		}
		for s, ps := range t.parentsByStk {
			for _, p := range ps {
				errLSH := db.NewThreadHierLinkForStk(ths[p], ths[t.name], s, "example.com")
				if errLSH != nil {
					return nil, nil, nil, fmt.Errorf(
						"Could not link threads %v and %v for stakeholder %v: %v",
						p,
						t.name,
						s,
						errLSH,
					)
				}
			}
		}
	}
	return db, stks, ths, nil
}
