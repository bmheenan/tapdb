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
	db.NewThreadStkLink(ths["A"], stks[0], "example.com", "2020 Oct", 50, 2)
	db.NewThreadStkLink(ths["AB"], stks[0], "example.com", "2020 Oct", 10, 1)
	db.NewThreadStkLink(ths["B"], stks[0], "example.com", "2020 Oct", 1, 1)
	res := db.GetThreadsByStkIter(stks[0], "2020 Oct")
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
	res := db.GetThreadsByParentIter(ths["A"], "2020 Oct")
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
	db.NewThreadStkLink(ths["AA"], stks[0], "example.com", "2020 Oct", 1, 1)
	db.NewThreadStkLink(ths["AB"], stks[0], "example.com", "2020 Oct", 1, 1)
	db.NewThreadStkLink(ths["B"], stks[0], "example.com", "2020 Oct", 1, 1)
	ths["AAA"] = db.NewThread("AAA", "example.com", stks[0], "2020 Oct", string(taps.NotStarted), 1, 1)
	db.NewThreadHierLink(ths["AA"], ths["AAA"], "2020 Oct", 1, "example.com")
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
	db.NewThreadStkLink(ths["A"], stks[0], "example.com", "2020 Oct", 1, 1)
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
	res := db.GetThreadrowsByStkIter(stks[0], "2020 Oct")
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

func TestGetThreadrowsByStkIterComplex(t *testing.T) {
	db, stks := setupWithStks()
	iter := "2020-12 Dec"
	dom := "example.com"
	a := db.NewThread("A", dom, stks[0], iter, string(taps.NotStarted), 1, 1)
	aa := db.NewThread("AA", dom, stks[0], iter, string(taps.NotStarted), 1, 1)
	db.NewThreadHierLink(a, aa, iter, 0, dom)
	aaa := db.NewThread("AAA", dom, stks[0], iter, string(taps.NotStarted), 1, 1)
	db.NewThreadHierLink(aa, aaa, iter, 0, dom)
	ab := db.NewThread("AB", dom, stks[0], iter, string(taps.NotStarted), 1, 1)
	db.NewThreadHierLink(a, ab, iter, 0, dom)
	b := db.NewThread("B", dom, stks[0], iter, string(taps.NotStarted), 1, 1)

	db.NewThreadStkLink(a, stks[1], dom, iter, 0, 1)
	db.NewThreadStkLink(aaa, stks[1], dom, iter, 0, 1)
	db.NewThreadStkLink(ab, stks[1], dom, iter, 0, 1)
	db.NewThreadStkLink(b, stks[1], dom, iter, 0, 1)

	res := db.GetThreadrowsByStkIter(stks[1], iter)
	if x, g := 2, len(res); x != g {
		t.Fatalf("Expected result length %v; got %v", x, g)
	}
	if x, g := 2, len(res[0].Children); x != g {
		t.Fatalf("Expected length of A's children %v; got %v", x, g)
	}
}

