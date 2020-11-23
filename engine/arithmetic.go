/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Date: 2020/11/23 12:03 PM
 */

package engine

import (
	"github.com/dengsgo/math-engine/common"
	"github.com/dengsgo/math-engine/source"
	elgamel "github.com/lyszhang/go-homomorphic/elGamel"
	paillier "github.com/roasbeef/go-go-gadget-paillier"
	"math"
	"math/big"
)

// 明文加法
func execAddCC(l, r *common.ArithmeticFactor) *common.ArithmeticFactor {
	if l.Offset > r.Offset {
		return &common.ArithmeticFactor{
			Factor: common.TypeConst,
			Number: l.Number + r.Number*Pow10(l.Offset-r.Offset),
			Offset: l.Offset,
		}
	}
	return &common.ArithmeticFactor{
		Factor: common.TypeConst,
		Number: r.Number + l.Number*Pow10(r.Offset-l.Offset),
		Offset: r.Offset,
	}
}

// 明文与密文加法
func execAddCE(constFactor, cipherFactor *common.ArithmeticFactor) *common.ArithmeticFactor {
	pub := cipherFactor.Cipher.PublicKey

	if constFactor.Offset > cipherFactor.Offset {
		// 如果常数位的小数位大于密文部分
		mulEandOffset := paillier.Mul(pub, cipherFactor.Cipher.Data,
			new(big.Int).SetInt64(Pow10(constFactor.Offset-cipherFactor.Offset)).Bytes())

		plusEandC := paillier.Add(pub, mulEandOffset,
			new(big.Int).SetInt64(constFactor.Number).Bytes())

		return &common.ArithmeticFactor{
			Factor: common.TypePaillier,
			Cipher: common.NumberEncrypted{
				Data:      plusEandC,
				PublicKey: pub,
			},
			Offset: constFactor.Offset,
		}
	} else {
		// 如果常数位的小数位小于等于密文部分
		plusEandC := paillier.Add(pub, cipherFactor.Cipher.Data,
			new(big.Int).SetInt64(constFactor.Number*Pow10(cipherFactor.Offset-constFactor.Offset)).Bytes())

		return &common.ArithmeticFactor{
			Factor: common.TypePaillier,
			Cipher: common.NumberEncrypted{
				Data:      plusEandC,
				PublicKey: pub,
			},
			Offset: cipherFactor.Offset,
		}
	}
}

//密文加法
func execAddEE(l, r *common.ArithmeticFactor) *common.ArithmeticFactor {
	lh := l.Cipher.Data
	rh := r.Cipher.Data
	///TODO: 公钥比对
	pub := l.Cipher.PublicKey
	if l.Offset > r.Offset {
		mulEandOffset := paillier.Mul(pub, rh, new(big.Int).SetInt64(Pow10(l.Offset-r.Offset)).Bytes())

		// Add the Cipher.
		plusEandE := paillier.AddCipher(pub, lh, mulEandOffset)
		return &common.ArithmeticFactor{
			Factor: common.TypePaillier,
			Cipher: common.NumberEncrypted{
				Data:      plusEandE,
				PublicKey: pub,
			},
			Offset: l.Offset,
		}
	} else {
		mulEandOffset := paillier.Mul(pub, lh, new(big.Int).SetInt64(Pow10(r.Offset-l.Offset)).Bytes())

		// Add the Cipher.
		plusEandE := paillier.AddCipher(pub, rh, mulEandOffset)
		return &common.ArithmeticFactor{
			Factor: common.TypePaillier,
			Cipher: common.NumberEncrypted{
				Data:      plusEandE,
				PublicKey: pub,
			},
			Offset: r.Offset,
		}
	}
}

// 明文减法
func execSubCC(l, r *common.ArithmeticFactor) *common.ArithmeticFactor {
	if l.Offset > r.Offset {
		return &common.ArithmeticFactor{
			Factor: common.TypeConst,
			Number: l.Number - r.Number*Pow10(l.Offset-r.Offset),
			Offset: l.Offset,
		}
	}
	return &common.ArithmeticFactor{
		Factor: common.TypeConst,
		Number: l.Number - r.Number*Pow10(r.Offset-l.Offset),
		Offset: r.Offset,
	}
}

// 密文减明文
func execSubEncryptToConst(l, r *common.ArithmeticFactor) *common.ArithmeticFactor {
	lh := l.Cipher.Data
	pub := l.Cipher.PublicKey

	if l.Offset > r.Offset {
		rh := new(big.Int).SetInt64(r.Number * Pow10(l.Offset-r.Offset)).Bytes()
		subEandC := paillier.SubCipherWithConstant(pub, lh, rh)
		return &common.ArithmeticFactor{
			Factor: common.TypePaillier,
			Cipher: common.NumberEncrypted{
				Data:      subEandC,
				PublicKey: pub,
			},
			Offset: l.Offset,
		}
	} else {
		rh := new(big.Int).SetInt64(r.Number).Bytes()
		mulEandOffset := paillier.Mul(pub, lh, new(big.Int).SetInt64(Pow10(r.Offset-l.Offset)).Bytes())

		subEandC := paillier.SubCipherWithConstant(pub, mulEandOffset, rh)
		return &common.ArithmeticFactor{
			Factor: common.TypePaillier,
			Cipher: common.NumberEncrypted{
				Data:      subEandC,
				PublicKey: pub,
			},
			Offset: r.Offset,
		}
	}
}

