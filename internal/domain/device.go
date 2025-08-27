package domain

import (
	"errors"
	"time"

	"github.com/piotrklosek/signing-service-challenge-go/internal/utils/crypto"
)

// AlgorithmType supported alrorithm
type AlgorithmType string

const (
	AlgorithmRSA AlgorithmType = "RSA"
	AlgorithmECC AlgorithmType = "ECC"
)

// SignatureDevice represent signature device
type SignatureDevice struct {
	ID               string        `json:"id"`
	UserID           string        `json:"user_id"`
	Algorithm        AlgorithmType `json:"algorithm"`
	Label            string        `json:"label,omitempty"`
	PublicKey        []byte        `json:"public_key"`
	PrivateKey       []byte        `json:"-"`
	SignatureCounter uint64        `json:"signature_counter"`
	LastSignature    string        `json:"last_signature"`
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`
}

// GenerateKeys generate pair of keys based on choosen algorithm
func (d *SignatureDevice) GenerateKeys() error {
	switch d.Algorithm {
	case AlgorithmRSA:
		gen := crypto.RSAGenerator{}
		marshaler := crypto.NewRSAMarshaler()

		keyPair, err := gen.Generate()
		if err != nil {
			return err
		}
		pub, priv, err := marshaler.Marshal(*keyPair)
		if err != nil {
			return err
		}
		d.PublicKey = pub
		d.PrivateKey = priv

	case AlgorithmECC:
		gen := crypto.ECCGenerator{}
		marshaler := crypto.NewECCMarshaler()

		keyPair, err := gen.Generate()
		if err != nil {
			return err
		}
		pub, priv, err := marshaler.Encode(*keyPair)
		if err != nil {
			return err
		}
		d.PublicKey = pub
		d.PrivateKey = priv

	default:
		return errors.New("unsupported algorithm: " + string(d.Algorithm))
	}

	return nil
}

// IncrementCounter used to increment signed docs by device and record last signed time
func (d *SignatureDevice) IncrementCounter(newLastSignature string) {
	d.SignatureCounter++
	d.LastSignature = newLastSignature
	d.UpdatedAt = time.Now()
}
