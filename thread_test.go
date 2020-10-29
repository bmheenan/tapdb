package tapdb

import (
	"fmt"
	"testing"

	"github.com/bmheenan/tapstruct"
)

func TestSingleThreadInsertAndGetThreadrowByPT(t *testing.T) {
	db, pts, errSetup := setupForNewThread()
	if errSetup != nil {
		t.Errorf("Setup failed: %v", errSetup)
		return
	}
	_, err := db.NewThread(getThread1(pts[0]), []*tapstruct.Threadrow{}, []*tapstruct.Threadrow{})
	if err != nil {
		t.Errorf("NewThread returned error: %v", err)
		return
	}
	results, errGet := db.GetThreadrowsByPersonteamPlan(pts[0].Email, []string{"2020 Oct"})
	if errGet != nil {
		t.Errorf("GetThreadrowsByPersonteamPlan returned error: %v", errGet)
		return
	}
	if len(results) != 1 {
		t.Errorf("Expected 1 threadrow, but got %v", len(results))
		return
	}
	if results[0].Name != "Thread 1" ||
		results[0].Owner.Email != pts[0].Email {
		t.Errorf("Threadrow didn't have the expected data")
		return
	}
}

func TestThreadNewAndGetTree(t *testing.T) {
	db, pts, errSetup := setupForNewThread()
	if errSetup != nil {
		t.Errorf("Setup failed: %v", errSetup)
		return
	}
	id1, err1 := db.NewThread(getThread1(pts[0]), []*tapstruct.Threadrow{}, []*tapstruct.Threadrow{})
	if err1 != nil {
		t.Errorf("NewThread returned error: %v", err1)
	}
	id11, err11 := db.NewThread(&tapstruct.Threaddetail{
		Domain:     "example.com",
		Name:       "Thread 1.1",
		State:      tapstruct.NotStarted,
		Iteration:  "2020 Oct",
		CostDirect: 10,
		CostTotal:  10,
		Owner:      pts[0],
		Order:      0,
		Percentile: 10,
	}, []*tapstruct.Threadrow{&tapstruct.Threadrow{
		ID: id1,
	}}, []*tapstruct.Threadrow{})
	if err11 != nil {
		t.Errorf("NewThread returned error: %v", err11)
	}
	_, err12 := db.NewThread(&tapstruct.Threaddetail{
		Domain:     "example.com",
		Name:       "Thread 1.2",
		State:      tapstruct.NotStarted,
		Iteration:  "2020 Oct",
		CostDirect: 10,
		CostTotal:  10,
		Owner:      pts[0],
		Order:      0,
		Percentile: 10,
	}, []*tapstruct.Threadrow{&tapstruct.Threadrow{
		ID: id1,
	}}, []*tapstruct.Threadrow{})
	if err12 != nil {
		t.Errorf("NewThread returned error: %v", err12)
	}
	_, err111 := db.NewThread(&tapstruct.Threaddetail{
		Domain:     "example.com",
		Name:       "Thread 1.1.1",
		State:      tapstruct.NotStarted,
		Iteration:  "2020 Oct",
		CostDirect: 10,
		CostTotal:  10,
		Owner:      pts[0],
		Order:      0,
		Percentile: 10,
	}, []*tapstruct.Threadrow{&tapstruct.Threadrow{
		ID: id11,
	}}, []*tapstruct.Threadrow{})
	if err111 != nil {
		t.Errorf("NewThread returned error: %v", err111)
	}
	_, err2 := db.NewThread(getThread2(pts[0]), []*tapstruct.Threadrow{}, []*tapstruct.Threadrow{})
	if err2 != nil {
		t.Errorf("NewThread returned error: %v", err2)
	}
	res, errGet := db.GetThreadrowsByPersonteamPlan(pts[0].Email, []string{"2020 Oct"})
	if errGet != nil {
		t.Errorf("GetThreadrowsByPersonteamPlan returned error: %v", errGet)
		return
	}
	if len(res) != 2 ||
		res[0].Name != "Thread 1" ||
		res[1].Name != "Thread 2" ||
		len(res[0].Children) != 2 ||
		res[0].Children[0].Name != "Thread 1.1" ||
		res[0].Children[1].Name != "Thread 1.2" ||
		len(res[0].Children[0].Children) != 1 ||
		res[0].Children[0].Children[0].Name != "Thread 1.1.1" {

		t.Errorf("Tree returned from GetThreadrowsByPersonteamPlan was wrong")
		return
	}
}

