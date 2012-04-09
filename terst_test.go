package terst

import (
    "testing"
)

func Test(t *testing.T) {
    WithTester(t)
    Is("apple", "apple")
    IsNot("apple", "orange")
    Unlike("apple", `pp`)
}
