/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Date: 2020/11/9 10:25 AM
 */

package engine

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/dengsgo/math-engine/source"
	paillier "github.com/roasbeef/go-go-gadget-paillier"
	"math/big"
)

//密文与明文比较
func compareEandC(expr0, expr1 *ArithmeticFactor) (int64, error) {
	//生成随机数m，n
	m, _ := rand.Int(rand.Reader, big.NewInt(2^32))
	n, _ := rand.Int(rand.Reader, big.NewInt(2^32))

	if expr0.Factor == TypePaillier && expr1.Factor == TypeConst {
		// E(A)^m+n
		mulEandC := paillier.Mul(expr0.Cipher.PublicKey, expr0.Cipher.Data, m.Bytes())
		addEandC := paillier.Add(expr0.Cipher.PublicKey, mulEandC, n.Bytes())

		// E(mA+n)
		mulCandC := big.NewInt(0).Mul(m, big.NewInt(expr1.Number))
		addCandC := big.NewInt(0).Add(mulCandC, n)
		addCandCEncrypted, _ := paillier.Encrypt(expr0.Cipher.PublicKey, addCandC.Bytes())

		// 发送给数据提供方
		return source.CompareResult(addEandC, addCandCEncrypted)
	} else if expr0.Factor == TypeConst && expr1.Factor == TypePaillier {
		// E(mA+n)
		mulCandC := big.NewInt(0).Mul(m, big.NewInt(expr0.Number))
		addCandC := big.NewInt(0).Add(mulCandC, n)
		addCandCEncrypted, _ := paillier.Encrypt(expr1.Cipher.PublicKey, addCandC.Bytes())

		// E(A)^m+n
		mulEandC := paillier.Mul(expr1.Cipher.PublicKey, expr1.Cipher.Data, m.Bytes())
		addEandC := paillier.Add(expr1.Cipher.PublicKey, mulEandC, n.Bytes())

		// 发送给数据提供方
		return source.CompareResult(addCandCEncrypted, addEandC)
	} else {
		return 0, errors.New("check the compare parameter type")
	}
}

//密文比较
func compareEandE(expr0, expr1 *ArithmeticFactor) (int64, error) {
	//生成随机数m，n
	m, _ := rand.Int(rand.Reader, big.NewInt(2^32))
	n, _ := rand.Int(rand.Reader, big.NewInt(2^32))

	if expr0.Factor == TypePaillier && expr1.Factor == TypePaillier {
		// E(A)^m+n
		mulEandC0 := paillier.Mul(expr0.Cipher.PublicKey, expr0.Cipher.Data, m.Bytes())
		addEandC0 := paillier.Add(expr0.Cipher.PublicKey, mulEandC0, n.Bytes())

		// E(B)^m+n
		mulEandC1 := paillier.Mul(expr0.Cipher.PublicKey, expr0.Cipher.Data, m.Bytes())
		addEandC1 := paillier.Add(expr0.Cipher.PublicKey, mulEandC1, n.Bytes())

		// 发送给数据提供方
		return source.CompareResult(addEandC0, addEandC1)
	} else {
		return 0, errors.New("check the compare parameter type")
	}
}

// 明文比较
func compareCandC(expr0, expr1 *ArithmeticFactor) int64 {
	var resInt int64
	if expr0.Number < expr1.Number {
		resInt = -1
	} else if expr0.Number == expr1.Number {
		resInt = 0
	} else {
		resInt = 1
	}
	return resInt
}

// 比较明文和密文
// compare(E(A), B)
// 0: A < B
// 1: A = B
// 2: A > B
func Compare(expr ...ExprAST) *ArithmeticFactor {
	expr0 := ExprASTResult(expr[0])
	expr1 := ExprASTResult(expr[1])

	// 一个参数为密文，一个参数为明文
	if (expr0.Factor == TypePaillier && expr1.Factor == TypeConst) || (expr0.Factor == TypeConst && expr1.
		Factor == TypePaillier) {
		resInt, err := compareEandC(expr0, expr1)
		if err != nil {
			fmt.Printf("compare failed, err:%s\n", err.Error())
		}
		return &ArithmeticFactor{Factor: TypeConst, Number: resInt + 1}
	}

	// 均为密文
	if expr0.Factor == TypePaillier && expr1.Factor == TypePaillier {
		resInt, err := compareEandE(expr0, expr1)
		if err != nil {
			fmt.Printf("compare failed, err:%s\n", err.Error())
		}
		return &ArithmeticFactor{Factor: TypeConst, Number: resInt + 1}
	}

	// 均为明文
	if expr0.Factor == TypeConst && expr1.Factor == TypeConst {
		resInt := compareCandC(expr0, expr1)
		return &ArithmeticFactor{Factor: TypeConst, Number: resInt + 1}
	}

	fmt.Printf("compare unspport parameter type")
	return nil
}
