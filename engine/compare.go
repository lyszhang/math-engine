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
	"github.com/dengsgo/math-engine/common"
	"github.com/dengsgo/math-engine/source"
	paillier "github.com/roasbeef/go-go-gadget-paillier"
	"math/big"
)

//密文与明文比较
func compareEandC(expr0, expr1 *common.ArithmeticFactor) (int64, error) {
	//生成随机数m，n
	m, _ := rand.Int(rand.Reader, big.NewInt(2^32))
	n, _ := rand.Int(rand.Reader, big.NewInt(2^32))

	if expr0.Factor == common.TypePaillier && expr1.Factor == common.TypeConst {
		// E(A)^m+n
		mulEandC := paillier.Mul(expr0.Cipher.PublicKey, expr0.Cipher.Data, m.Bytes())
		addEandC := paillier.Add(expr0.Cipher.PublicKey, mulEandC, n.Bytes())

		// E(mA+n)
		mulCandC := big.NewInt(0).Mul(m, big.NewInt(expr1.Number))
		addCandC := big.NewInt(0).Add(mulCandC, n)
		addCandCEncrypted, _ := paillier.Encrypt(expr0.Cipher.PublicKey, addCandC.Bytes())

		// 发送给数据提供方
		return source.CompareResult(addEandC, addCandCEncrypted)
	} else if expr0.Factor == common.TypeConst && expr1.Factor == common.TypePaillier {
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
func compareEandE(expr0, expr1 *common.ArithmeticFactor) (int64, error) {
	//生成随机数m，n
	m, _ := rand.Int(rand.Reader, big.NewInt(2^32))
	n, _ := rand.Int(rand.Reader, big.NewInt(2^32))

	if expr0.Factor == common.TypePaillier && expr1.Factor == common.TypePaillier {
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
func compareCandC(expr0, expr1 *common.ArithmeticFactor) int64 {
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
func Compare(expr ...ExprAST) *common.ArithmeticFactor {
	expr0 := ExprASTResult(expr[0])
	expr1 := ExprASTResult(expr[1])

	// 一个参数为密文，一个参数为明文
	if (expr0.Factor == common.TypePaillier && expr1.Factor == common.TypeConst) || (expr0.Factor == common.TypeConst && expr1.
		Factor == common.TypePaillier) {
		resInt, err := compareEandC(expr0, expr1)
		if err != nil {
			fmt.Printf("compare failed, err:%s\n", err.Error())
		}
		return &common.ArithmeticFactor{Factor: common.TypeConst, Number: resInt + 1}
	}

	// 均为密文
	if expr0.Factor == common.TypePaillier && expr1.Factor == common.TypePaillier {
		resInt, err := compareEandE(expr0, expr1)
		if err != nil {
			fmt.Printf("compare failed, err:%s\n", err.Error())
		}
		return &common.ArithmeticFactor{Factor: common.TypeConst, Number: resInt + 1}
	}

	// 均为明文
	if expr0.Factor == common.TypeConst && expr1.Factor == common.TypeConst {
		resInt := compareCandC(expr0, expr1)
		return &common.ArithmeticFactor{Factor: common.TypeConst, Number: resInt + 1}
	}

	fmt.Printf("compare unspport parameter type")
	return nil
}

// 参数中有几个为1
// CountOne(A, B...)
// 参数必须为明文
func CountOne(exprs ...ExprAST) *common.ArithmeticFactor {
	var sum int64
	for _, expr := range exprs {
		exprAST := ExprASTResult(expr)

		if exprAST.Factor != common.TypeConst {
			fmt.Println("countOne must with const parameters")
		}
		if exprAST.Number == 1 {
			sum += 1
		}
	}
	return &common.ArithmeticFactor{
		Factor: common.TypeConst,
		Number: sum,
	}
}

// 确定参数是否落在区间内
// Ratio(0, 1, A, 3)
// [0,1] 区间
// A 参数
// 如果落在区间内，返回系数3，否则返回0
func Ratio(expr ...ExprAST) *common.ArithmeticFactor {
	exprStart := expr[0]
	exprEnd := expr[1]
	exprPara := expr[2]
	exprCoff := ExprASTResult(expr[3])

	// 落在区间内
	if Compare(exprStart, exprPara).Number <= int64(1) && Compare(exprEnd, exprPara).Number >= int64(1) {
		return exprCoff
	}

	return &common.ArithmeticFactor{
		Factor: common.TypeConst,
		Number: 0,
	}
}

// 求取区间内最大的数字
// Max(A，1，B，C ...)
// 参数可变长度，任意明文或密文
// 返回最大值
func Max(exprs ...ExprAST) *common.ArithmeticFactor {
	var exprTmp ExprAST
	for index, expr := range exprs {
		if index == 0 {
			exprTmp = expr
			continue
		}

		//当小于新的值
		if Compare(exprTmp, expr).Number == int64(0) {
			exprTmp = expr
		}
	}
	return ExprASTResult(exprTmp)
}

// 求取区间内最小的数字
// Min(A，1，B，C ...)
// 参数可变长度，任意明文或密文
// 返回最小值
func Min(exprs ...ExprAST) *common.ArithmeticFactor {
	var exprTmp ExprAST
	for index, expr := range exprs {
		if index == 0 {
			exprTmp = expr
			continue
		}

		//当大于新的值
		if Compare(exprTmp, expr).Number == int64(2) {
			exprTmp = expr
		}
	}
	return ExprASTResult(exprTmp)
}
