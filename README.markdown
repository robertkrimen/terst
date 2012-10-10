# terst
--
Package terst is a terse (terst = test + terse), easy-to-use testing library for Go.

terst is compatible with (and works via) the standard testing package: http://golang.org/pkg/testing

	import (
		"testing"
		. "github.com/robertkrimen/terst"
	)

	func Test(t *testing.T) {
		Terst(t) // Associate terst methods with t (the current testing.T)

		Is(getApple(), "apple") // Pass
		Is(getOrange(), "orange") // Fail: emits nice-looking diagnostic

		Compare(1, ">", 0) // Pass
		Compare(1, "==", 1.0) // Pass
	}

	func getApple() string {
		return "apple"
	}

	func getOrange() string {
		return "apple" // Intentional mistake
	}

At the top of your testing function, call Terst(), passing the testing.T you receive as the first argument:

	func TestExample(t *testing.T) {
		Terst(t)
		...
	}

After you initialize with the given *testing.T, you can use the following to test:

	Is
	IsNot
	Equal
	Unequal
	IsTrue
	IsFalse
	Like
	Unlike
	Compare

Each of the methods above can take an additional (optional) argument,
which is a string describing the test. If the test fails, this
description will be included with the test output For example:

	Is(2 + 2, float32(5), "This result is Doubleplusgood")

	--- FAIL: Test (0.00 seconds)
		test.go:17: This result is Doubleplusgood
			Failed test (Is)
			     got: 4 (int)
			expected: 5 (float32)

### Future

	- Add Catch() for testing panic()
	- Add Same() for testing via .DeepEqual && == (without panicking?)
	- Add StrictCompare to use {}= scoping
	- Add BigCompare for easier math/big.Int testing?
	- Support the complex type in Compare()
	- Equality test for NaN?
	- Better syntax for At*
	- Need IsType/TypeIs

## Usage

#### func  Compare

```go
func Compare(have interface{}, operator string, want interface{}, description ...interface{}) bool
```
Compare will compare <have> to <want> with the given operator. The operator can
be one of the following:

    ==
    !=
    <
    <=
    >
    >=

Compare is not strict when comparing numeric types, and will make a best effort
to promote <have> and <want> to the same type.

Compare will promote int and uint to big.Int for testing against each other.

Compare will promote int, uint, and float to float64 for float testing.

For example:

    Compare(float32(1.0), "<", int8(2)) // A valid test
    result := float32(1.0) < int8(2) // Will not compile because of the type mismatch

#### func  Equal

```go
func Equal(have, want interface{}, description ...interface{}) bool
```
Equal tests have against want via ==:

    Equal(have, want) // Pass if have == want

No special coercion or type inspection is done.

If the type is incomparable (e.g. type mismatch) this will panic.

#### func  Fail

```go
func Fail(description ...interface{}) bool
```
Fail will fail immediately, reporting a test failure with the (optional)
description

#### func  Is

```go
func Is(have, want interface{}, description ...interface{}) bool
```
Is tests <have> against <want> in different ways, depending on the type of
<want>.

If <want> is a string, then it will first convert <have> to a string before
doing the comparison:

    Is(fmt.Sprintf("%v", have), want) // Pass if have == want

Otherwise, Is is a shortcut for:

    Compare(have, "==", want)

If <want> is a slice, struct, or similar, Is will perform a reflect.DeepEqual()
comparison.

#### func  IsFalse

```go
func IsFalse(have bool, description ...interface{}) bool
```
IsFalse tests if <have> is false.

#### func  IsNot

```go
func IsNot(have, want interface{}, description ...interface{}) bool
```
IsNot tests <have> against <want> in different ways, depending on the type of
<want>.

If <want> is a string, then it will first convert <have> to a string before
doing the comparison:

    IsNot(fmt.Sprintf("%v", have), want) // Pass if have != want

Otherwise, Is is a shortcut for:

    Compare(have, "!=", want)

If <want> is a slice, struct, or similar, Is will perform a reflect.DeepEqual()
comparison.

#### func  IsTrue

```go
func IsTrue(have bool, description ...interface{}) bool
```
IsTrue tests if <have> is true.

#### func  Like

