package terst

import (
    "testing"
    "fmt"
    "reflect"
    "runtime"
    "strings"
    "unsafe"
    "strconv"
    "math"
    "regexp"
)

// Pass

func Pass(have bool, arguments ...interface{}) bool {
    return OurTester().AtPass(1, have, arguments...)
}

func (self *Tester) Pass(have bool, arguments ...interface{}) bool {
    return self.AtPass(1, have, arguments...)
}

func (self *Tester) AtPass(callDepth int, have bool, arguments ...interface{}) bool {
    return self.atPassOrFail(true, 1, have, arguments...)
}

// Fail

func Fail(have bool, arguments ...interface{}) bool {
    return OurTester().AtFail(1, have, arguments...)
}

func (self *Tester) Fail(have bool, arguments ...interface{}) bool {
    return self.AtFail(1, have, arguments...)
}

func (self *Tester) AtFail(callDepth int, have bool, arguments ...interface{}) bool {
    return self.atPassOrFail(false, 1, have, arguments...)
}

func (self *Tester) atPassOrFail(want bool, callDepth int, have bool, arguments ...interface{}) bool {
    kind := "Pass"
    if want == false {
        kind = "Fail"
    }
    test := newTest(kind, callDepth + 1, have, want, arguments)
    pass := have == want
    if (!pass) {
        self.Log(self.failMessageForPass(test))
        self.TestingT.Fail()
        return false
    }
    return true
}

// Is

func Is(have, want interface{}, arguments ...interface{}) bool {
    return OurTester().AtIs(1, have, want, arguments...)
}

func (self *Tester) Is(have, want interface{}, arguments ...interface{}) bool {
    return self.AtIs(1, have, want, arguments...)
}

func (self *Tester) AtIs(callDepth int, have, want interface{}, arguments ...interface{}) bool {
    test := newTest("Is", callDepth + 1, have, want, arguments)
	pass := have == want
    if (!pass) {
        self.Log(self.failMessageForIs(test))
        self.TestingT.Fail()
        return false
    }
    return true
}

// IsNot

func IsNot(have, want interface{}, arguments ...interface{}) bool {
    return OurTester().AtIsNot(1, have, want, arguments...)
}

func (self *Tester) IsNot(have, want interface{}, arguments ...interface{}) bool {
    return self.AtIsNot(1, have, want, arguments...)
}

func (self *Tester) AtIsNot(callDepth int, have, want interface{}, arguments ...interface{}) bool {
    test := newTest("IsNot", callDepth + 1, have, want, arguments)
	pass := have != want
    if (!pass) {
        self.Log(self.failMessageForIs(test))
        self.TestingT.Fail()
        return false
    }
    return true
}

// Like

func Like(have, want interface{}, arguments ...interface{}) bool {
    return OurTester().AtLike(1, have, want, arguments...)
}

func (self *Tester) Like(have, want interface{}, arguments ...interface{}) bool {
    return self.AtLike(1, have, want, arguments...)
}

func (self *Tester) AtLike(callDepth int, have, want interface{}, arguments ...interface{}) bool {
    return self.atLikeOrUnlike(true, 1, have, want, arguments...)
}

// Unlike

func Unlike(have, want interface{}, arguments ...interface{}) bool {
    return OurTester().AtUnlike(1, have, want, arguments...)
}

func (self *Tester) Unlike(have, want interface{}, arguments ...interface{}) bool {
    return self.AtUnlike(1, have, want, arguments...)
}

func (self *Tester) AtUnlike(callDepth int, have, want interface{}, arguments ...interface{}) bool {
    return self.atLikeOrUnlike(false, 1, have, want, arguments...)
}

func (self *Tester) atLikeOrUnlike(wantLike bool, callDepth int, have, want interface{}, arguments ...interface{}) bool {
    test := newTest("Like", callDepth + 1, have, want, arguments)
    switch want0 := want.(type) {
    case string:
        haveString := ToString(have)
        pass, error := regexp.Match(want0, []byte(haveString))
        if !wantLike {
            pass = !pass
        }
        if error != nil {
            panic("regexp.Match(" + want0 + ", ...): " + error.Error())
        }
        if (!pass) {
            self.Log(self.failMessageForMatch(test, haveString, want0, wantLike))
            self.TestingT.Fail()
            return false
        }
    }
    return true
}

// failMessage*

func (self *Tester) failMessageForPass(test *aTest) string {
    return self.FormatMessage(`
        %s:%d: %s 
           Failed test (%s)
                  got: %s
             expected: %s
    `, test.file, test.line, test.Description(), test.kind, ToString(test.have), ToString(test.want))
}

func (self *Tester) failMessageForIs(test *aTest) string {
    return self.FormatMessage(`
        %s:%d: %s 
           Failed test (%s)
                  got: %s
             expected: %s
    `, test.file, test.line, test.Description(), test.kind, test.have, test.want)
}

func (self *Tester) failMessageForMatch(test *aTest, have, want string, wantMatch bool) string {
    expect := "mismatched"
    if !wantMatch {
        expect = "   matched"
    }
    return self.FormatMessage(`
        %s:%d: %s 
           Failed test (%s)
                  got: %s
           %s: %s
    `, test.file, test.line, test.Description(), test.kind, test.have, expect, test.want)
}

// ...

type Tester struct {
    TestingT *testing.T
}

var _ourTester *Tester = nil

func WithTester(t *testing.T) *Tester {
    _ourTester = NewTester(t)
    return _ourTester
}

func OurTester() *Tester {
    if _ourTester == nil {
        panic("_ourTester == nil")
    }
    return _ourTester
}

func HaveTester() bool {
    return _ourTester != nil
}

func NewTester(t *testing.T) *Tester {
    return &Tester{t}
}

type aTest struct {
    kind string
    have interface{}
    want interface{}
    arguments []interface{}

    file string
    line int
    functionPC uintptr
    function string
}

func newTest(kind string, callDepth int, have, want interface{}, arguments ...interface{}) *aTest {
    file, line, functionPC, function, _ := AtFileLineFunction(callDepth + 1)
    return &aTest{kind, have, want, arguments, file, line, functionPC, function}
}

func (self *aTest) Description() string {
    description := ""
    if len(self.arguments) > 0 {
        description = fmt.Sprintf("%s", self.arguments...)
    }
    return description
}

func AtFileLineFunction(callDepth int) (string, int, uintptr, string, bool) {
    functionPC, file, line, good := runtime.Caller(callDepth + 1)
    function := runtime.FuncForPC(functionPC).Name()
    if (good) {
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
    *(*[]byte)(unsafe.Pointer(outputValue.UnsafeAddr())) = output;
}

func ToString(value interface{}) string {
    switch value0 := value.(type) {
        case bool:
            return strconv.FormatBool(value0)
        case int, int8, int16, int32, int64:
        case uint, uint8, uint16, uint32, uint64:
            return value0.(string)
        case float32:
            return strconv.FormatFloat(float64(value0), 'd', -1, 32)
        case float64:
            if math.IsNaN(value0) {
                return "NaN"
            } else if math.IsInf(value0, 0) {
                return "Infinity"
            }
            return strconv.FormatFloat(float64(value0), 'd', -1, 64)
        case string:
            return value0
    }
    return reflect.ValueOf(value).String()
}

