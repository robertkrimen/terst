package terst

import (
	"fmt"
	"math"
	"math/big"
	"reflect"
	"regexp"
	"runtime"
	/*"strconv"*/
	"strings"
	"testing"
	"unsafe"
)

var (
    SanityCheck bool = true

    expectResult bool
    isTesting bool = false
)

func (self *Tester) hadResult(result bool, test *test, onFail func()) bool {
    if isTesting {
        if expectResult != result {
            self.Log(fmt.Sprintf("Expect %v but got %v (%v) (%v) (%v)\n", expectResult, result, test.kind, test.have, test.want))
            onFail()
            self.fail()
        }
        return result
    }
    if !result {
        onFail()
        self.fail()
    }
    return result
}

// Pass

func Pass(have bool, arguments ...interface{}) bool {
	return OurTester().AtPass(1, have, arguments...)
}

func (self *Tester) Pass(have bool, arguments ...interface{}) bool {
	return self.AtPass(1, have, arguments...)
}

func (self *Tester) AtPass(callDepth int, have bool, arguments ...interface{}) bool {
	return self.atPassOrFail(true, callDepth+1, have, arguments...)
}

// Fail

func Fail(have bool, arguments ...interface{}) bool {
	return OurTester().AtFail(1, have, arguments...)
}

func (self *Tester) Fail(have bool, arguments ...interface{}) bool {
	return self.AtFail(1, have, arguments...)
}

func (self *Tester) AtFail(callDepth int, have bool, arguments ...interface{}) bool {
	return self.atPassOrFail(false, callDepth+1, have, arguments...)
}

func (self *Tester) atPassOrFail(want bool, callDepth int, have bool, arguments ...interface{}) bool {
	kind := "Pass"
	if want == false {
		kind = "Fail"
	}
	test := newTest(kind, callDepth+1, have, want, arguments)
	didPass := have == want
    return self.hadResult(didPass, test, func(){
		self.Log(self.failMessageForPass(test))
    })
}

// Equal

func Equal(have, want interface{}, arguments ...interface{}) bool {
	return OurTester().AtEqual(1, have, want, arguments...)
}

func (self *Tester) Equal(have, want interface{}, arguments ...interface{}) bool {
	return self.AtEqual(1, have, want, arguments...)
}

func (self *Tester) AtEqual(callDepth int, have, want interface{}, arguments ...interface{}) bool {
	test := newTest("==", callDepth+1, have, want, arguments)
	didPass := have == want
    return self.hadResult(didPass, test, func(){
		self.Log(self.failMessageForEqual(test))
    })
}

// Unequal

func Unequal(have, want interface{}, arguments ...interface{}) bool {
	return OurTester().AtUnequal(1, have, want, arguments...)
}

func (self *Tester) Unequal(have, want interface{}, arguments ...interface{}) bool {
	return self.AtUnequal(1, have, want, arguments...)
}

func (self *Tester) AtUnequal(callDepth int, have, want interface{}, arguments ...interface{}) bool {
	test := newTest("!=", callDepth+1, have, want, arguments)
	didPass := have != want
    return self.hadResult(didPass, test, func(){
		self.Log(self.failMessageForIs(test))
    })
}

// Is

func Is(have, want interface{}, arguments ...interface{}) bool {
	return OurTester().AtIs(1, have, want, arguments...)
}

func (self *Tester) Is(have, want interface{}, arguments ...interface{}) bool {
	return self.AtIs(1, have, want, arguments...)
}

func (self *Tester) AtIs(callDepth int, have, want interface{}, arguments ...interface{}) bool {
	test := newTest("Is", callDepth+1, have, want, arguments)
	didPass := have == want
    return self.hadResult(didPass, test, func(){
		self.Log(self.failMessageForIs(test))
    })
}

// IsNot

func IsNot(have, want interface{}, arguments ...interface{}) bool {
	return OurTester().AtIsNot(1, have, want, arguments...)
}

func (self *Tester) IsNot(have, want interface{}, arguments ...interface{}) bool {
	return self.AtIsNot(1, have, want, arguments...)
}

func (self *Tester) AtIsNot(callDepth int, have, want interface{}, arguments ...interface{}) bool {
	test := newTest("IsNot", callDepth+1, have, want, arguments)
	didPass := have != want
    return self.hadResult(didPass, test, func(){
		self.Log(self.failMessageForIs(test))
    })
}

// Like

func Like(have, want interface{}, arguments ...interface{}) bool {
	return OurTester().AtLike(1, have, want, arguments...)
}

func (self *Tester) Like(have, want interface{}, arguments ...interface{}) bool {
	return self.AtLike(1, have, want, arguments...)
}