// 明文减密文
func execSubConstToEncrypt(l, r *common.ArithmeticFactor) *common.ArithmeticFactor {
	rh := r.Cipher.Data
	pub := r.Cipher.PublicKey
	if l.Offset > r.Offset {
		lh := new(big.Int).SetInt64(l.Number).Bytes()
		mulEandOffset := paillier.Mul(pub, rh, new(big.Int).SetInt64(Pow10(l.Offset-r.Offset)).Bytes())

		subCandE := paillier.SubConstWithCipher(pub, lh, mulEandOffset)
		return &common.ArithmeticFactor{
			Factor: common.TypePaillier,
			Cipher: common.NumberEncrypted{
				Data:      subCandE,
				PublicKey: pub,
			},
			Offset: l.Offset,
		}
	} else {
		lh := new(big.Int).SetInt64(l.Number * Pow10(r.Offset-l.Offset)).Bytes()
		subCandE := paillier.SubConstWithCipher(pub, lh, rh)
		return &common.ArithmeticFactor{
			Factor: common.TypePaillier,
			Cipher: common.NumberEncrypted{
				Data:      subCandE,
				PublicKey: pub,
			},
			Offset: r.Offset,
		}
	}
}

// 密文减密文
func execSubEE(l, r *common.ArithmeticFactor) *common.ArithmeticFactor {
	lh := l.Cipher.Data
	rh := r.Cipher.Data
	///TODO: 公钥比对
	pub := l.Cipher.PublicKey

	if l.Offset > r.Offset {
		mulEandOffset := paillier.Mul(pub, rh, new(big.Int).SetInt64(Pow10(l.Offset-r.Offset)).Bytes())

		//相减
		subEandE := paillier.SubCipher(pub, lh, mulEandOffset)
		return &common.ArithmeticFactor{
			Factor: common.TypePaillier,
			Cipher: common.NumberEncrypted{
				Data:      subEandE,
				PublicKey: pub,
			},
			Offset: l.Offset,
		}
	} else {
		mulEandOffset := paillier.Mul(pub, lh, new(big.Int).SetInt64(Pow10(r.Offset-l.Offset)).Bytes())

		//相减
		subEandE := paillier.SubCipher(pub, mulEandOffset, rh)
		return &common.ArithmeticFactor{
			Factor: common.TypePaillier,
			Cipher: common.NumberEncrypted{
				Data:      subEandE,
				PublicKey: pub,
			},
			Offset: r.Offset,
		}
	}
}

// 明文乘法
func execMulCC(l, r *common.ArithmeticFactor) *common.ArithmeticFactor {
	// 检测小数位是否超过
	if (l.Offset + r.Offset) <= MaxLengthDecimal {
		return &common.ArithmeticFactor{
			Factor: common.TypeConst,
			Number: l.Number * r.Number,
			Offset: l.Offset + r.Offset,
		}
	} else {
		return &common.ArithmeticFactor{
			Factor: common.TypeConst,
			Number: l.Number * r.Number / int64(math.Pow10(int(l.Offset+r.Offset-MaxLengthDecimal))),
			Offset: MaxLengthDecimal,
		}
	}
}

// 明文乘以密文乘法
func execMulCE(constFactor, cipherFactor *common.ArithmeticFactor) *common.ArithmeticFactor {
	pub := cipherFactor.Cipher.PublicKey
	// TODO: 检测小数位长度
	mulEandC := paillier.Mul(pub, cipherFactor.Cipher.Data, new(big.Int).SetInt64(constFactor.Number).Bytes())
	return &common.ArithmeticFactor{
		Factor: common.TypePaillier,
		Cipher: common.NumberEncrypted{
			Data:      mulEandC,
			PublicKey: pub,
		},
		Offset: constFactor.Offset + cipherFactor.Offset,
	}
}

// 明文乘以密文乘法
func execMulEE(l, r *common.ArithmeticFactor) *common.ArithmeticFactor {
	///TODO: 校验是否使用相同的密钥加密出的
	// TODO: 检测小数位长度

	// re-encrypt
	lpc := &common.CipherCompression{
		T:            common.TypePaillier,
		PaillierData: l.Cipher.Data,
		Offset:       l.Offset,
	}
	lec, _ := source.TransformExternal(lpc)

	rpc := &common.CipherCompression{
		T:            common.TypePaillier,
		PaillierData: r.Cipher.Data,
		Offset:       r.Offset,
	}
	rec, _ := source.TransformExternal(rpc)

	// mul
	m := elgamel.Mul(&lec.ElGamalData, &rec.ElGamalData)

	// re-encrypt
	mec := &common.CipherCompression{
		T:           common.TypeElgmel,
		ElGamalData: *m,
		Offset:      lpc.Offset + rpc.Offset,
	}
	mpc, _ := source.TransformExternal(mec)

	return &common.ArithmeticFactor{
		Factor: common.TypePaillier,
		Cipher: common.NumberEncrypted{
			Data:      mpc.PaillierData,
			PublicKey: l.Cipher.PublicKey,
		},
		Offset: mec.Offset,
	}
}
