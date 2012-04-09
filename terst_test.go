package terst

import (
	"testing"
    "fmt"
    "math"
)

func Test(t *testing.T) {
    a := uint16(0x10F0)
    b := int8(a)
    c := uint32(b)
    fmt.Printf("%v %v %v %v\n", a, b, c, int64(a))


	WithTester(t)
    Compare(false, "<", true)
    return
    Compare("apple", "==", "banana")
    Compare("apple", "==", true)
    Compare(false, "==", true)
    Compare(uint(1), "==", int(2))
    Compare(uint(1), "==", 1.1)

	Equal("apple", "orange")
	Is("apple", "apple")
	IsNot("apple", "orange")
	Unlike("apple", `pp`)
	Pass(false)
	Fail(true)
    Compare(1, "==", 1)
    Compare(10, "<", 4.0)
    Compare(6, ">", 6.0)
    Compare("abcd", "<", "abc")
    Compare("ab", ">=", "abc")
    Compare("abc", ">=", "abc")
    Compare(math.Inf(0), ">", 2)
}
