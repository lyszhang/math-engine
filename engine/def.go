package engine

import "github.com/dengsgo/math-engine/common"

const (
	RadianMode = iota
	AngleMode
)

type defS struct {
	argc int
	fun  func(expr ...ExprAST) *common.ArithmeticFactor
}

// enum "RadianMode", "AngleMode"
var TrigonometricMode = RadianMode

var defConst = map[string]float64{}

var defFunc map[string]defS

func init() {
	defFunc = map[string]defS{
		//"sin": {1, defSin},
		//"cos": {1, defCos},
		//"tan": {1, defTan},
		//"cot": {1, defCot},
		//"sec": {1, defSec},
		//"csc": {1, defCsc},
		//
		//"abs":   {1, defAbs},
		//"ceil":  {1, defCeil},
		//"floor": {1, defFloor},
		//"round": {1, defRound},
		//"sqrt":  {1, defSqrt},
		//"cbrt":  {1, defCbrt},
		//
		//"noerr": {1, defNoerr},
		//
		//"max": {2, defMax},
		//"min": {2, defMin},
	}
}

// noerr(1/0) = 0
// noerr(2.5/(1-1)) = 0
func defNoerr(expr ...ExprAST) (r *common.ArithmeticFactor) {
	defer func() {
		if e := recover(); e != nil {
			r = nil
		}
	}()
	return ExprASTResult(expr[0])
}
