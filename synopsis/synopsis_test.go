package terst

import (
	"testing"
)

func float32_3() float32 {
	return 3
}

func uint8_5() uint8 {
	return 5
}

func Test(t *testing.T) {
    Terst(t) // Associate terst methods with t (the current testing.T)
    Is("apple", "apple") // Pass
    Is(1, "1") // Pass: 1 is converted to a string before testing
    Is("apple", "orange") // Fail: emits nice-looking diagnostic 
    Compare(1, ">", 0) // Pass

    Is(1, 1.0) // Fail: comparing an integer to a float
    Compare(1, "==", 1.0) // Pass
    Like(1, 1.0) // Pass
	Is(2 + 2, 5.0, "Doubleplusgood")

	Compare(float32_3(), "<", uint8_5())
	// Pass(float32_3() < uint8_5()) Will not compile
}