func TestDontGetOtherOwners(t *testing.T) {
	db, pts, errSetup := setupForNewThread()
	if errSetup != nil {
		t.Errorf("Setup failed: %v", errSetup)
		return
	}
	_, err1 := db.NewThread(getThread1(pts[0]), []*tapstruct.Threadrow{}, []*tapstruct.Threadrow{})
	if err1 != nil {
		t.Errorf("NewThread returned error: %v", err1)
	}
	_, err2 := db.NewThread(getThread2(pts[1]), []*tapstruct.Threadrow{}, []*tapstruct.Threadrow{})
	if err2 != nil {
		t.Errorf("NewThread returned error: %v", err2)
	}
	results, errGet := db.GetThreadrowsByPersonteamPlan(pts[0].Email, []string{"2020 Oct"})
	if errGet != nil {
		t.Errorf("GetThreadrowsByPersonteamPlan returned error: %v", errGet)
		return
	}
	if len(results) != 1 {
		t.Errorf("Expected 1 result, but got %v", len(results))
		return
	}
	if results[0].Name != "Thread 1" {
		t.Errorf("The wrong thread was returned: %v", results[0].Name)
		return
	}
}

func TestStakeholders(t *testing.T) {
	db, pts, errSetup := setupForNewThread()
	if errSetup != nil {
		t.Errorf("Setup failed: %v", errSetup)
		return
	}
	id1, err1 := db.NewThread(getThread1(pts[0]), []*tapstruct.Threadrow{}, []*tapstruct.Threadrow{})
	if err1 != nil {
		t.Errorf("NewThread returned error: %v", err1)
	}
	_, err2 := db.NewThread(getThread2(pts[1]), []*tapstruct.Threadrow{}, []*tapstruct.Threadrow{})
	if err2 != nil {
		t.Errorf("NewThread returned error: %v", err2)
	}
	errSt1 := db.NewStakeholder(id1, &pts[1])
	if errSt1 != nil {
		t.Errorf("Could not add stakeholder: %v", errSt1)
	}
	// Setting a stakeholder multiple times should succeed
	errSt2 := db.NewStakeholder(id1, &pts[1])
	if errSt2 != nil {
		t.Errorf("Could not add stakeholder: %v", errSt2)
		return
	}
	res0, err0 := db.GetThreadrowsByPersonteamPlan(pts[0].Email, []string{"2020 Oct"})
	if err0 != nil {
		t.Errorf("Could not get results for %v: %v", pts[1].Email, err0)
		return
	}
	if len(res0) != 1 {
		t.Errorf("For %v, expected 1 result, but got %v", pts[0].Email, len(res0))
		return
	}
	res1, err1 := db.GetThreadrowsByPersonteamPlan(pts[1].Email, []string{"2020 Oct"})
	if err1 != nil {
		t.Errorf("Could not get results for %v: %v", pts[1].Email, err1)
		return
	}
	if len(res1) != 2 {
		t.Errorf("For %v, expected 2 results, but got %v", pts[1].Email, len(res1))
		return
	}
}

