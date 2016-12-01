package unittest_test

import (
	"github.com/wowplus2/golang/usages/unittest"
	"testing"
)

func TestSum(t *testing.T) {
	s := unittest.Sum(1,2,3)

	if s != 6 {
		t.Error("Wrong result!")
	}
}