func (self *Tester) AtLike(callDepth int, have, want interface{}, arguments ...interface{}) bool {
	return self.atLikeOrUnlike(true, callDepth+1, have, want, arguments...)
}

// Unlike

func Unlike(have, want interface{}, arguments ...interface{}) bool {
	return OurTester().AtUnlike(1, have, want, arguments...)
}

func (self *Tester) Unlike(have, want interface{}, arguments ...interface{}) bool {
	return self.AtUnlike(1, have, want, arguments...)
}

func (self *Tester) AtUnlike(callDepth int, have, want interface{}, arguments ...interface{}) bool {
	return self.atLikeOrUnlike(false, callDepth+1, have, want, arguments...)
}

func (self *Tester) atLikeOrUnlike(wantLike bool, callDepth int, have, want interface{}, arguments ...interface{}) bool {
	test := newTest("Like", callDepth+1, have, want, arguments)
    if !wantLike {
        test.kind = "Unlike"
    }
    didPass := true
	switch want0 := want.(type) {
	case string:
		haveString := ToString(have)
		didPass, error := regexp.Match(want0, []byte(haveString))
		if !wantLike {
			didPass = !didPass
		}
		if error != nil {
			panic("regexp.Match(" + want0 + ", ...): " + error.Error())
		}
        want = fmt.Sprintf("(?:%v)", want) // Make it look like a regular expression
        return self.hadResult(didPass, test, func(){
            self.Log(self.failMessageForMatch(test, ToString(have), ToString(want), wantLike))
        })
    }
    operator := "=="
    if !wantLike {
        operator = "!="
    }
    didPass = compare(have, operator, want)
    return self.hadResult(didPass, test, func(){
        self.Log(self.failMessageForLike(test, ToString(have), ToString(want), wantLike))
    })
}

// Compare 

func Compare(have interface{}, operator string, want interface{}, arguments ...interface{}) bool {
	return OurTester().AtCompare(1, have, operator, want, arguments...)
}

func (self *Tester) Compare(have interface{}, operator string, want interface{}, arguments ...interface{}) bool {
	return self.AtCompare(1, have, operator, want, arguments...)
}

func (self *Tester) AtCompare(callDepth int, left interface{}, operator string, right interface{}, arguments ...interface{}) bool {
    test := newTest("Compare " + operator, callDepth+1, left, right, arguments)
    test.operator = operator
    didPass := compare(left, operator, right)
    return self.hadResult(didPass, test, func(){
        self.Log(self.failMessageForCompare(test))
    })
}

func compare(left interface{}, operator string, right interface{}) bool {
    pass := true
    comparator := newComparator(left, right)
    switch operator {
    case "==":
        pass = comparator.isEqual()
    case "!=":
        pass = !comparator.isEqual()
    default:
        if !comparator.hasOrder() {
            panic(fmt.Errorf("Comparison (%v) %v (%v) is invalid", left, operator, right))
        }
        switch operator {
        case "<":
            pass = comparator.compare() == -1
        case "<=":
            pass = comparator.compare() <= 0
        case ">":
            pass = comparator.compare() == 1
        case ">=":
            pass = comparator.compare() >= 0
        default:
            panic(fmt.Errorf("Compare operator (%v) is invalid", operator))
        }
    }
    return pass
}

// Compare / Comparator

type kind int
const (
    isInvalid = iota
    isInteger
    isFloat
    isString
    isBoolean
)

func comparatorValue(value interface{}) (reflect.Value, int) {
    reflectValue := reflect.ValueOf(value)
    kind := isInvalid
    switch value.(type) {
    case int, int8, int16, int32, int64:
        kind = isInteger
    case uint, uint8, uint16, uint32, uint64:
        kind = isInteger
    case float32, float64:
        kind = isFloat
    case string:
        kind = isString
    case bool:
        kind = isBoolean
    }
    return reflectValue, kind
}

func toFloat(value reflect.Value) float64 {
    switch result := value.Interface().(type) {
    case int, int8, int16, int32, int64:
        return float64(value.Int())
    case uint, uint8, uint16, uint32, uint64:
        return float64(value.Uint())
    case float32, float64:
        return float64(value.Float())
    default:
        panic(fmt.Errorf("toFloat( %v )", result))
    }
    panic(0)
}

func toInteger(value reflect.Value) *big.Int {
    switch result := value.Interface().(type) {
    case int, int8, int16, int32, int64:
        return big.NewInt(value.Int())
    case uint, uint8, uint16, uint32, uint64:
        yield := big.NewInt(0)
        yield.SetString(fmt.Sprintf("%v", value.Uint()), 10)
        return yield
    default:
        panic(fmt.Errorf("toInteger( %v )", result))
    }
    panic(0)
}

