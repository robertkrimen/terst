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

// Is

func Is(have, want interface{}, arguments ...interface{}) bool {
    return OurTester().AtIs(1, have, want, arguments...)
}

func (self *Tester) Is(have, want interface{}, arguments ...interface{}) bool {
    return self.AtIs(1, have, want, arguments...)
}

// IsNot

func IsNot(have, want interface{}, arguments ...interface{}) bool {
    return OurTester().AtIsNot(1, have, want, arguments...)
}

func (self *Tester) IsNot(have, want interface{}, arguments ...interface{}) bool {
    return self.AtIsNot(1, have, want, arguments...)
}

// Like

func Like(have, want interface{}, arguments ...interface{}) bool {
    return OurTester().AtLike(1, have, want, arguments...)
}

func (self *Tester) Like(have, want interface{}, arguments ...interface{}) bool {
    return self.AtLike(1, have, want, arguments...)
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
    callDepth int
    have interface{}
    want interface{}
    arguments []interface{}
}

func newTest(callDepth int, have, want interface{}, arguments ...interface{}) *aTest {
    return &aTest{callDepth, have, want, arguments}
}

func (self *aTest) Sink() *aTest {
    test := *self
    return &test
}

func (self *aTest) AtFileLineFunction() (string, int, string) {
    return AtFileLineFunction(self.callDepth + 1)
}

func (self *aTest) Description() string {
    description := ""
    if len(self.arguments) > 0 {
        description = fmt.Sprintf("%s", self.arguments...)
    }
    return description
}

func AtFileLineFunction(callDepth int) (string, int, string) {
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
        file = "?"
        line = 1
    }
    return file, line, function
}

func (self *Tester) AtFailMessage(test *aTest, kind string) string {
    file, line, _ := test.AtFileLineFunction()
    description := test.Description()

    message := fmt.Sprintf(`
        %s:%d: %s 
           Failed test (%s)
                  got: %s
             expected: %s
    `, file, line, description, kind, test.have, test.want)
    message = strings.TrimLeft(message, "\n")
    message = strings.TrimRight(message, " \n")
    return message + "\n\n"
}


func (self *Tester) AtFailMessageForMatch(callDepth int, kind string, have, want string, arguments ...interface{}) string {
    file, line, _ := AtFileLineFunction(callDepth + 1)

    description := ""
    if len(arguments) > 0 {
        description = fmt.Sprintf("%s", arguments...)
    }

    message := fmt.Sprintf(`
        %s:%d: %s 
           Failed test (%s)
                  got: %s
             expected: %s
    `, file, line, description, kind, have, want)
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

func (self *Tester) AtIs(callDepth int, have, want interface{}, arguments ...interface{}) bool {
    test := newTest(callDepth, have, want, arguments)
	pass := have == want
    if (!pass) {
        self.Log(self.AtFailMessage(test, "Is"))
        self.TestingT.Fail()
        return false
    }
    return true
}

func (self *Tester) AtIsNot(callDepth int, have, want interface{}, arguments ...interface{}) bool {
    test := newTest(callDepth, have, want, arguments)
	pass := have != want
    if (!pass) {
        self.Log(self.AtFailMessage(test, "IsNot"))
        self.TestingT.Fail()
        return false
    }
    return true
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

func (self *Tester) AtLike(callDepth int, have, want interface{}, arguments ...interface{}) bool {
    test := newTest(callDepth, have, want, arguments)
    switch want0 := want.(type) {
    case string:
        haveString := ToString(have)
        pass, error := regexp.Match(want0, []byte(haveString))
        if error != nil {
            panic("regexp.Match(" + want0 + ", ...): " + error.Error())
        }
        if (!pass) {
            self.Log(self.AtFailMessage(test.Sink(), "Like"))
            self.TestingT.Fail()
            return false
        }
    }
    return true
}
