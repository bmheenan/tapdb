package tapdb

import (
	"testing"

	"github.com/bmheenan/tapstruct"
)

func TestSingleThreadInsertAndGetThreadrowByPT(t *testing.T) {
	db, errInit := InitDB()
	if errInit != nil {
		t.Errorf("Init returned error: %v", errInit)
		return
	}
	errClear := db.ClearDomain("example.com")
	if errClear != nil {
		t.Errorf("Clear domain returned error: %v", errClear)
		return
	}
	pt := tapstruct.Personteam{
		Email:      "brandon@example.com",
		Domain:     "example.com",
		Name:       "Brandon",
		Abbrev:     "BR",
		ColorF:     "#ffffff",
		ColorB:     "#1f57cf",
		IterTiming: tapstruct.Monthly,
	}
	errPT := db.NewPersonteam(&pt, "")
	if errPT != nil {
		t.Errorf("NewPersonteam returned an error: %v", errPT)
		return
	}
	_, err := db.NewThread(&tapstruct.Threaddetail{
		Domain:     "example.com",
		Name:       "Example thread",
		State:      tapstruct.NotStarted,
		Iteration:  "2020 Q4",
		CostDirect: 10,
		CostTotal:  10,
		Owner:      pt,
		Order:      0,
		Percentile: 10,
	}, []int64{}, []int64{})
	if err != nil {
		t.Errorf("NewThread returned error: %v", err)
	}
	results, errGet := db.GetThreadrowsByPersonteamPlan("brandon@example.com", []string{"2020 Q4"})
	if errGet != nil {
		t.Errorf("GetThreadrowsByPersonteamPlan returned error: %v", errGet)
		return
	}
	if len(results) != 1 {
		t.Errorf("Expected 1 threadrow, but got %v", len(results))
		return
	}
	if results[0].Name != "Example thread" ||
		results[0].Owner.Email != "brandon@example.com" {
		t.Errorf("Threadrow didn't have the expected data")
		return
	}
}
