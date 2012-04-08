package terst

import (
    "testing"
    "fmt"
    "reflect"
    "runtime"
    "strings"
    "unsafe"
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

func (self *Tester) AtFileLineFunction(callDepth int) (string, int, string) {
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

func (self *Tester) AtFailMessage(callDepth int, kind string, have, want interface{}, arguments ...interface{}) string {
    file, line, _ := self.AtFileLineFunction(callDepth + 1)

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
    output = append( output, moreOutput... )
    *(*[]byte)(unsafe.Pointer(outputValue.UnsafeAddr())) = output;
}

func (self *Tester) AtIs(callDepth int, have, want interface{}, arguments ...interface{}) bool {
	pass := have == want
    if (!pass) {
        self.Log(self.AtFailMessage(callDepth + 1, "Is", have, want, arguments...))
        self.TestingT.Fail()
        return false
    }
    return true
}

func (self *Tester) AtIsNot(callDepth int, have, want interface{}, arguments ...interface{}) bool {
	pass := have != want
    if (!pass) {
        self.Log(self.AtFailMessage(callDepth + 1, "IsNot", have, "Anything else", arguments...))
        self.TestingT.Fail()
        return false
    }
    return true
}


