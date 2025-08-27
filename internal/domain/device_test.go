package domain_test

import (
	"testing"
	"time"

	"github.com/piotrklosek/signing-service-challenge-go/internal/domain"
)

func TestSignatureDevice_GenerateKeys(t *testing.T) {
	tests := []struct {
		name      string
		algorithm domain.AlgorithmType
		wantErr   bool
	}{
		{
			name:      "RSA success",
			algorithm: domain.AlgorithmRSA,
			wantErr:   false,
		},
		{
			name:      "ECC success",
			algorithm: domain.AlgorithmECC,
			wantErr:   false,
		},
		{
			name:      "Unsupported algorithm",
			algorithm: "FOO",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			device := &domain.SignatureDevice{
				ID:        "dev-1",
				UserID:    "user-1",
				Algorithm: tt.algorithm,
			}

			err := device.GenerateKeys()

			if tt.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.wantErr {
				if len(device.PublicKey) == 0 {
					t.Errorf("expected PublicKey to be set")
				}
				if len(device.PrivateKey) == 0 {
					t.Errorf("expected PrivateKey to be set")
				}
			}
		})
	}
}

func TestSignatureDevice_IncrementCounter(t *testing.T) {
	tests := []struct {
		name       string
		initialCnt uint64
		newSig     string
		wantCnt    uint64
		wantSig    string
	}{
		{
			name:       "Increment from zero",
			initialCnt: 0,
			newSig:     "sig-1",
			wantCnt:    1,
			wantSig:    "sig-1",
		},
		{
			name:       "Increment from non-zero",
			initialCnt: 5,
			newSig:     "sig-6",
			wantCnt:    6,
			wantSig:    "sig-6",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			device := &domain.SignatureDevice{
				ID:               "dev-1",
				UserID:           "user-1",
				Algorithm:        domain.AlgorithmRSA,
				SignatureCounter: tt.initialCnt,
				LastSignature:    "old-sig",
				UpdatedAt:        time.Now().Add(-time.Hour),
			}

			before := device.UpdatedAt
			device.IncrementCounter(tt.newSig)

			if device.SignatureCounter != tt.wantCnt {
				t.Errorf("expected counter %d, got %d", tt.wantCnt, device.SignatureCounter)
			}
			if device.LastSignature != tt.wantSig {
				t.Errorf("expected LastSignature %q, got %q", tt.wantSig, device.LastSignature)
			}
			if !device.UpdatedAt.After(before) {
				t.Errorf("expected UpdatedAt to be updated, got %v (before %v)", device.UpdatedAt, before)
			}
		})
	}
}
