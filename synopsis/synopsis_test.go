package terst

import (
	"fmt"
	"testing"
)

func float32_3() float32 {
	return 3
}

func uint8_5() uint8 {
	return 5
}

func getApple() string {
	return "apple"
}

func getOrange() string {
	return "apple"
}

func get1() int {
	return 1
}

func Test(t *testing.T) {
	Terst(t) // Associate terst methods with t (the current testing.T)

	Is(getApple(), "apple")   // Pass
	Is(getOrange(), "orange") // Fail: emits nice-looking diagnostic 

	Compare(1, ">", 0)    // Pass
	Compare(1, "==", 1.0) // Pass
}

func Test_(t *testing.T) {
	Terst(t) // Associate terst methods with t (the current testing.T)

	Is(get1(), "1") // Pass: 1 is converted to a string before testing
	Is(2+2, float32(5.0), "Doubleplusgood")

	Is(1, 1.0)   // Pass
	Like(1, 1.0) // Fail: comparing an integer to a float

	Compare(float32_3(), "<", uint8_5())
	// Pass(float32_3() < uint8_5()) // Will not compile

	if false {
		fmt.Printf("Xyzzy")
	}
}
