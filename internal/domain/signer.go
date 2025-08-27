package domain

import (
	"encoding/base64"
	"errors"

	"github.com/piotrklosek/signing-service-challenge-go/internal/utils/crypto"
)

// SignData business logic for signing data
func SignData(device *SignatureDevice, data string) (signedData string, signature string, err error) {
	signedData = PrepareSignedData(device, data)

	var s crypto.Signer
	switch device.Algorithm {
	case AlgorithmRSA:
		s, err = crypto.NewRSASigner(device.PrivateKey)
	case AlgorithmECC:
		s, err = crypto.NewECCSigner(device.PrivateKey)
	default:
		return "", "", errors.New("unsupported algorithm")
	}
	if err != nil {
		return "", "", err
	}

	sigBytes, err := s.Sign([]byte(signedData))
	if err != nil {
		return "", "", err
	}

	// base64-encode signature
	signature = base64.StdEncoding.EncodeToString(sigBytes)

	return signedData, signature, nil
}