func toString(value reflect.Value) string {
    switch result := value.Interface().(type) {
    case string:
        return result
    default:
        panic(fmt.Errorf("toString( %v )", result))
    }
    panic(0)
}

func toBoolean(value reflect.Value) bool {
    switch result := value.Interface().(type) {
    case bool:
        return result
    default:
        panic(fmt.Errorf("toBoolean( %v )", result))
    }
    panic(0)
}

type ofComparator interface {
    compare() int
    isEqual() bool
    hasOrder() bool
}

type baseComparator struct {
    order bool
}
func (self *baseComparator) compare() int {
    panic(fmt.Errorf("Attempt to .compare() on type without order"))
}
func (self *baseComparator) hasOrder() bool {
    return self.order
}

type floatComparator struct {
    *baseComparator
    left float64
    right float64
}
func (self *floatComparator) compare() int {
    if self.left == self.right {
        return 0
    } else if self.left < self.right {
        return -1
    }
    return 1
}
func (self *floatComparator) isEqual() bool {
    return self.left == self.right
}

type integerComparator struct {
    *baseComparator
    left *big.Int
    right *big.Int
}
func (self *integerComparator) compare() int {
    return self.left.Cmp(self.right)
}
func (self *integerComparator) isEqual() bool {
    return 0 == self.left.Cmp(self.right)
}

type stringComparator struct {
    *baseComparator
    left string
    right string
}
func (self *stringComparator) compare() int {
    if self.left == self.right {
        return 0
    } else if self.left < self.right {
        return -1
    }
    return 1
}
func (self *stringComparator) isEqual() bool {
    return self.left == self.right
}

type booleanComparator struct {
    *baseComparator
    left bool
    right bool
}
func (self *booleanComparator) isEqual() bool {
    return self.left == self.right
}

func newComparator(left interface{}, right interface{}) ofComparator {
    leftValue, leftKind := comparatorValue(left)
    rightValue, rightKind := comparatorValue(right)
    if false {
    } else if leftKind == isFloat || rightKind == isFloat {
        return &floatComparator{
            &baseComparator{true},
            toFloat(leftValue),
            toFloat(rightValue),
        }
    } else if leftKind == isInteger || rightKind == isInteger {
        return &integerComparator{
            &baseComparator{true},
            toInteger(leftValue),
            toInteger(rightValue),
        }
    } else if leftKind == isString {
        return &stringComparator{
            &baseComparator{true},
            toString(leftValue),
            toString(rightValue),
        }
    } else if leftKind == isBoolean {
        return &booleanComparator{
            &baseComparator{false},
            toBoolean(leftValue),
            toBoolean(rightValue),
        }
    }
    panic(fmt.Errorf("Comparing (%v) to (%v) not implemented", left, right))
}

// failMessage*

func (self *Tester) failMessageForPass(test *test) string {
	return self.FormatMessage(`
        %s:%d: %s 
           Failed test (%s)
                  got: %s
             expected: %s
    `, test.file, test.line, test.Description(), test.kind, ToString(test.have), ToString(test.want))
}

func (self *Tester) failMessageForCompare(test *test) string {
	return self.FormatMessage(`
        %s:%d: %s 
           Failed test (%s)
                  %s
                       %s
                  %s
    `, test.file, test.line, test.Description(), test.kind, ToString(test.have), test.operator, ToString(test.want))
}

func (self *Tester) failMessageForEqual(test *test) string {
    return self.failMessageForIs(test)
}

func (self *Tester) failMessageForIs(test *test) string {
	return self.FormatMessage(`
        %s:%d: %s 
           Failed test (%s)
                  got: %s
             expected: %s
    `, test.file, test.line, test.Description(), test.kind, test.have, test.want)
}

func (self *Tester) failMessageForMatch(test *test, have, want string, wantMatch bool) string {
    expect := "  like"
	if !wantMatch {
        expect = "unlike"
	}
	return self.FormatMessage(`
        %s:%d: %s 
           Failed test (%s)
                  got: %s
               %s: %s
    `, test.file, test.line, test.Description(), test.kind, have, expect, want)
}

func (self *Tester) failMessageForLike(test *test, have, want string, wantLike bool) string {
	if !wantLike {
        want = "Anything else"
	}
	return self.FormatMessage(`
        %s:%d: %s 
           Failed test (%s)
                  got: %s
             expected: %s
    `, test.file, test.line, test.Description(), test.kind, have, want)
}

// ...

