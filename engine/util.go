package engine

import (
	"errors"
	"strconv"
	"strings"
)

// Top level function
// Analytical expression and execution
// err is not nil if an error occurs (including arithmetic runtime errors)
func ParseAndExec(s string) (r *ArithmeticFactor, err error) {
	toks, err := Parse(s)
	if err != nil {
		return nil, err
	}
	ast := NewAST(toks, s)
	if ast.Err != nil {
		return nil, ast.Err
	}
	ar := ast.ParseExpression()
	if ast.Err != nil {
		return nil, ast.Err
	}
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	return ExprASTResult(ar), err
}

func ErrPos(s string, pos int) string {
	r := strings.Repeat("-", len(s)) + "\n"
	s += "\n"
	for i := 0; i < pos; i++ {
		s += " "
	}
	s += "^\n"
	return r + s + r
}

// the integer power of a number
func Pow(x float64, n int) float64 {
	if x == 0 {
		return 0
	}
	r := calPow(x, n)
	if n < 0 {
		r = 1 / r
	}
	return r
}

func calPow(x float64, n int) float64 {
	if n == 0 {
		return 1
	}
	if n == 1 {
		return x
	}
	r := calPow(x, n>>1) // move right 1 byte
	r *= r
	if n&1 == 1 {
		r *= x
	}
	return r
}

// Float64ToStr float64 -> string
func Float64ToStr(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// RegFunction is Top level function
// register a new function to use in expressions
func RegFunction(name string, argc int, fun func(...ExprAST) *ArithmeticFactor) error {
	if len(name) == 0 {
		return errors.New("RegFunction name is not empty.")
	}
	if argc < 1 {
		return errors.New("RegFunction argc is must has one arg at least.")
	}
	if _, ok := defFunc[name]; ok {
		return errors.New("RegFunction name is already exist.")
	}
	defFunc[name] = defS{argc, fun}
	return nil
}

