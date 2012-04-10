package terst

import (
	"testing"
	/*"fmt"*/
	"math"
)

func init() {
	isTesting = true
	SanityCheck = true
}

type Apple struct{}

func (self Apple) String() string {
	return "This is an Apple object"
}

func TestCompareOperator(t *testing.T) {
	isTesting = false
	Terst(t)

    operator := newCompareOperator("#= ==")
	Is(operator.scope, compareSame)
	Is(operator.operator, "==")
}

func TestPass(t *testing.T) {
	Terst(t)
	isTesting = true
	expectResult = true
	Is(1, 1)
	Compare(1, "==", 1.0)
	Is("apple", "apple")
	IsNot("apple", "orange")
	Compare(1, "==", 1)
	Compare("abc", ">=", "abc")
	Compare(math.Inf(0), "==", 2)

	Compare(1, "#= ==", 1)
	Compare(compareFamily, "#= ==", int32(1))
	Compare("test", "#= ==", int32(1))
	Is("1", 1)

}

func TestIs(t *testing.T) {
	Terst(t)
	isTesting = false

	Is(true, "true")
	Is(1, "1")
	Is(Apple{}, "This is an Apple object")
}

func TestFail(t *testing.T) {
	Terst(t)
	isTesting = true
	expectResult = false
	Unlike("apple", `pp`)
	Like(1, 1.1)
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

	// Is coerce
	Is("true", true)
	Is("1", 1)
}