type Tester struct {
	TestingT *testing.T
    TestEntry uintptr
    sanityCheck bool
}

var _ourTester *Tester = nil

func testFunctionEntry() uintptr {
    height := 2
    for {
	    functionPC, _, _, good := runtime.Caller(height)
	    function := runtime.FuncForPC(functionPC)
	    functionName := function.Name()
        if !good {
            return 0
        }
		if index := strings.LastIndex(functionName, ".Test"); index >= 0 {
            // Assume we have an instance of TestXyzzy in a _test file
            return function.Entry()
		}
        height += 1
    }
    return 0
}

func Terst(arguments ...interface{}) *Tester {
    if len(arguments) == 0 {
        _ourTester = nil
    } else {
	    _ourTester = NewTester(arguments[0].(*testing.T))
        _ourTester.sanityCheck = SanityCheck
        if _ourTester.sanityCheck {
            _ourTester.TestEntry = testFunctionEntry()
        }
    }
	return _ourTester
}

func UnTerst() {
    _ourTester = nil
}

func OurTester() *Tester {
	if _ourTester == nil {
		panic("_ourTester == nil")
	}
	return _ourTester.checkSanity()
}

func HaveTester() bool {
	return _ourTester != nil
}

// Tester

func NewTester(t *testing.T) *Tester {
	return &Tester{t, 0, true}
}

func (self *Tester) FormatMessage(format string, arguments ...interface{}) string {
	message := fmt.Sprintf(format, arguments...)
	message = strings.TrimLeft(message, "\n")
	message = strings.TrimRight(message, " \n")
	return message + "\n\n"
}

func (self *Tester) Log(moreOutput string) {
	outputValue := reflect.ValueOf(self.TestingT).Elem().FieldByName("output")
	output := outputValue.Bytes()
	output = append(output, moreOutput...)
	*(*[]byte)(unsafe.Pointer(outputValue.UnsafeAddr())) = output
}

func (self *Tester) fail() {
    self.TestingT.Fail()
}

func (self *Tester) checkSanity() *Tester {
    if self.sanityCheck && self.TestEntry != 0 {
        foundEntry := testFunctionEntry()
        if self.TestEntry != foundEntry {
            panic(fmt.Errorf("TestEntry(%v) does not match foundEntry(%v): Did you call Terst when entering a new Test* function?", self.TestEntry, foundEntry))
        }
    }
    return self
}

// test

type test struct {
	kind      string
	have      interface{}
	want      interface{}
	arguments []interface{}
    operator string // Only used for Compare 

	file       string
	line       int
	functionPC uintptr
	function   string
}

func newTest(kind string, callDepth int, have, want interface{}, arguments ...interface{}) *test {
	file, line, functionPC, function, _ := AtFileLineFunction(callDepth + 1)
    operator := ""
	return &test{kind, have, want, arguments, operator, file, line, functionPC, function}
}

func (self *test) Description() string {
	description := ""
	if len(self.arguments) > 0 {
		description = fmt.Sprintf("%s", self.arguments...)
	}
	return description
}

func AtFileLineFunction(callDepth int) (string, int, uintptr, string, bool) {
	functionPC, file, line, good := runtime.Caller(callDepth + 1)
	function := runtime.FuncForPC(functionPC).Name()
	if good {
		if index := strings.LastIndex(file, "/"); index >= 0 {
			file = file[index+1:]
		} else if index = strings.LastIndex(file, "\\"); index >= 0 {
			file = file[index+1:]
		}
		if index := strings.LastIndex(function, ".Test"); index >= 0 {
			function = function[index+1:]
		}
	} else {
		functionPC = 0
		file = "?"
		line = 1
	}
	return file, line, functionPC, function, good
}

// Conversion

func integerValue(value interface{}) int64 {
    return reflect.ValueOf(value).Int()
}

func unsignedIntegerValue(value interface{}) uint64 {
    return reflect.ValueOf(value).Uint()
}

func floatValue(value interface{}) float64 {
    return reflect.ValueOf(value).Float()
}

func ToString(value interface{}) string {
	switch value0 := value.(type) {
	case bool:
		return fmt.Sprintf("%v", value)
    case int, int8, int16, int32, int64:
		return fmt.Sprintf("%v", value)
    case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%v", value)
	case string:
		return fmt.Sprintf("%v", value)
    case float32:
		return fmt.Sprintf("%v", value)
    case float64:
		if math.IsNaN(value0) {
			return "NaN"
		} else if math.IsInf(value0, 0) {
			return "Infinity"
		}
		return fmt.Sprintf("%v", value)
	}
    return fmt.Sprintf("%v", value)
}
