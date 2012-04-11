package terst

import (
	"testing"
	/*"fmt"*/
	"math"
)

type Apple struct{}

func (self Apple) String() string {
	return "This is an Apple object"
}

func TestNewCompareOperator(t *testing.T) {
	Terst(t)
	test := func(input string, expect []string) {
		result := newCompareOperatorRE.FindStringSubmatch(input)
		Is(result[1:], expect)
	}
	test("#= ==", []string{"#=", "=="})
	test("  {}* ==", []string{"{}*", "=="})
	test("  {}* ==  ", []string{"{}*", "=="})
	test("   ==  ", []string{"", "=="})
}

func TestCompareOperator(t *testing.T) {
	Terst(t)

    operator := newCompareOperator("#= ==")
	Is(operator.scope, compareScopeEqual)
	Is(operator.comparison, "==")
}

func TestPass(t *testing.T) {
	Terst(t).EnableSelfTesting()
	Is(1, 1)
	Compare(1, "==", 1.0)
	Is("apple", "apple")
	IsNot("apple", "orange")
	Compare(1, "==", 1)
	Compare(&Apple{}, "#* ==", &Apple{})
	Is(&Apple{}, &Apple{})
	Compare("abc", ">=", "abc")
	Compare(1, "#= ==", 1)
}

func TestIs(t *testing.T) {
	Terst(t)

	Is(true, "true")
	Is(1, "1")
	Is(Apple{}, "This is an Apple object")
}

func TestFail(t *testing.T) {
	Terst(t).EnableSelfTesting().FailIsPass()
	Unlike("apple", `pp`)
	Like(1, 1.1)
	Compare(true, ">", false)
	Compare(math.Inf(0), "==", 2)
	Is("1", 1)
	Compare("test", "#= ==", int32(1))

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
