package terst

import (
	"testing"
	/*"fmt"*/
	"math"
)

type Xyzzy struct{}

func (self Xyzzy) String() string {
	return "Nothing happens."
}

func TestIsWithoutTerst(t *testing.T) {
	// The following will panic
	// UnTerst()
	// Is(1, 1)
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

func TestCompare(t *testing.T) {
	Terst(t)
	Compare([]string{""}, "==", []string{""})
	Compare([]string{""}, "!=", []string{"Xyzzy"})
	Compare([]string{""}, "{}* ==", []string{""})
	Compare([]string{""}, "{}* !=", []string{"Xyzzy"})
	Compare([]string{""}, "{}~ ==", []string{""})
	Compare([]string{""}, "{}~ !=", []string{"Xyzzy"})
	if false {
		// These fail because you cannot do []type == []type
		Compare([]string{""}, "{}= !=", []string{""})
		Compare([]string{""}, "{}= !=", []string{"Xyzzy"})
	}
	Compare(&Xyzzy{}, "==", &Xyzzy{})
	Compare(&Xyzzy{}, "{}= !=", &Xyzzy{})
	Compare(float32(1.1), "<", int8(2))
	if false {
		// This will not parse/compile because of a type mismatch
		// Pass(float32(1.1) < int(2))
	}
}

func TestCompareOperator(t *testing.T) {
	Terst(t)

	operator := newCompareOperator("#= ==")
	Is(operator.scope, compareScopeEqual)
	Is(operator.comparison, "==")
}

func TestIs(t *testing.T) {
	Terst(t)

	Is(true, "true")
	Is(1, "1")
	Is(Xyzzy{}, "Nothing happens.")
}

func TestPassing(t *testing.T) {
	Terst(t).EnableSelfTesting()

	Is(1, 1)
	Is("apple", "apple")
	IsNot("apple", "orange")
	Is(&Xyzzy{}, &Xyzzy{})

	Compare(1, "==", 1.0)
	Compare(1, "==", 1)
	Compare(&Xyzzy{}, "#* ==", &Xyzzy{})
	Compare("abc", ">=", "abc")
	Compare(1, "#= ==", 1)
}

func TestFailing(t *testing.T) {
	Terst(t).EnableSelfTesting().FailIsPass()

	Equal("apple", "orange")

	Pass(false)
	Fail(true)

	Is("1", 1)
	Is("true", true)
	Is("1", 1)

	Like(1, 1.1)
	Unlike("apple", `pp`)

	Compare(true, ">", false)
	Compare(math.Inf(0), "==", 2)
	Compare("test", "#= ==", int32(1))
	Compare(uint64(math.MaxUint64), "<", int64(math.MinInt32))
	Compare("apple", "==", "banana")
	Compare(false, "==", true)
	Compare(uint(1), "==", int(2))
	Compare(uint(1), "==", 1.1)
	Compare(10, "<", 4.0)
	Compare(6, ">", 6.0)
	Compare("abcd", "<", "abc")
	Compare("ab", ">=", "abc")

}
