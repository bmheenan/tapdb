package tapdb

import (
	"testing"

	"github.com/bmheenan/tapstruct"
)

func TestNewAndGetPersonteam(t *testing.T) {
	_, err := InitDB()
	if err != nil {
		t.Errorf("Init returned error: %v", err)
	}
	pt := &tapstruct.Personteam{}
	pt.Domain = "hi"
}
