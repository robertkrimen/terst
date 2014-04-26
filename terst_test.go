package terst

import (
	"fmt"
	"math"
	"regexp"
	"strings"
	"testing"
)

func fail() {
	Is("abc", ">", regexp.MustCompile(`abc`))
}

func TestSynopsis(t *testing.T) {
	Terst(t, func() {
		Is(1, ">", 0)
		Is(2+2, "!=", 5)
		Is("Nothing happens.", "=~", `ing(\s+)happens\.$`)
	})
}

func Test(t *testing.T) {
	Terst(t, func() {
		Is("abc", "abc")

		Is("abc", "==", "abc")

		Is("abc", "!=", "def")

		Is(0, "!=", 3.14159)

		Is(math.NaN(), math.NaN())

		Is(nil, nil)

		var abc map[string]string

		Is(abc, nil)

		Is([]byte("Nothing happens."), "=~", `ing(\s+)happens\.$`)

		{
			err := IsErr(abc, "!=", nil)
			Is(err, "!=", nil)

			if false {
				func() {
					fail()
				}()
			}

			err = IsErr("abc", ">", "def")
			Is(err, "!=", nil)
		}
	})
}

func Test_findTestFunc(t *testing.T) {
	Terst(t, func() {
		cl := Caller()
		Is(cl.TestFunc().Name(), "github.com/robertkrimen/terst.Test_findTestFunc")
	})
}

func Test_IsErr(t *testing.T) {

	{
		// NaN == NaN
		result, err := compareNumber(math.NaN(), math.NaN())
		if err != nil {
			t.Errorf("NaN =? NaN: %s", err.Error())
		} else if result != 0 {
			t.Errorf("NaN =? NaN: %d", result)
		}
	}

	if err := IsErr(math.NaN(), math.NaN()); err != nil {
		t.Error(err)
	}

	if err := IsErr("", ""); err != nil {
		t.Error(err)
	}

	if err := IsErr("abc", ""); err == nil {
		t.Error(err)
	}

	if err := IsErr(1, "<=", 1); err != nil {
		t.Error(err)
	}

	if err := IsErr(0, "<", 1); err != nil {
		t.Error(err)
	}

	if err := IsErr(int64(math.MaxInt64), "<", uint64(math.MaxUint64)); err != nil {
		t.Error(err)
	}

	if err := IsErr(int64(math.MaxInt64), ">=", uint64(math.MaxUint64)); err == nil {
		t.Error(err)
	}

	if err := IsErr(fmt.Errorf("abc"), "abc"); err != nil {
		t.Error(err)
	}

	if err := IsErr(fmt.Errorf("abc"), ""); err == nil {
		t.Error(err)
	}

	test := func(arguments ...interface{}) bool {
		expect := arguments[len(arguments)-1]
		arguments = arguments[:len(arguments)-1]

		input0 := arguments[0]
		input1 := arguments[len(arguments)-1]

		err := IsErr(arguments...)
		if expect == nil && err == nil {
			return true
		} else if expect != nil && err != nil {
			expect := expect.(string)
			got := strings.Join(strings.Fields(err.Error()), " ")
			if got == expect {
				return true
			}
			t.Errorf("\nerr != expect:\n          got: %v\n       expect: %v", got, expect)
		} else if err == nil {
			t.Errorf("\nerr == nil: %v\n       got: %v\n    expect: %v", expect, input0, input1)
		} else {
			got := strings.Join(strings.Fields(err.Error()), " ")
			t.Errorf("\nerr != nil: %v\n       got: %v\n    expect: %v", got, input0, input1)
		}
		return false
	}

	test("", "", nil)

	test("abc", "def", "FAIL (==) got: abc expected: def")

	test("abc", ">", "def", "INVALID (>): got: abc (string) expected: def (string)")

	test(1, ">", 0, nil)

	test(1, "<", 0, "FAIL (<) got: 1 (int) expected: < 0 (int)")

	test([]byte("def"), "=~", `abc$`, "FAIL (=~) got: def [100 101 102] ([]uint8=slice) expected: abc$")

}