func TestThreadOrderingBeforeSameOwner(t *testing.T) {
	db, pts, errSetup := setupForNewThread()
	if errSetup != nil {
		t.Errorf("Setup failed: %v", errSetup)
		return
	}
	id0, errN0 := db.NewThread(getThread1(pts[0]), []*tapstruct.Threadrow{}, []*tapstruct.Threadrow{})
	if errN0 != nil {
		t.Errorf("Could not create new thread: %v", errN0)
		return
	}
	id1, errN1 := db.NewThread(getThread2(pts[0]), []*tapstruct.Threadrow{}, []*tapstruct.Threadrow{})
	if errN1 != nil {
		t.Errorf("Could not create new thread: %v", errN1)
		return
	}
	res0, errG0 := db.GetThreadrowsByPersonteamPlan(pts[0].Email, []string{"2020 Oct"})
	if errG0 != nil {
		t.Errorf("Could not get threads: %v", errG0)
		return
	}
	if len(res0) != 2 {
		t.Errorf("Expected 2 results but got %v", len(res0))
		return
	}
	if res0[0].ID != id0 {
		t.Errorf("The first thread wasn't the first one added. It was %v", res0[0].ID)
		return
	}
	if res0[1].ID != id1 {
		t.Errorf("The second thread wasn't the second one added. It was %v", res0[1].ID)
		return
	}
	if res0[0].Order >= res0[1].Order {
		t.Errorf("The first thread's order is not less than the second thread")
		return
	}
	if res0[0].Percentile >= res0[1].Percentile {
		t.Errorf("The first thread's percentile is not less than the second thread")
		return
	}
	// Switch the threads
	errM := db.MoveThread(&res0[1], Before, &res0[0])
	if errM != nil {
		t.Errorf("Could not move thread: %v", errM)
		return
	}
	res1, errG1 := db.GetThreadrowsByPersonteamPlan(pts[0].Email, []string{"2020 Oct"})
	if errG1 != nil {
		t.Errorf("Could not get threads after move: %v", errG1)
		return
	}
	if len(res1) != 2 {
		t.Errorf("Expected 2 results but got %v", len(res1))
		return
	}
	if res1[0].ID != id1 {
		t.Errorf("The first thread wasn't the second one added. It was %v", res1[0].ID)
		return
	}
	if res1[1].ID != id0 {
		t.Errorf("The second thread wasn't the first one added. It was %v", res1[1].ID)
		return
	}
}

func TestThreadOrderingAfterSameOwner(t *testing.T) {
	db, pts, errSetup := setupForNewThread()
	if errSetup != nil {
		t.Errorf("Setup failed: %v", errSetup)
		return
	}
	id0, errN0 := db.NewThread(getThread1(pts[0]), []*tapstruct.Threadrow{}, []*tapstruct.Threadrow{})
	if errN0 != nil {
		t.Errorf("Could not create new thread: %v", errN0)
		return
	}
	id1, errN1 := db.NewThread(getThread2(pts[0]), []*tapstruct.Threadrow{}, []*tapstruct.Threadrow{})
	if errN1 != nil {
		t.Errorf("Could not create new thread: %v", errN1)
		return
	}
	id2, errN2 := db.NewThread(&tapstruct.Threaddetail{
		Domain:     "example.com",
		Name:       "Thread 3",
		State:      tapstruct.NotStarted,
		Iteration:  "2020 Oct",
		CostDirect: 5,
		CostTotal:  5,
		Owner:      pts[0],
	}, []*tapstruct.Threadrow{}, []*tapstruct.Threadrow{})
	if errN2 != nil {
		t.Errorf("Could not create new thread: %v", errN2)
		return
	}
	res0, errG0 := db.GetThreadrowsByPersonteamPlan(pts[0].Email, []string{"2020 Oct"})
	if errG0 != nil {
		t.Errorf("Could not get threads: %v", errG0)
		return
	}
	if len(res0) != 3 {
		t.Errorf("Expected 3 results but got %v", len(res0))
		return
	}
	if res0[0].ID != id0 {
		t.Errorf("The first thread wasn't the first one added. It was %v", res0[0].ID)
		return
	}
	if res0[1].ID != id1 {
		t.Errorf("The second thread wasn't the second one added. It was %v", res0[1].ID)
		return
	}
	if res0[2].ID != id2 {
		t.Errorf("The third thread wasn't the third one added. It was %v", res0[2].ID)
		return
	}
	if res0[0].Order >= res0[1].Order {
		t.Errorf("The first thread's order is not less than the second thread")
		return
	}
	if res0[1].Order >= res0[2].Order {
		t.Errorf("The second thread's order is not less than the third thread")
		return
	}
	if res0[0].Percentile >= res0[1].Percentile {
		t.Errorf("The first thread's percentile is not less than the second thread")
		return
	}
	if res0[1].Percentile >= res0[2].Percentile {
		t.Errorf("The second thread's percentile is not less than the third thread")
		return
	}
	// Move the first thread to the end
	errM := db.MoveThread(&res0[0], After, &res0[2])
	if errM != nil {
		t.Errorf("Could not move thread: %v", errM)
		return
	}
	res1, errG1 := db.GetThreadrowsByPersonteamPlan(pts[0].Email, []string{"2020 Oct"})
	if errG1 != nil {
		t.Errorf("Could not get threads after move: %v", errG1)
		return
	}
	if len(res1) != 3 {
		t.Errorf("Expected 3 results but got %v", len(res1))
		return
	}
	if res1[0].ID != id1 {
		t.Errorf("The first thread wasn't the second one added. It was %v", res1[0].ID)
		return
	}
	if res1[1].ID != id2 {
		t.Errorf("The second thread wasn't the third one added. It was %v", res1[1].ID)
		return
	}
	if res1[2].ID != id0 {
		t.Errorf("The third thread wasn't the first one added. It was %v", res1[2].ID)
		return
	}
}

