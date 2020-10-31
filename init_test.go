package tapdb

import (
	"testing"
)

func TestInitDb(t *testing.T) {
	_, err := Init(getTestCredentials())
	if err != nil {
		t.Errorf("Init returned error: %v", err)
		return
	}
}
