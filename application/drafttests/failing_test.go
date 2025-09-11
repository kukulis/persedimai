package drafttests

import "testing"

func TestFailing(t *testing.T) {
	t.Errorf("Failing test")
}
