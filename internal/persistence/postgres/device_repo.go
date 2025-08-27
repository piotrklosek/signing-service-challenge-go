package postgres

import (
	"context"
	"database/sql"

	"github.com/piotrklosek/signing-service-challenge-go/internal/domain"
)

type deviceRepo struct {
	db *sql.DB
}

func NewDeviceRepo(db *sql.DB) *deviceRepo {
	return &deviceRepo{db: db}
}

func (r *deviceRepo) Create(ctx context.Context, d *domain.SignatureDevice) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO signature_devices
        (id, user_id, algorithm, label, public_key, private_key,
         signature_counter, last_signature, created_at, updated_at)
         VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		d.ID, d.UserID, d.Algorithm, d.Label, d.PublicKey, d.PrivateKey,
		d.SignatureCounter, d.LastSignature, d.CreatedAt, d.UpdatedAt,
	)
	return err
}

func (r *deviceRepo) GetByID(ctx context.Context, id string) (*domain.SignatureDevice, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, algorithm, label, public_key, private_key,
                signature_counter, last_signature, created_at, updated_at
         FROM signature_devices WHERE id=$1`, id)

	var d domain.SignatureDevice
	if err := row.Scan(
		&d.ID, &d.UserID, &d.Algorithm, &d.Label, &d.PublicKey, &d.PrivateKey,
		&d.SignatureCounter, &d.LastSignature, &d.CreatedAt, &d.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *deviceRepo) List(ctx context.Context) ([]*domain.SignatureDevice, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, algorithm, label, public_key, private_key,
                signature_counter, last_signature, created_at, updated_at
         FROM signature_devices ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []*domain.SignatureDevice
	for rows.Next() {
		var d domain.SignatureDevice
		if err := rows.Scan(
			&d.ID, &d.UserID, &d.Algorithm, &d.Label, &d.PublicKey, &d.PrivateKey,
			&d.SignatureCounter, &d.LastSignature, &d.CreatedAt, &d.UpdatedAt,
		); err != nil {
			return nil, err
		}
		devices = append(devices, &d)
	}
	return devices, nil
}

func (r *deviceRepo) Update(ctx context.Context, d *domain.SignatureDevice) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE signature_devices
         SET algorithm=$2, label=$3, public_key=$4, private_key=$5,
             signature_counter=$6, last_signature=$7, updated_at=$8
         WHERE id=$1`,
		d.ID, d.Algorithm, d.Label, d.PublicKey, d.PrivateKey,
		d.SignatureCounter, d.LastSignature, d.UpdatedAt,
	)
	return err
}
