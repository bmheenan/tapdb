package tapdb

import (
	"testing"
)

func TestInitDb(t *testing.T) {
	_, err := InitDB()
	if err != nil {
		t.Errorf("Init returned error: %v", err)
		return
	}
}
