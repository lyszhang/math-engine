/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Date: 2020/11/5 3:05 PM
 */

package engine

import (
	"fmt"
	"github.com/dengsgo/math-engine/common"
	"github.com/dengsgo/math-engine/entry"
	"github.com/dengsgo/math-engine/source"
)

// ExprASTResultHE is a Top level function
// AST traversal
// if an arithmetic runtime error occurs, a panic exception is thrown
// TODO: 默认当前支持小数点6位，统一转换，暂不考虑计算过程中出现小数位不一致的情况
func ExprASTResult(expr ExprAST) (res *common.ArithmeticFactor) {
	fmt.Printf("ExprAST: %+v\n", expr)
	entry.Append(fmt.Sprintf("分解运算: ExprAST: %+v\n", expr))
	defer func() { entry.Append(fmt.Sprintf("中间结果 %s\n", res.String())) }()
	var l, r *common.ArithmeticFactor
	switch expr.(type) {
	case BinaryExprAST:
		ast := expr.(BinaryExprAST)
		l = ExprASTResult(ast.Lhs)
		r = ExprASTResult(ast.Rhs)
		switch ast.Op {
		case "+":
			// 如果双方都是明文数字
			if l.Factor == common.TypeConst && r.Factor == common.TypeConst {
				return execAddCC(l, r)
			}
			// 如果左侧为常数，右侧为密文
			if l.Factor == common.TypeConst && r.Factor == common.TypePaillier {
				return execAddCE(l, r)
			}
			// 如果左侧为密文，右侧为明文
			if l.Factor == common.TypePaillier && r.Factor == common.TypeConst {
				return execAddCE(r, l)
			}
			// 如果双方均为密文
			if l.Factor == common.TypePaillier && r.Factor == common.TypePaillier {
				return execAddEE(l, r)
			}

		case "-":
			// 如果双方都是明文数字
			//TODO：如何检测负数结果的出现
			if l.Factor == common.TypeConst && r.Factor == common.TypeConst {
				return execSubCC(l, r)
			}

			// 如果是密文减明文
			if l.Factor == common.TypePaillier && r.Factor == common.TypeConst {
				return execSubEncryptToConst(l, r)
			}

			// 如果是明文减密文
			if l.Factor == common.TypeConst && r.Factor == common.TypePaillier {
				return execSubConstToEncrypt(l, r)
			}

			// 如果双方都是密文数字
			if l.Factor == common.TypePaillier && r.Factor == common.TypePaillier {
				return execSubEE(l, r)
			}

		case "*":
			// 如果双方都是明文数字
			// 由于为了支持小数位，乘法比较特殊，结果需要除以10^6
			if l.Factor == common.TypeConst && r.Factor == common.TypeConst {
				return execMulCC(l, r)
			}

			// 如果左侧为常数，右侧为密文
			if l.Factor == common.TypeConst && r.Factor == common.TypePaillier {
				return execMulCE(l, r)
			}

			// 如果左侧为密文，右侧为常数
			if l.Factor == common.TypePaillier && r.Factor == common.TypeConst {
				return execMulCE(r, l)
			}

			// 如果两方均是paillier密文， 则需要重新使用ElGamel重新加密
			if l.Factor == common.TypePaillier && r.Factor == common.TypePaillier {
				return execMulEE(l, r)
			}

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
		// 将float64转化为整数
		value, offset := Float64ToInterger(expr.(NumberExprAST).Val)
		return &common.ArithmeticFactor{Factor: common.TypeConst, Number: value, Offset: offset}
	case ParameterExprAST:
		f, _ := source.GetExternalGravity(expr.(ParameterExprAST).Str)
		return f
	case FunCallerExprAST:
		f := expr.(FunCallerExprAST)
		def := defFunc[f.Name]
		return def.fun(f.Arg...)
	}

	return &common.ArithmeticFactor{Factor: common.TypeConst, Number: 0}
}
