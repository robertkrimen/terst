package terst

import (
	"testing"
)

func Test(t *testing.T) {
	WithTester(t)
	Is("apple", "apple")
	IsNot("apple", "orange")
	Unlike("apple", `pp`)
	Pass(false)
	Fail(true)
    Compare(int64(1), "==", 2)
}
