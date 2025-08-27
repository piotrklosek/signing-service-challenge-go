package domain_test

import (
	"encoding/base64"
	"testing"

	"github.com/piotrklosek/signing-service-challenge-go/internal/domain"
)

func TestPrepareSignedData(t *testing.T) {
	deviceID := "dev123"
	idB64 := base64.StdEncoding.EncodeToString([]byte(deviceID))

	tests := []struct {
		name   string
		device domain.SignatureDevice
		data   string
		want   string
	}{
		{
			name: "first signature uses base64(device.id) even if lastSignature empty",
			device: domain.SignatureDevice{
				ID:               deviceID,
				SignatureCounter: 0,
				LastSignature:    "",
			},
			data: "payload",
			want: "0_payload_" + idB64,
		},
		{
			name: "first signature uses base64(device.id) even if lastSignature preset",
			device: domain.SignatureDevice{
				ID:               deviceID,
				SignatureCounter: 0,
				LastSignature:    "should_be_ignored",
			},
			data: "payload",
			want: "0_payload_" + idB64,
		},
		{
			name: "next signature uses device.LastSignature",
			device: domain.SignatureDevice{
				ID:               deviceID,
				SignatureCounter: 1,
				LastSignature:    "abc==",
			},
			data: "invoice-42",
			want: "1_invoice-42_abc==",
		},
		{
			name: "data may contain underscores",
			device: domain.SignatureDevice{
				ID:               deviceID,
				SignatureCounter: 5,
				LastSignature:    "prevSigB64",
			},
			data: "a_b_c",
			want: "5_a_b_c_prevSigB64",
		},
		{
			name: "empty data allowed",
			device: domain.SignatureDevice{
				ID:               deviceID,
				SignatureCounter: 0,
				LastSignature:    "",
			},
			data: "",
			want: "0__" + idB64,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := domain.PrepareSignedData(&tc.device, tc.data)
			if got != tc.want {
				t.Fatalf("PrepareSignedData() = %q, want %q", got, tc.want)
			}
		})
	}
}
