package engine

import (
	"errors"
	"fmt"
	"github.com/dengsgo/math-engine/common"
	"math"
	"strconv"
	"strings"
)

// 默认支持的小数位
const (
	//一般默认为6位小数
	DefaultLengthDecimal = 6
	//最大支持10位小数
	MaxLengthDecimal = 10
)

// Top level function
// Analytical expression and execution
// err is not nil if an error occurs (including arithmetic runtime errors)
func ParseAndExec(s string) (r *common.ArithmeticFactor, err error) {
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

// the integer power of a Number
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

// 浮点数转化为 value/10^offset
func Float64ToInterger(f float64) (value int64, offset int64) {
	fstr := fmt.Sprintf("%v", f)
	if strings.Contains(fstr, ".") {
		length := len(strings.Split(fmt.Sprintf("%v", f), ".")[1])
		if length > 6 {
			length = 6
		}

		return int64(math.Pow10(length) * f), int64(length)
	}
	return int64(f), 0
}

// 浮点数转化为 value/10^offset
// 浮点数默认解析为6位小数，不足位补0
// 0.34 ====> 340000 6
func Float64ToIntergerOffsetDefault(f float64) (value int64, offset int64) {
	return int64(math.Pow10(DefaultLengthDecimal) * f), DefaultLengthDecimal
}

func Pow10(n int64) int64 {
	return int64(math.Pow10(int(n)))
}

// RegFunction is Top level function
// register a new function to use in expressions
func RegFunction(name string, argc int, fun func(...ExprAST) *common.ArithmeticFactor) error {
	if len(name) == 0 {
		return errors.New("RegFunction name is not empty.")
	}

	///TODO: 长度为0时，表示可变参数
	//if argc < 1 {
	//	return errors.New("RegFunction argc is must has one arg at least.")
	//}
	if _, ok := defFunc[name]; ok {
		return errors.New("RegFunction name is already exist.")
	}
	defFunc[name] = defS{argc, fun}
	return nil
}
