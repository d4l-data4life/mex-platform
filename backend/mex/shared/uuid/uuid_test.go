package uuid

import (
	"strings"
	"testing"
)

func TestUUIDs(t *testing.T) {
	for i := 0; i < 100000; i++ {
		u := MustNewV4()
		if strings.Contains(u, "000000000") {
			t.Errorf("invalid UUID: %s", u)
		}
	}
}
