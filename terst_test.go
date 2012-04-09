package terst

import (
	"testing"
    /*"fmt"*/
    "math"
)

func init() {
    isTesting = true
}

func TestPass(t *testing.T) {
	WithTester(t)
    expectResult = true
    Is(1, 1)
    Compare(1, "==", 1.0)
	Is("apple", "apple")
	IsNot("apple", "orange")
    Compare(1, "==", 1)
    Compare("abc", ">=", "abc")
    Compare(math.Inf(0), ">", 2)
    Like(1, 1.1)
}

func TestFail(t *testing.T) {
	WithTester(t)
    expectResult = false
	Unlike("apple", `pp`)
    Like(1, 1.1);
    Compare(uint64(math.MaxUint64), "<", int64(math.MinInt32))
    Compare("apple", "==", "banana")
    Compare(false, "==", true)
    Compare(uint(1), "==", int(2))
    Compare(uint(1), "==", 1.1)
	Equal("apple", "orange")
	Pass(false)
	Fail(true)
    Compare(10, "<", 4.0)
    Compare(6, ">", 6.0)
    Compare("abcd", "<", "abc")
    Compare("ab", ">=", "abc")
}
