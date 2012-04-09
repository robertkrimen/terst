package terst

import (
	"testing"
)

func Test(t *testing.T) {
	Terst(t)
    Terst(t) // Associate terst methods with t (the current testing.T)
    Is("apple", "apple") // Pass
    Is(1, "1") // Pass: 1 is converted to a string before testing
    Is("apple", "orange") // Fail: emits nice-looking diagnostic 
    Compare(1, ">", 0) // Pass

    Is(1, 1.0) // Fail: comparing an integer to a float
    Compare(1, "==", 1.0) // Pass
    Like(1, 1.0) // Pass
}
