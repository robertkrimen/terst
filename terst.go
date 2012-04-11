package terst

import (
	"fmt"
	"math/big"
	"reflect"
	"regexp"
	"runtime"
	/*"strconv"*/
	"strings"
	"testing"
	"unsafe"
)

func (self *Tester) hadResult(result bool, test *test, onFail func()) bool {
	if self.SelfTesting {
		expect := true
		if self.FailIsPassing {
			expect = false
		}
		if expect != result {
			self.Log(fmt.Sprintf("Expect %v but got %v (%v) (%v) (%v)\n", expect, result, test.kind, test.have, test.want))
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
	return self.hadResult(didPass, test, func() {
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
	return self.hadResult(didPass, test, func() {
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
	return self.hadResult(didPass, test, func() {
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
	return self.atIsOrIsNot(true, callDepth+1, have, want, arguments...)
}

// IsNot

func IsNot(have, want interface{}, arguments ...interface{}) bool {
	return OurTester().AtIsNot(1, have, want, arguments...)
}

func (self *Tester) IsNot(have, want interface{}, arguments ...interface{}) bool {
	return self.AtIsNot(1, have, want, arguments...)
}

func (self *Tester) AtIsNot(callDepth int, have, want interface{}, arguments ...interface{}) bool {
	return self.atIsOrIsNot(false, callDepth+1, have, want, arguments...)
}

func (self *Tester) atIsOrIsNot(wantIs bool, callDepth int, have, want interface{}, arguments ...interface{}) bool {
	test := newTest("Is", callDepth+1, have, want, arguments)
	if !wantIs {
		test.kind = "IsNot"
	}
	didPass := false
	switch want.(type) {
	case string:
		didPass = stringValue(have) == want
	default:
		didPass, _ = compare(have, "{}* ==", want)
	}
	if !wantIs {
		didPass = !didPass
	}
	return self.hadResult(didPass, test, func() {
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
	didPass := false
	switch want0 := want.(type) {
	case string:
		haveString := stringValue(have)
		didPass, error := regexp.Match(want0, []byte(haveString))
		if !wantLike {
			didPass = !didPass
		}
		if error != nil {
			panic("regexp.Match(" + want0 + ", ...): " + error.Error())
		}
		want = fmt.Sprintf("(?:%v)", want) // Make it look like a regular expression
		return self.hadResult(didPass, test, func() {
			self.Log(self.failMessageForMatch(test, stringValue(have), stringValue(want), wantLike))
		})
	}
	didPass, operator := compare(have, "{}~ ==", want)
	if !wantLike {
		didPass = !didPass
	}
	test.operator = operator
	return self.hadResult(didPass, test, func() {
		self.Log(self.failMessageForLike(test, stringValue(have), stringValue(want), wantLike))
	})
}

// Compare 

func Compare(have interface{}, operator string, want interface{}, arguments ...interface{}) bool {
	return OurTester().AtCompare(1, have, operator, want, arguments...)
}

func (self *Tester) Compare(have interface{}, operator string, want interface{}, arguments ...interface{}) bool {
	return self.AtCompare(1, have, operator, want, arguments...)
}

func (self *Tester) AtCompare(callDepth int, left interface{}, operatorString string, right interface{}, arguments ...interface{}) bool {
	operatorString = strings.TrimSpace(operatorString)
	test := newTest("Compare "+operatorString, callDepth+1, left, right, arguments)
	didPass, operator := compare(left, operatorString, right)
	test.operator = operator
	return self.hadResult(didPass, test, func() {
		self.Log(self.failMessageForCompare(test))
	})
}

type (
	compareScope int
)

const (
	compareScopeEqual compareScope = iota
	compareScopeTilde
	compareScopeAsterisk
)

type compareOperator struct {
	scope      compareScope
	comparison string
}

var newCompareOperatorRE *regexp.Regexp = regexp.MustCompile(`^\s*(?:((?:{}|#)[*~=])\s+)?(==|!=|<|<=|>|>=)\s*$`)

func newCompareOperator(operatorString string) compareOperator {

	if operatorString == "" {
		return compareOperator{compareScopeEqual, ""}
	}

	result := newCompareOperatorRE.FindStringSubmatch(operatorString)
	if result == nil {
		panic(fmt.Errorf("Unable to parse %v into a compareOperator", operatorString))
	}

	scope := compareScopeAsterisk
	switch result[1] {
	case "#*", "{}*":
		scope = compareScopeAsterisk
	case "#~", "{}~":
		scope = compareScopeTilde
	case "#=", "{}=":
		scope = compareScopeEqual
	}

	comparison := result[2]

	return compareOperator{scope, comparison}
}

func compare(left interface{}, operatorString string, right interface{}) (bool, compareOperator) {
	pass := true
	operator := newCompareOperator(operatorString)
	comparator := newComparator(left, operator, right)
	// FIXME Confusing
	switch operator.comparison {
	case "==":
		pass = comparator.IsEqual()
	case "!=":
		pass = !comparator.IsEqual()
	default:
		if comparator.HasOrder() {
			switch operator.comparison {
			case "<":
				pass = comparator.Compare() == -1
			case "<=":
				pass = comparator.Compare() <= 0
			case ">":
				pass = comparator.Compare() == 1
			case ">=":
				pass = comparator.Compare() >= 0
			default:
				panic(fmt.Errorf("Compare operator (%v) is invalid", operator.comparison))
			}
		} else {
			pass = false
		}
	}
	return pass, operator
}

// Compare / Comparator

type compareKind int

const (
	kindInterface compareKind = iota
	kindInteger
	kindUnsignedInteger
	kindFloat
	kindString
	kindBoolean
)

func comparatorValue(value interface{}) (reflect.Value, compareKind) {
	reflectValue := reflect.ValueOf(value)
	kind := kindInterface
	switch value.(type) {
	case int, int8, int16, int32, int64:
		kind = kindInteger
	case uint, uint8, uint16, uint32, uint64:
		kind = kindUnsignedInteger
	case float32, float64:
		kind = kindFloat
	case string:
		kind = kindString
	case bool:
		kind = kindBoolean
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

type aComparator interface {
	Compare() int
	HasOrder() bool
	IsEqual() bool
	CompareScope() compareScope
}

type baseComparator struct {
	hasOrder bool
	operator compareOperator
}

func (self *baseComparator) Compare() int {
	panic(fmt.Errorf("Invalid .Compare()"))
}
func (self *baseComparator) HasOrder() bool {
	return self.hasOrder
}
func (self *baseComparator) CompareScope() compareScope {
	return self.operator.scope
}
func comparatorWithOrder(operator compareOperator) *baseComparator {
	return &baseComparator{true, operator}
}
func comparatorWithoutOrder(operator compareOperator) *baseComparator {
	return &baseComparator{false, operator}
}

type interfaceComparator struct {
	*baseComparator
	left  interface{}
	right interface{}
}

func (self *interfaceComparator) IsEqual() bool {
	if self.CompareScope() != compareScopeEqual {
		return reflect.DeepEqual(self.left, self.right)
	}
	return self.left == self.right
}

type floatComparator struct {
	*baseComparator
	left  float64
	right float64
}

func (self *floatComparator) Compare() int {
	if self.left == self.right {
		return 0
	} else if self.left < self.right {
		return -1
	}
	return 1
}
func (self *floatComparator) IsEqual() bool {
	return self.left == self.right
}

type integerComparator struct {
	*baseComparator
	left  *big.Int
	right *big.Int
}

func (self *integerComparator) Compare() int {
	return self.left.Cmp(self.right)
}
func (self *integerComparator) IsEqual() bool {
	return 0 == self.left.Cmp(self.right)
}

type stringComparator struct {
	*baseComparator
	left  string
	right string
}

func (self *stringComparator) Compare() int {
	if self.left == self.right {
		return 0
	} else if self.left < self.right {
		return -1
	}
	return 1
}
func (self *stringComparator) IsEqual() bool {
	return self.left == self.right
}

type booleanComparator struct {
	*baseComparator
	left  bool
	right bool
}

func (self *booleanComparator) IsEqual() bool {
	return self.left == self.right
}

func newComparator(left interface{}, operator compareOperator, right interface{}) aComparator {
	leftValue, _ := comparatorValue(left)
	rightValue, rightKind := comparatorValue(right)

	// The simplest comparator is comparing interface{} =? interface{}
	targetKind := kindInterface
	// Are left and right of the same kind?
	// (reflect.Value.Kind() is different from compareKind)
	scopeEqual := leftValue.Kind() == rightValue.Kind()
	scopeTilde := false
	scopeAsterisk := false
	if scopeEqual {
		targetKind = rightKind // Since left and right are the same, the targetKind is Integer/Float/String/Boolean
	} else {
		// Examine the prefix of reflect.Value.Kind().String() to see if there is a similarity of 
		// the left value to right value
		lk := leftValue.Kind().String()
		hasPrefix := func(prefix string) bool {
			return strings.HasPrefix(lk, prefix)
		}

		switch right.(type) {
		case float32, float64:
			// Right is float*
			if hasPrefix("float") {
				// Left is also float*
				targetKind = kindFloat
				scopeTilde = true
			} else if hasPrefix("int") || hasPrefix("uint") {
				// Left is a kind of numeric (int* or uint*)
				targetKind = kindFloat
				scopeAsterisk = true
			} else {
				// Otherwise left is a non-numeric
			}
		case uint, uint8, uint16, uint32, uint64:
			// Right is uint*
			if hasPrefix("uint") {
				// Left is also uint*
				targetKind = kindInteger
				scopeTilde = true
			} else if hasPrefix("int") {
				// Left is an int* (a numeric)
				targetKind = kindInteger
				scopeAsterisk = true
			} else if hasPrefix("float") {
				// Left is an float* (a numeric)
				targetKind = kindFloat
				scopeAsterisk = true
			} else {
				// Otherwise left is a non-numeric
			}
		case int, int8, int16, int32, int64:
			// Right is int*
			if hasPrefix("int") {
				// Left is also int*
				targetKind = kindInteger
				scopeTilde = true
			} else if hasPrefix("uint") {
				// Left is a uint* (a numeric)
				targetKind = kindInteger
				scopeAsterisk = true
			} else if hasPrefix("float") {
				// Left is an float* (a numeric)
				targetKind = kindFloat
				scopeAsterisk = true
			} else {
				// Otherwise left is a non-numeric
			}
		default:
			// Right is a non-numeric
			// Can only really compare string to string or boolean to boolean, so
			// we will either have a string/boolean/interfaceComparator
		}
	}

	/*fmt.Println("%v %v %v %v %s %s", operator.scope, same, sibling, family, leftValue, rightValue)*/
	{
		mismatch := false
		switch operator.scope {
		case compareScopeEqual:
			mismatch = !scopeEqual
		case compareScopeTilde:
			mismatch = !scopeEqual && !scopeTilde
		case compareScopeAsterisk:
			mismatch = !scopeEqual && !scopeTilde && !scopeAsterisk
		}
		if mismatch {
			targetKind = kindInterface
		}
	}

	switch targetKind {
	case kindFloat:
		return &floatComparator{
			comparatorWithOrder(operator),
			toFloat(leftValue),
			toFloat(rightValue),
		}
	case kindInteger:
		return &integerComparator{
			comparatorWithOrder(operator),
			toInteger(leftValue),
			toInteger(rightValue),
		}
	case kindString:
		return &stringComparator{
			comparatorWithOrder(operator),
			toString(leftValue),
			toString(rightValue),
		}
	case kindBoolean:
		return &booleanComparator{
			comparatorWithoutOrder(operator),
			toBoolean(leftValue),
			toBoolean(rightValue),
		}
	}

	// As a last resort, we can always compare left (interface{}) to right (interface{})
	return &interfaceComparator{
		comparatorWithoutOrder(operator),
		left,
		right,
	}
}

// failMessage*

func (self *Tester) failMessageForPass(test *test) string {
	return self.FormatMessage(`
        %s:%d: %s 
           Failed test (%s)
                  got: %s
             expected: %s
    `, test.file, test.line, test.Description(), test.kind, stringValue(test.have), stringValue(test.want))
}

func typeKindString(value interface{}) string {
	reflectValue := reflect.ValueOf(value)
	kind := reflectValue.Kind().String()
	result := fmt.Sprintf("%T", value)
	if kind == result {
		if kind == "string" {
			return ""
		}
		return fmt.Sprintf(" (%T)", value)
	}
	return fmt.Sprintf(" (%T=%s)", value, kind)
}

func (self *Tester) failMessageForCompare(test *test) string {
	return self.FormatMessage(`
        %s:%d: %s 
           Failed test (%s)
                  %v%s
                       %s
                  %v%s
    `, test.file, test.line, test.Description(), test.kind, test.have, typeKindString(test.have), test.operator.comparison, test.want, typeKindString(test.want))
}

func (self *Tester) failMessageForEqual(test *test) string {
	return self.failMessageForIs(test)
}

func (self *Tester) failMessageForIs(test *test) string {
	return self.FormatMessage(`
        %s:%d: %v
           Failed test (%s)
                  got: %v%s
             expected: %v%s
    `, test.file, test.line, test.Description(), test.kind, test.have, typeKindString(test.have), test.want, typeKindString(test.want))
}

func (self *Tester) failMessageForMatch(test *test, have, want string, wantMatch bool) string {
	expect := "  like"
	if !wantMatch {
		expect = "unlike"
	}
	return self.FormatMessage(`
        %s:%d: %s 
           Failed test (%s)
                  got: %v%s
               %s: %s
    `, test.file, test.line, test.Description(), test.kind, have, typeKindString(have), expect, want)
}

func (self *Tester) failMessageForLike(test *test, have, want string, wantLike bool) string {
	if !wantLike {
		want = "Anything else"
	}
	return self.FormatMessage(`
        %s:%d: %s 
           Failed test (%s)
                  got: %v%s
             expected: %v%s
    `, test.file, test.line, test.Description(), test.kind, have, typeKindString(have), want, typeKindString(want))
}

// ...

type Tester struct {
	TestingT       *testing.T
	TestEntry      uintptr
	SanityChecking bool
	SelfTesting    bool
	FailIsPassing  bool
}

var ourTester *Tester = nil

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
		ourTester = nil
	} else {
		ourTester = NewTester(arguments[0].(*testing.T))
		ourTester.EnableSanityChecking()
		ourTester.TestEntry = testFunctionEntry()
	}
	return ourTester
}

func UnTerst() {
	ourTester = nil
}

func OurTester() *Tester {
	if ourTester == nil {
		panic("ourTester == nil")
	}
	return ourTester.CheckSanity()
}

func HaveTester() bool {
	return ourTester != nil
}

// Tester

func NewTester(t *testing.T) *Tester {
	return &Tester{t, 0, false, false, false}
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

func (self *Tester) EnableSanityChecking() *Tester {
	self.SanityChecking = true
	return self
}

func (self *Tester) DisableSanityChecking() *Tester {
	self.SanityChecking = false
	return self
}

func (self *Tester) EnableSelfTesting() *Tester {
	self.SelfTesting = true
	return self
}

func (self *Tester) DisableSelfTesting() *Tester {
	self.SelfTesting = false
	return self
}

func (self *Tester) FailIsPass() *Tester {
	self.FailIsPassing = true
	return self
}

func (self *Tester) PassIsPass() *Tester {
	self.FailIsPassing = false
	return self
}

func (self *Tester) CheckSanity() *Tester {
	if self.SanityChecking && self.TestEntry != 0 {
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
	operator  compareOperator

	file       string
	line       int
	functionPC uintptr
	function   string
}

func newTest(kind string, callDepth int, have, want interface{}, arguments []interface{}) *test {
	file, line, functionPC, function, _ := AtFileLineFunction(callDepth + 1)
	operator := newCompareOperator("")
	return &test{kind, have, want, arguments, operator, file, line, functionPC, function}
}

func (self *test) Description() string {
	description := ""
	arguments := self.arguments
	if len(arguments) > 0 {
		description = fmt.Sprintf("%s", arguments[0])
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

func stringValue(value interface{}) string {
	return fmt.Sprintf("%v", value)
}
