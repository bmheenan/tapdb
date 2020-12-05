package tapdb

import (
	"testing"

	"github.com/bmheenan/taps"
)

func TestThStkLinkAndGetHasStks(t *testing.T) {
	db, stks, ths, errSet := setupWithThreads()
	if errSet != nil {
		t.Errorf("Could not set up test: %v", errSet)
	}
	db.NewThreadStkLink(ths["A"], stks[0], "example.com", "2020 Oct", 1, 1)
	db.NewThreadStkLink(ths["AA"], stks[0], "example.com", "2020 Oct", 1, 1)
	db.NewThreadStkLink(ths["AB"], stks[0], "example.com", "2020 Oct", 1, 1)
	db.NewThreadStkLink(ths["B"], stks[0], "example.com", "2020 Oct", 1, 1)
	res := db.GetThreadDes(ths["A"])
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
	ob := db.GetOrdBeforeForStk(stks[0], "2020 Oct", th.Stks[stks[0]].Ord)
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
	db, stks := setupWithStks()
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
		id := db.NewThread(t.n, "example.com", stks[0], "2020 Oct", string(taps.NotStarted), 1, 1)
		ths[t.n] = id
		for _, p := range t.ps {
			db.NewThreadHierLink(ths[p], ths[t.n], "2020 Oct", t.o, "example.com")
		}
	}
	return db, stks, ths, nil
}

/*
func TestGetChildrenByParentStkLinks(t *testing.T) {
	db, stks := setupWithStks()
	ths := map[string]int64{}
	for _, n := range []string{"A", "B", "C"} {
		id := db.NewThread(n, "example.com", stks[0], "2020 Oct", string(taps.NotStarted), 1, 1)
		ths[n] = id
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
	db, stks := setupWithStks()
	ths := map[string]int64{}
	for _, n := range []string{"A", "B", "C"} {
		id := db.NewThread(n, "example.com", stks[0], "2020 Oct", string(taps.NotStarted), 1, 1)
		ths[n] = id
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
*/
