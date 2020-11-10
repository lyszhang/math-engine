/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Date: 2020/11/10 4:49 PM
 */

package common

import (
	"crypto/rand"
	"fmt"
	elgamel "github.com/lyszhang/go-homomorphic/elGamel"
	paillier "github.com/roasbeef/go-go-gadget-paillier"
	"github.com/stretchr/testify/assert"
	"math/big"
	"os"
	"testing"
)

func TestTransform(t *testing.T) {
	// Generate a 128-bit private key.
	var err error
	privKeyPaillier, err := paillier.GenerateKey(rand.Reader, 128)
	if err != nil {
		os.Exit(0)
		return
	}
	privKeyElGamel := elgamel.CreatePrivateKey()

	// Encrypt the number "15".
	m15 := new(big.Int).SetInt64(15)
	c15, err := paillier.Encrypt(&privKeyPaillier.PublicKey, m15.Bytes())
	if err != nil {
		fmt.Println(err)
		return
	}

	c15CipherPaillier := CipherCompression{
		T:            TypePaillier,
		PaillierData: c15,
		ElGamalData:  elgamel.CipherData{},
	}
	c15CipherElgamel, err := c15CipherPaillier.TransformP2E(privKeyPaillier, &privKeyElGamel.PublicKey)
	assert.Nil(t, err)

	denc, _ := elgamel.Decrypt(privKeyElGamel, c15CipherElgamel.ElGamalData.X, c15CipherElgamel.ElGamalData.Y)
	fmt.Printf("\n\n====Decrypted: %d", elgamel.Valint(denc))
}
