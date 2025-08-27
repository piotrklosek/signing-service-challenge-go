package domain_test

import (
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/piotrklosek/signing-service-challenge-go/internal/domain"
	"github.com/piotrklosek/signing-service-challenge-go/internal/utils/crypto"
)

func TestSignData_Table(t *testing.T) {
	deviceID := "dev123"
	idB64 := base64.StdEncoding.EncodeToString([]byte(deviceID))

	// rsa
	rsaGen := crypto.RSAGenerator{}
	rsaMar := crypto.NewRSAMarshaler()
	rsaKP, err := rsaGen.Generate()
	if err != nil {
		t.Fatalf("rsa generate: %v", err)
	}
	_, rsaPrivPEM, err := rsaMar.Marshal(*rsaKP)
	if err != nil {
		t.Fatalf("rsa marshal: %v", err)
	}

	// ecc
	eccGen := crypto.ECCGenerator{}
	eccMar := crypto.NewECCMarshaler()
	eccKP, err := eccGen.Generate()
	if err != nil {
		t.Fatalf("ecc generate: %v", err)
	}
	_, eccPrivPEM, err := eccMar.Encode(*eccKP)
	if err != nil {
		t.Fatalf("ecc encode: %v", err)
	}

	tests := []struct {
		name      string
		device    domain.SignatureDevice
		data      string
		wantData  string
		wantErr   bool
		verifyRSA bool
	}{
		{
			name: "RSA success - first signature uses base64(device.id)",
			device: domain.SignatureDevice{
				ID:               deviceID,
				Algorithm:        domain.AlgorithmRSA,
				PrivateKey:       rsaPrivPEM,
				SignatureCounter: 0,
				LastSignature:    "",
			},
			data:      "payload",
			wantData:  fmt.Sprintf("0_%s_%s", "payload", idB64),
			wantErr:   false,
			verifyRSA: true,
		},
		{
			name: "RSA success - non-zero counter uses LastSignature",
			device: domain.SignatureDevice{
				ID:               deviceID,
				Algorithm:        domain.AlgorithmRSA,
				PrivateKey:       rsaPrivPEM,
				SignatureCounter: 2,
				LastSignature:    "prevSigB64",
			},
			data:      "data-123",
			wantData:  "2_data-123_prevSigB64",
			wantErr:   false,
			verifyRSA: true,
		},
		{
			name: "ECC success - first signature",
			device: domain.SignatureDevice{
				ID:               deviceID,
				Algorithm:        domain.AlgorithmECC,
				PrivateKey:       eccPrivPEM,
				SignatureCounter: 0,
				LastSignature:    "",
			},
			data:      "invoice",
			wantData:  fmt.Sprintf("0_%s_%s", "invoice", idB64),
			wantErr:   false,
			verifyRSA: false,
		},
		{
			name: "unsupported algorithm",
			device: domain.SignatureDevice{
				ID:         deviceID,
				Algorithm:  "FOO",
				PrivateKey: []byte("irrelevant"),
			},
			data:     "x",
			wantData: "",
			wantErr:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			signedData, signature, err := domain.SignData(&tc.device, tc.data)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if signedData != tc.wantData {
				t.Fatalf("signedData = %q, want %q", signedData, tc.wantData)
			}
			if signature == "" {
				t.Fatalf("signature should not be empty")
			}

			if tc.verifyRSA {
				rawSig, decErr := base64.StdEncoding.DecodeString(signature)
				if decErr != nil {
					t.Fatalf("decode signature: %v", decErr)
				}
				hash := sha256.Sum256([]byte(signedData))
				if verr := rsa.VerifyPKCS1v15(rsaKP.Public, 0, hash[:], rawSig); verr != nil {
					t.Fatalf("rsa verify failed: %v", verr)
				}
			}
		})
	}
}
