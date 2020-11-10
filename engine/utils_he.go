/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Date: 2020/11/5 3:05 PM
 */

package engine

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/dengsgo/math-engine/entry"
	"github.com/dengsgo/math-engine/source"
	paillier "github.com/roasbeef/go-go-gadget-paillier"
	"math/big"
)

type ArithmeticType int

const (
	_ ArithmeticType = iota
	TypePaillier
	TypeElgmel
	TypeConst
	TypeEnd
)

func (t ArithmeticType) String() string {
	switch t {
	case TypePaillier:
		return "TypePaillier"
	case TypeElgmel:
		return "TypeElgmel"
	case TypeConst:
		return "TypeConst"
	case TypeEnd:
		return "TypeEnd"
	default:
		return "Unknown Type"
	}
}

type numberEncrypted struct {
	Data      []byte
	PublicKey *paillier.PublicKey
}

type ArithmeticFactor struct {
	Factor ArithmeticType
	Number int64
	Cipher numberEncrypted
}

func (a *ArithmeticFactor) String() string {
	if a == nil {
		return "nil"
	}
	switch a.Factor {
	case TypeConst:
		return fmt.Sprintf("{Factor: %s, Number: %d}", a.Factor.String(), a.Number)
	case TypePaillier:
		buf, _ := json.Marshal(a.Cipher.PublicKey)
		return fmt.Sprintf("{Factor: %s, Data: %s, Pubkey: %s}", a.Factor, hex.EncodeToString(a.Cipher.Data),
			string(buf))
	default:
		return "unknown ArithmeticFactor"
	}
}

// ExprASTResultHE is a Top level function
// AST traversal
// if an arithmetic runtime error occurs, a panic exception is thrown
func ExprASTResult(expr ExprAST) (res *ArithmeticFactor) {
	defer func() { entry.Append(res.String() + "\n") }()
	var l, r *ArithmeticFactor
	switch expr.(type) {
	case BinaryExprAST:
		ast := expr.(BinaryExprAST)
		l = ExprASTResult(ast.Lhs)
		r = ExprASTResult(ast.Rhs)
		switch ast.Op {
		case "+":
			// 如果双方都是明文数字
			if l.Factor == TypeConst && r.Factor == TypeConst {
				return &ArithmeticFactor{
					Factor: TypeConst,
					Number: l.Number + r.Number,
				}
			}
			// 如果左侧为常数，右侧为密文
			if l.Factor == TypeConst && r.Factor == TypePaillier {
				pub := r.Cipher.PublicKey
				plusEandC := paillier.Add(pub, r.Cipher.Data,
					new(big.Int).SetInt64(l.Number).Bytes())
				return &ArithmeticFactor{
					Factor: TypePaillier,
					Cipher: numberEncrypted{
						Data:      plusEandC,
						PublicKey: pub,
					},
				}
			}
			// 如果左侧为密文，右侧为明文
			if l.Factor == TypePaillier && r.Factor == TypeConst {
				pub := l.Cipher.PublicKey
				plusCandE := paillier.Add(pub, l.Cipher.Data,
					new(big.Int).SetInt64(r.Number).Bytes())
				return &ArithmeticFactor{
					Factor: TypePaillier,
					Cipher: numberEncrypted{
						Data:      plusCandE,
						PublicKey: pub,
					},
				}
			}
			// 如果双方均为密文
			if l.Factor == TypePaillier && r.Factor == TypePaillier {
				lh := l.Cipher.Data
				rh := r.Cipher.Data
				///TODO: 公钥比对
				pub := l.Cipher.PublicKey
				// Add the Cipher integers 15 and 15 together.
				plusEandE := paillier.AddCipher(pub, lh, rh)
				return &ArithmeticFactor{
					Factor: TypePaillier,
					Number: 0,
					Cipher: numberEncrypted{
						Data:      plusEandE,
						PublicKey: pub,
					},
				}
			}

		case "-":
			// 如果双方都是明文数字
			//TODO：如何检测负数结果的出现
			if l.Factor == TypeConst && r.Factor == TypeConst {
				return &ArithmeticFactor{
					Factor: TypeConst,
					Number: l.Number - r.Number,
				}
			}

			// 如果双方都是密文数字
			if l.Factor == TypePaillier && r.Factor == TypePaillier {
				lh := l.Cipher.Data
				rh := r.Cipher.Data
				///TODO: 公钥比对
				pub := l.Cipher.PublicKey
				// Add the Cipher integers 15 and 15 together.
				subEandE := paillier.SubCipher(pub, lh, rh)
				return &ArithmeticFactor{
					Factor: TypePaillier,
					Number: 0,
					Cipher: numberEncrypted{
						Data:      subEandE,
						PublicKey: pub,
					},
				}
			}

		case "*":
			// 如果双方都是明文数字
			if l.Factor == TypeConst && r.Factor == TypeConst {
				return &ArithmeticFactor{
					Factor: TypeConst,
					Number: l.Number * r.Number,
				}
			}

			// 如果左侧为常数，右侧为密文
			if l.Factor == TypeConst && r.Factor == TypePaillier {
				pub := r.Cipher.PublicKey
				mulEandC := paillier.Mul(pub, r.Cipher.Data,
					new(big.Int).SetInt64(l.Number).Bytes())
				return &ArithmeticFactor{
					Factor: TypePaillier,
					Cipher: numberEncrypted{
						Data:      mulEandC,
						PublicKey: pub,
					},
				}
			}

			// 如果左侧为密文，右侧为常数
			if l.Factor == TypePaillier && r.Factor == TypeConst {
				pub := l.Cipher.PublicKey
				mulEandC := paillier.Mul(pub, l.Cipher.Data,
					new(big.Int).SetInt64(r.Number).Bytes())
				return &ArithmeticFactor{
					Factor: TypePaillier,
					Cipher: numberEncrypted{
						Data:      mulEandC,
						PublicKey: pub,
					},
				}
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
		return &ArithmeticFactor{Factor: TypeConst, Number: expr.(NumberExprAST).Val}
	case ParameterExprAST:
		data, pub, _ := source.FetchExternalGravity(nil, expr.(ParameterExprAST).Str)
		return &ArithmeticFactor{Factor: TypePaillier, Cipher: numberEncrypted{
			Data:      data,
			PublicKey: pub,
		}}
	case FunCallerExprAST:
		f := expr.(FunCallerExprAST)
		def := defFunc[f.Name]
		return def.fun(f.Arg...)
	}

	return &ArithmeticFactor{Factor: TypeConst, Number: 0}
}