func setupForNewThread() (DBInterface, []tapstruct.Personteam, error) {
	db, errInit := InitDB()
	if errInit != nil {
		return nil, []tapstruct.Personteam{}, fmt.Errorf("Init returned error: %v", errInit)
	}
	errClear := db.ClearDomain("example.com")
	if errClear != nil {
		return nil, []tapstruct.Personteam{}, fmt.Errorf("Clear domain returned error: %v", errClear)
	}
	pt1 := tapstruct.Personteam{
		Email:      "brandon@example.com",
		Domain:     "example.com",
		Name:       "Brandon",
		Abbrev:     "BR",
		ColorF:     "#ffffff",
		ColorB:     "#1f57cf",
		IterTiming: tapstruct.Monthly,
	}
	pt2 := tapstruct.Personteam{
		Email:      "brenda@example.com",
		Domain:     "example.com",
		Name:       "Brenda",
		Abbrev:     "BR",
		ColorF:     "#ffffff",
		ColorB:     "#1f57cf",
		IterTiming: tapstruct.Quarterly,
	}
	errPT := db.NewPersonteam(&pt1, "")
	if errPT != nil {
		return nil, []tapstruct.Personteam{}, fmt.Errorf("NewPersonteam returned an error: %v", errPT)
	}
	errPT = db.NewPersonteam(&pt2, "")
	if errPT != nil {
		return nil, []tapstruct.Personteam{}, fmt.Errorf("NewPersonteam returned an error: %v", errPT)
	}
	return db, []tapstruct.Personteam{pt1, pt2}, nil
}

func getThread1(owner tapstruct.Personteam) *tapstruct.Threaddetail {
	return &tapstruct.Threaddetail{
		Domain:     "example.com",
		Name:       "Thread 1",
		State:      tapstruct.NotStarted,
		Iteration:  "2020 Oct",
		CostDirect: 10,
		CostTotal:  10,
		Owner:      owner,
	}
}

func getThread2(owner tapstruct.Personteam) *tapstruct.Threaddetail {
	return &tapstruct.Threaddetail{
		Domain:     "example.com",
		Name:       "Thread 2",
		State:      tapstruct.NotStarted,
		Iteration:  "2020 Oct",
		CostDirect: 20,
		CostTotal:  20,
		Owner:      owner,
	}
}
