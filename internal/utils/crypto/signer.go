package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

// Signer defines a contract for different types of signing implementations.
type Signer interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
}

// RSASigner represents signer object
type RSASigner struct {
	PrivateKey *rsa.PrivateKey
}

// NewRSASigner creates new signer based on private key
func NewRSASigner(privateKey []byte) (*RSASigner, error) {
	m := NewRSAMarshaler()
	keyPair, err := m.Unmarshal(privateKey)
	if err != nil {
		return nil, err
	}
	return &RSASigner{PrivateKey: keyPair.Private}, nil
}

// Sign provided data using RSA algorithm and return signed data
func (s *RSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	hash := sha256.Sum256(dataToBeSigned)
	return rsa.SignPKCS1v15(rand.Reader, s.PrivateKey, 0, hash[:])
}

// ECCSigner represents signer object
type ECCSigner struct {
	PrivateKey *ecdsa.PrivateKey
}

// NewECCSigner creates new signer based on private key
func NewECCSigner(privateKey []byte) (*ECCSigner, error) {
	m := NewECCMarshaler()
	keyPair, err := m.Decode(privateKey)
	if err != nil {
		return nil, err
	}
	return &ECCSigner{PrivateKey: keyPair.Private}, nil
}

// Sign provided data using ECC algorithm and return signed data
func (s *ECCSigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	hash := sha256.Sum256(dataToBeSigned)

	r, sigS, err := ecdsa.Sign(rand.Reader, s.PrivateKey, hash[:])
	if err != nil {
		return nil, err
	}

	signature := append(r.Bytes(), sigS.Bytes()...)
	return signature, nil
}
