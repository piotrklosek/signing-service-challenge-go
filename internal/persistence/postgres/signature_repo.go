package postgres

import (
	"context"
	"database/sql"

	"github.com/piotrklosek/signing-service-challenge-go/internal/domain"
)

type signatureRepo struct {
	db *sql.DB
}

// NewSignatureRepo create interface for signature database
func NewSignatureRepo(db *sql.DB) *signatureRepo {
	return &signatureRepo{db: db}
}

// Create used to create new signature record
func (r *signatureRepo) Create(ctx context.Context, s *domain.SignatureRecord) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO signatures (id, device_id, signed_data, signature, created_at)
         VALUES ($1,$2,$3,$4,$5)`,
		s.ID, s.DeviceID, s.SignedData, s.Signature, s.CreatedAt,
	)
	return err
}

// ListByDevice used to return all signature record by deviceID
func (r *signatureRepo) ListByDevice(ctx context.Context, deviceID string) ([]*domain.SignatureRecord, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, device_id, signed_data, signature, created_at
         FROM signatures WHERE device_id=$1 ORDER BY created_at ASC`, deviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*domain.SignatureRecord
	for rows.Next() {
		var s domain.SignatureRecord
		if err := rows.Scan(&s.ID, &s.DeviceID, &s.SignedData, &s.Signature, &s.CreatedAt); err != nil {
			return nil, err
		}
		records = append(records, &s)
	}
	return records, nil
}
