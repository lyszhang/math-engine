/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Date: 2020/11/10 4:16 PM
 */

package common

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	elgamel "github.com/lyszhang/go-homomorphic/elGamel"
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

type NumberEncrypted struct {
	Data      []byte
	PublicKey *paillier.PublicKey
}

type ArithmeticFactor struct {
	Factor ArithmeticType  `json:"Factor"`
	Number int64           `json:"Number, omitempty"`
	Cipher NumberEncrypted `json:"Cipher, omitempty"`
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

type CipherCompression struct {
	T            ArithmeticType     `json:"T"`
	PaillierData []byte             `json:"paillier,omitempty"`
	ElGamalData  elgamel.CipherData `json:"elgamel,omitempty"`
}

func (c *CipherCompression) TransformP2E(ppriv *paillier.PrivateKey, epub *elgamel.PublicKey) (*CipherCompression, error) {
	if c.T == TypePaillier {
		// Decrypt.
		decryptedAddition, err := paillier.Decrypt(ppriv, c.PaillierData)
		if err != nil {
			return nil, err
		}
		valueStr := new(big.Int).SetBytes(decryptedAddition).String()

		// Re encrypt by ElGamel
		e1, _ := elgamel.Encrypt(rand.Reader, epub, valueStr)
		return &CipherCompression{
			T:            TypeElgmel,
			PaillierData: nil,
			ElGamalData:  *e1,
		}, nil
	}
	return nil, errors.New("c not a paillier cipher data")
}

func (c *CipherCompression) TransformE2P(epriv *elgamel.PrivateKey, ppub *paillier.PublicKey) (*CipherCompression, error) {
	if c.T == TypeElgmel {
		// Decrypt.
		decryptBytes, _ := elgamel.Decrypt(epriv, c.ElGamalData.X, c.ElGamalData.Y)
		decryptValue := elgamel.Valint(decryptBytes)
		value := new(big.Int).SetInt64(int64(decryptValue))

		// Re encrypt by Paillier
		encrypted, err := paillier.Encrypt(ppub, value.Bytes())
		if err != nil {
			return nil, err
		}
		return &CipherCompression{
			T:            TypeElgmel,
			PaillierData: encrypted,
		}, nil
	}
	return nil, errors.New("c not a elgmel cipher data")
}
