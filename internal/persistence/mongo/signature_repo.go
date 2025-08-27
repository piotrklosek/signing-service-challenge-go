package mongo

import (
	"context"

	"github.com/piotrklosek/signing-service-challenge-go/internal/domain"
	"gopkg.in/mgo.v2"
)

const signatureCollectioName = "signature"

type signatureRepo struct {
	sess         *mgo.Session
	databaseName string
}

func NewSignatureRepo(sess *mgo.Session, databaseName string) (*signatureRepo, error) {
	c := sess.DB(databaseName).C(userCollectioName)
	key := "signature_uuid"
	index := mgo.Index{
		Key:        []string{key},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	if err := c.EnsureIndex(index); err != nil {
		return nil, err
	}
	return &signatureRepo{
		sess:         sess,
		databaseName: databaseName,
	}, nil
}

func (r *signatureRepo) Create(ctx context.Context, s *domain.SignatureRecord) error {
	panic("implement me")
}

func (r *signatureRepo) ListByDevice(ctx context.Context, deviceID string) ([]*domain.SignatureRecord, error) {
	panic("implement me")
}
