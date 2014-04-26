# terst
--
    import "github.com/robertkrimen/terst"

Package terst is a terse (terst = test + terse), easy-to-use testing library for
Go.

terst is compatible with (and works via) the standard testing package:
http://golang.org/pkg/testing

    var is = terst.Is

    func Test(t *testing.T) {
        terst.Terst(t, func() {
            is("abc", "abc")

            is(1, ">", 0)

            var abc []int
            is(abc, nil)
        }
    }

Do not import terst directly, instead use `terst-import` to copy it into your
testing environment:

https://github.com/robertkrimen/terst/tree/master/terst-import

    $ go get github.com/robertkrimen/terst/terst-import

    $ terst-import

## Usage

#### func  Is

```go
func Is(arguments ...interface{}) bool
```
Is compares two values (got & expect) and returns true if the comparison is
true, false otherwise. In addition, if the comparison is false, Is will report
the error in a manner similar to testing.T.Error(...). Is also takes an optional
argument, an operator, that changes how the comparison is made. The following
operators are available:

    ==      # got == expect, This is the default
    !=      # got != expect

    >       # got > expect (float32, uint, uint16, int, int64, ...)
    >=      # got >= expect
    <       # got < expect
    <=      # got <= expect

    =~      # regexp.MustCompile(expect).Match{String}(got)
    !~      # !regexp.MustCompile(expect).Match{String}(got)

A simple comparison:

    Is(2 + 2, 4)

A bit trickier:

    Is(1, ">", 0)
    Is(2 + 2, "!=", 5)
    Is("Nothing happens.", "=~", `ing(\s+)happens\.$`)

Is should only be called under a Terst(t, ...) call. For a standalone version,
use IsErr(...). If no scope is found and the comparison is false, then Is will
panic the error.

#### func  IsErr

```go
func IsErr(arguments ...interface{}) error
```
IsErr compares two values (got & expect) and returns nil if the comparison is
true, an ErrFail if the comparison is false, or an ErrInvalid if the comparison
is invalid. IsErr also takes an optional argument, an operator, that changes how
the comparison is made.

#### func  Terst

```go
func Terst(t *testing.T, arguments ...func())
```
Terst creates a testing scope, where Is can be called and errors will be
reported according to the top-level location of the comparison, and not where
the Is call actually takes place. For example:

    func test() {
        Is(2 + 2, 5) // <--- This failure is reported below.
    }

    Terst(t, func(){

        Is(2, ">", 3) // <--- An error is reported here.

        test() // <--- An error is reported here.

    })

#### type Call

```go
type Call struct {
}
```

Call is a reference to a line immediately under a Terst testing scope.

#### func  Caller

```go
func Caller() *Call
```
Caller will search the stack, looking for a Terst testing scope. If a scope is
found, then Caller returns a Call for logging errors, accessing testing.T, etc.
If no scope is found, Caller returns nil.

#### func (*Call) Error

```go
func (cl *Call) Error(arguments ...interface{})
```
Error is the terst version of `testing.T.Error`

#### func (*Call) Errorf

```go
func (cl *Call) Errorf(format string, arguments ...interface{})
```
Errorf is the terst version of `testing.T.Errorf`

#### func (*Call) Log

```go
func (cl *Call) Log(arguments ...interface{})
```
Log is the terst version of `testing.T.Log`

#### func (*Call) Logf

```go
func (cl *Call) Logf(format string, arguments ...interface{})
```
Logf is the terst version of `testing.T.Logf`

#### func (*Call) Skip

```go
func (cl *Call) Skip(arguments ...interface{})
```
Skip is the terst version of `testing.T.Skip`

#### func (*Call) Skipf

```go
func (cl *Call) Skipf(format string, arguments ...interface{})
```
Skipf is the terst version of `testing.T.Skipf`

#### func (*Call) T

```go
func (cl *Call) T() *testing.T
```
T returns the original testing.T passed to Terst(...)

#### func (*Call) TestFunc

```go
func (cl *Call) TestFunc() *runtime.Func
```
TestFunc returns the *runtime.Func entry for the top-level Test...(t testing.T)
function.

#### type ErrFail

```go
type ErrFail error
```

ErrFail indicates a comparison failure (e.g. 0 > 1).

#### type ErrInvalid

```go
type ErrInvalid error
```

ErrInvalid indicates an invalid comparison (e.g. bool == string).

--
**godocdown** http://github.com/robertkrimen/godocdown