func TestGetThreadrowsByParentIter(t *testing.T) {
	db, _, ths, errSet := setupWithThreadsStks()
	if errSet != nil {
		t.Errorf("Could not set up test: %v", errSet)
		return
	}
	res := db.GetThreadrowsByParentIter(ths["A"], "2020 Oct")
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

func TestGetParentsForAnc(t *testing.T) {
	db, _, ths, err := setupWithThreadsStks()
	if err != nil {
		t.Errorf("Could not set up test: %v", err)
		return
	}
	pas, err := db.GetThreadParentsForAnc(ths["AA"], ths["A"])
	if err != nil {
		t.Errorf("Could not get parents: %v", err)
		return
	}
	if x, m := 1, len(pas); x != m {
		t.Errorf("Expected length %v; got %v", x, m)
	}
	if x, m := "A", pas[0].Name; x != m {
		t.Errorf("Expected name %v; got %v", x, m)
	}
}

func TestGetParentsOfChildThread(t *testing.T) {
	db, _, ths, err := setupWithThreadsStks()
	if err != nil {
		t.Errorf("Could not set up test: %v", err)
		return
	}
	res := db.GetThreadrowsByChild(ths["AA"])
	if x, g := 1, len(res); x != g {
		t.Fatalf("Expected len %d; got %d", x, g)
	}
}

func TestThreadrowByStkOrdering(t *testing.T) {
	db, stks := setupWithStks()
	ths := map[string]int64{}

	ths["A"] = db.NewThread("A", "example.com", stks[0], "2020-12 Dec", "not started", 1, 1)
	db.NewThreadStkLink(ths["A"], stks[0], "example.com", "2020-12 Dec", 1, 1)

	ths["AA"] = db.NewThread("AA", "example.com", stks[0], "2020-12 Dec", "not started", 1, 1)
	db.NewThreadHierLink(ths["A"], ths["AA"], "2020-12 Dec", 1, "example.com")

	ths["AAA"] = db.NewThread("AAA", "example.com", stks[0], "2020-12 Dec", "not started", 1, 1)
	db.NewThreadHierLink(ths["AA"], ths["AAA"], "2020-12 Dec", 1, "example.com")
	db.NewThreadStkLink(ths["AAA"], stks[0], "example.com", "2020-12 Dec", 3, 1)

	ths["AB"] = db.NewThread("AB", "example.com", stks[0], "2020-12 Dec", "not started", 1, 1)
	db.NewThreadHierLink(ths["A"], ths["AB"], "2020-12 Dec", 2, "example.com")
	db.NewThreadStkLink(ths["AB"], stks[0], "example.com", "2020-12 Dec", 2, 1)

	res := db.GetThreadrowsByStkIter(stks[0], "2020-12 Dec")
	if x, g := 1, len(res); x != g {
		t.Fatalf("Expected length %v; got %v", x, g)
	}
	if x, g := 2, len(res[0].Children); x != g {
		t.Fatalf("Expected length %v; got %v", x, g)
	}
	if x, g := "AB", res[0].Children[0].Name; x != g {
		t.Fatalf("Expected first child %v; got %v", x, g)
	}
	if x, g := "AAA", res[0].Children[1].Name; x != g {
		t.Fatalf("Expected second child %v; got %v", x, g)
	}
}

func TestGetThreadrosByStkWithMultiOwnersInTree(t *testing.T) {
	db, _ := setupWithStks()
	type setupItem struct {
		name       string
		owner      string
		stks       []string
		parent     string
		ord        int
		cost       int
		percentile float64
	}
	domain := "example.com"
	iter := "2020-10 Oct"
	state := "notstarted"
	setups := []setupItem{
		setupItem{
			name:       "A",
			owner:      "a@" + domain,
			stks:       []string{"a@" + domain},
			ord:        10,
			cost:       1,
			percentile: 0,
		},
		setupItem{
			name:       "B",
			parent:     "A",
			owner:      "b@" + domain,
			stks:       []string{"a@" + domain, "b@" + domain},
			ord:        8,
			cost:       1,
			percentile: 0,
		},
		setupItem{
			name:       "C",
			parent:     "B",
			owner:      "a@" + domain,
			stks:       []string{"a@" + domain, "b@" + domain},
			ord:        6,
			cost:       1,
			percentile: 0,
		},
	}
	ths := map[string]int64{}
	for _, s := range setups {
		ths[s.name] = db.NewThread(s.name, domain, s.owner, iter, state, s.percentile, s.cost)
		if s.parent != "" {
			db.NewThreadHierLink(ths[s.parent], ths[s.name], iter, s.ord, domain)
		}
		for _, stk := range s.stks {
			db.NewThreadStkLink(ths[s.name], stk, domain, iter, s.ord, s.cost)
		}
	}
	res := db.GetThreadrowsByStkIter("a@"+domain, iter)
	if x, g := 1, len(res); x != g {
		t.Fatalf("Expected overall length %v; got %v", x, g)
	}
	if x, g := 1, len(res[0].Children); x != g {
		t.Fatalf("Expected length of A.children %v; got %v", x, g)
	}
	if x, g := 1, len(res[0].Children[0].Children); x != g {
		t.Fatalf("Expected length of B.children %v; got %v", x, g)
	}
	if x, g := 0, len(res[0].Children[0].Children[0].Children); x != g {
		t.Fatalf("Expected length of C.children %v; got %v", x, g)
	}
}

func setupWithThreadsStks() (DBInterface, []string, map[string](int64), error) {
	db, stks := setupWithStks()
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
		id := db.NewThread(t.name, "example.com", t.owner, "2020 Oct", string(taps.NotStarted), 1, 1)
		ths[t.name] = id
		for _, p := range t.parents {
			db.NewThreadHierLink(ths[p], ths[t.name], "2020 Oct", t.ord, "example.com")
		}
		for _, s := range t.stks {
			db.NewThreadStkLink(ths[t.name], s, "example.com", "2020 Oct", t.ord, t.cost)
		}
	}
	return db, stks, ths, nil
}
