/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Date: 2020/11/5 3:05 PM
 */

package engine

import (
	"github.com/dengsgo/math-engine/common"
	"github.com/dengsgo/math-engine/entry"
	"github.com/dengsgo/math-engine/source"
	elgamel "github.com/lyszhang/go-homomorphic/elGamel"
	paillier "github.com/roasbeef/go-go-gadget-paillier"
	"math/big"
)

// ExprASTResultHE is a Top level function
// AST traversal
// if an arithmetic runtime error occurs, a panic exception is thrown
func ExprASTResult(expr ExprAST) (res *common.ArithmeticFactor) {
	defer func() { entry.Append(res.String() + "\n") }()
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
				return &common.ArithmeticFactor{
					Factor: common.TypeConst,
					Number: l.Number + r.Number,
				}
			}
			// 如果左侧为常数，右侧为密文
			if l.Factor == common.TypeConst && r.Factor == common.TypePaillier {
				pub := r.Cipher.PublicKey
				plusEandC := paillier.Add(pub, r.Cipher.Data,
					new(big.Int).SetInt64(l.Number).Bytes())
				return &common.ArithmeticFactor{
					Factor: common.TypePaillier,
					Cipher: common.NumberEncrypted{
						Data:      plusEandC,
						PublicKey: pub,
					},
				}
			}
			// 如果左侧为密文，右侧为明文
			if l.Factor == common.TypePaillier && r.Factor == common.TypeConst {
				pub := l.Cipher.PublicKey
				plusCandE := paillier.Add(pub, l.Cipher.Data,
					new(big.Int).SetInt64(r.Number).Bytes())
				return &common.ArithmeticFactor{
					Factor: common.TypePaillier,
					Cipher: common.NumberEncrypted{
						Data:      plusCandE,
						PublicKey: pub,
					},
				}
			}
			// 如果双方均为密文
			if l.Factor == common.TypePaillier && r.Factor == common.TypePaillier {
				lh := l.Cipher.Data
				rh := r.Cipher.Data
				///TODO: 公钥比对
				pub := l.Cipher.PublicKey
				// Add the Cipher integers 15 and 15 together.
				plusEandE := paillier.AddCipher(pub, lh, rh)
				return &common.ArithmeticFactor{
					Factor: common.TypePaillier,
					Number: 0,
					Cipher: common.NumberEncrypted{
						Data:      plusEandE,
						PublicKey: pub,
					},
				}
			}

		case "-":
			// 如果双方都是明文数字
			//TODO：如何检测负数结果的出现
			if l.Factor == common.TypeConst && r.Factor == common.TypeConst {
				return &common.ArithmeticFactor{
					Factor: common.TypeConst,
					Number: l.Number - r.Number,
				}
			}

			// 如果双方都是密文数字
			if l.Factor == common.TypePaillier && r.Factor == common.TypePaillier {
				lh := l.Cipher.Data
				rh := r.Cipher.Data
				///TODO: 公钥比对
				pub := l.Cipher.PublicKey
				// Add the Cipher integers 15 and 15 together.
				subEandE := paillier.SubCipher(pub, lh, rh)
				return &common.ArithmeticFactor{
					Factor: common.TypePaillier,
					Number: 0,
					Cipher: common.NumberEncrypted{
						Data:      subEandE,
						PublicKey: pub,
					},
				}
			}

			if (l.Factor == common.TypePaillier && r.Factor == common.TypeConst) || (l.Factor == common.TypeConst && r.
				Factor == common.TypePaillier) {
				entry.Append("!!!!!!Warning: don't support a cipher and a const number subtract operation!!!!!")
			}

		case "*":
			// 如果双方都是明文数字
			if l.Factor == common.TypeConst && r.Factor == common.TypeConst {
				return &common.ArithmeticFactor{
					Factor: common.TypeConst,
					Number: l.Number * r.Number,
				}
			}

			// 如果左侧为常数，右侧为密文
			if l.Factor == common.TypeConst && r.Factor == common.TypePaillier {
				pub := r.Cipher.PublicKey
				mulEandC := paillier.Mul(pub, r.Cipher.Data,
					new(big.Int).SetInt64(l.Number).Bytes())
				return &common.ArithmeticFactor{
					Factor: common.TypePaillier,
					Cipher: common.NumberEncrypted{
						Data:      mulEandC,
						PublicKey: pub,
					},
				}
			}

			// 如果左侧为密文，右侧为常数
			if l.Factor == common.TypePaillier && r.Factor == common.TypeConst {
				pub := l.Cipher.PublicKey
				mulEandC := paillier.Mul(pub, l.Cipher.Data,
					new(big.Int).SetInt64(r.Number).Bytes())
				return &common.ArithmeticFactor{
					Factor: common.TypePaillier,
					Cipher: common.NumberEncrypted{
						Data:      mulEandC,
						PublicKey: pub,
					},
				}
			}

			// 如果两方均是paillier密文， 则需要重新使用ElGamel重新加密
			if l.Factor == common.TypePaillier && r.Factor == common.TypePaillier {
				///TODO: 校验是否使用相同的密钥加密出的

				// re-encrypt
				lpc := &common.CipherCompression{
					T:            common.TypePaillier,
					PaillierData: l.Cipher.Data,
				}
				lec, _ := source.TransformExternal(lpc)

				rpc := &common.CipherCompression{
					T:            common.TypePaillier,
					PaillierData: r.Cipher.Data,
				}
				rec, _ := source.TransformExternal(rpc)

				// mul
				m := elgamel.Mul(&lec.ElGamalData, &rec.ElGamalData)

				// re-encrypt
				mec := &common.CipherCompression{
					T:           common.TypeElgmel,
					ElGamalData: *m,
				}
				mpc, _ := source.TransformExternal(mec)

				return &common.ArithmeticFactor{
					Factor: common.TypePaillier,
					Cipher: common.NumberEncrypted{
						Data:      mpc.PaillierData,
						PublicKey: l.Cipher.PublicKey,
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
		return &common.ArithmeticFactor{Factor: common.TypeConst, Number: expr.(NumberExprAST).Val}
	case ParameterExprAST:
		data, pub, _ := source.FetchExternalGravity(nil, expr.(ParameterExprAST).Str)
		return &common.ArithmeticFactor{Factor: common.TypePaillier, Cipher: common.NumberEncrypted{
			Data:      data,
			PublicKey: pub,
		}}
	case FunCallerExprAST:
		f := expr.(FunCallerExprAST)
		def := defFunc[f.Name]
		return def.fun(f.Arg...)
	}

	return &common.ArithmeticFactor{Factor: common.TypeConst, Number: 0}
}
