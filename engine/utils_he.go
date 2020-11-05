/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Date: 2020/11/5 3:05 PM
 */

package engine

import (
	paillier "github.com/roasbeef/go-go-gadget-paillier"
)

type ArithmeticType int
const (
	_ ArithmeticType = iota
	TypePaillier
	TypeElgmel
	TypeConst
)

type numberEncrypted struct {
	Data []byte
	PublicKey *paillier.PublicKey
}

type ArithmeticFactor struct {
	factor ArithmeticType
	number int64
	Cipher numberEncrypted
}

// ExprASTResultHE is a Top level function
// AST traversal
// if an arithmetic runtime error occurs, a panic exception is thrown
func ExprASTResult(expr ExprAST) *ArithmeticFactor {
	var l, r *ArithmeticFactor
	switch expr.(type) {
	case BinaryExprAST:
		ast := expr.(BinaryExprAST)
		l = ExprASTResult(ast.Lhs)
		r = ExprASTResult(ast.Rhs)
		switch ast.Op {
		case "+":
			lh := l.Cipher.Data
			rh := r.Cipher.Data
			pub := l.Cipher.PublicKey

			// Add the Cipher integers 15 and 15 together.
			plusM16M20 := paillier.AddCipher(pub, lh, rh)
			return &ArithmeticFactor{
				factor:    TypePaillier,
				number:    0,
				Cipher: numberEncrypted{
					Data: plusM16M20,
					PublicKey: pub,
				},
			}

		///TODO: 逐个完善各类型运算符操作
		//case "-":
		//	lh, _ := new(big.Float).SetString(Float64ToStr(l))
		//	rh, _ := new(big.Float).SetString(Float64ToStr(r))
		//	f, _ := new(big.Float).Sub(lh, rh).Float64()
		//	return f
		//case "*":
		//	f, _ := new(big.Float).Mul(new(big.Float).SetFloat64(l), new(big.Float).SetFloat64(r)).Float64()
		//	return f
		//case "/":
		//	if r == 0 {
		//		panic(errors.New(
		//			fmt.Sprintf("violation of arithmetic specification: a division by zero in ExprASTResult: [%g/%g]",
		//				l,
		//				r)))
		//	}
		//	f, _ := new(big.Float).Quo(new(big.Float).SetFloat64(l), new(big.Float).SetFloat64(r)).Float64()
		//	return f
		//case "%":
		//	return float64(int(l) % int(r))
		default:

		}
	case NumberExprAST:
		return &ArithmeticFactor{factor: TypeConst, number: expr.(NumberExprAST).Val}
	case ParameterExprAST:
		data, pub, _ := FetchExternalGravity(nil)
		return &ArithmeticFactor{factor: TypePaillier, Cipher: numberEncrypted{
			Data:      data,
			PublicKey: pub,
		}}
	case FunCallerExprAST:
		f := expr.(FunCallerExprAST)
		def := defFunc[f.Name]
		return def.fun(f.Arg...)
	}

	return &ArithmeticFactor{factor: TypeConst, number: 0}
}