```go
func Like(have, want interface{}, description ...interface{}) bool
```
Like tests <have> against <want> in different ways, depending on the type of
<want>.

If <want> is a string, then it will first convert <have> to a string before
doing a regular expression comparison:

    Like(fmt.Sprintf("%v", have), want) // Pass if regexp.Match(want, have)

Otherwise, Like is a shortcut for:

    Compare(have, "{}~ ==", want)

If <want> is a slice, struct, or similar, Like will perform a
reflect.DeepEqual() comparison.

#### func  Unequal

```go
func Unequal(have, want interface{}, description ...interface{}) bool
```
Unequal tests have against want via !=:

    Unequal(have, want) // Pass if have != want

No special coercion or type inspection is done.

If the type is incomparable (e.g. type mismatch) this will panic.

#### func  Unlike

```go
func Unlike(have, want interface{}, description ...interface{}) bool
```
Unlike tests <have> against <want> in different ways, depending on the type of
<want>.

If <want> is a string, then it will first convert <have> to a string before
doing a regular expression comparison:

    Unlike(fmt.Sprintf("%v", have), want) // Pass if !regexp.Match(want, have)

Otherwise, Unlike is a shortcut for:

    Compare(have, "{}~ !=", want)

If <want> is a slice, struct, or similar, Unlike will perform a
reflect.DeepEqual() comparison.

#### type Tester

```go
type Tester struct {
	TestingT *testing.T
	// contains filtered or unexported fields
}
```


#### func  Terst

```go
func Terst(terst ...interface{}) *Tester
```

    Terst(*testing.T)
Create a new terst Tester and return it. Associate calls to Is, Compare, Like,
etc. with the newly created terst.

    Terst()

Return the current Tester (if any).

    Terst(nil)

Clear out the current Tester (if any).

#### func (*Tester) Compare

```go
func (self *Tester) Compare(have interface{}, operator string, want interface{}, description ...interface{}) bool
```

#### func (*Tester) Equal

```go
func (self *Tester) Equal(have, want interface{}, description ...interface{}) bool
```

#### func (*Tester) Fail

```go
func (self *Tester) Fail(description ...interface{}) bool
```
Fail will fail immediately, reporting a test failure with the (optional)
description

#### func (*Tester) Focus

```go
func (self *Tester) Focus()
```
Focus will focus the entry point of the test to the current method.

This is important for test failures in getting feedback on which line was at
fault.

Consider the following scenario:

    func testingMethod( ... ) {
    	Is( ..., ... )
    }
    func TestExample(t *testing.T) {
    	Terst(t)
    	testingMethod( ... )
    	testingMethod( ... ) // If something in testingMethod fails, this line number will come up
    	testingMethod( ... )
    }

By default, when a test fails, terst will report the outermost line that led to
the failure. Usually this is what you want, but if you need to drill down, you
can by inserting a special call at the top of your testing method:

    func testingMethod( ... ) {
    	Terst().Focus() // Grab the global Tester and tell it to focus on this method
    	Is( ..., ... ) // Now if this test fails, this line number will come up
    }

#### func (*Tester) Is

```go
func (self *Tester) Is(have, want interface{}, description ...interface{}) bool
```

#### func (*Tester) IsFalse

```go
func (self *Tester) IsFalse(have bool, description ...interface{}) bool
```
IsFalse tests if <have> is false.

#### func (*Tester) IsNot

```go
func (self *Tester) IsNot(have, want interface{}, description ...interface{}) bool
```

#### func (*Tester) IsTrue

```go
func (self *Tester) IsTrue(have bool, description ...interface{}) bool
```
IsTrue tests if <have> is true.

#### func (*Tester) Like

```go
func (self *Tester) Like(have, want interface{}, description ...interface{}) bool
```

#### func (*Tester) Log

```go
func (self *Tester) Log(output string)
```
Log is a utility method that will append the given output to the normal output
stream.

#### func (*Tester) Unequal

```go
func (self *Tester) Unequal(have, want interface{}, description ...interface{}) bool
```

#### func (*Tester) Unlike

```go
func (self *Tester) Unlike(have, want interface{}, description ...interface{}) bool
```


