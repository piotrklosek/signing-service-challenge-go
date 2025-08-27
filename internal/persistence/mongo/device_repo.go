package mongo

import (
	"context"

	"github.com/piotrklosek/signing-service-challenge-go/internal/domain"
	"gopkg.in/mgo.v2"
)

const deviceCollectioName = "device"

type deviceRepo struct {
	sess         *mgo.Session
	databaseName string
}

func NewDeviceRepo(sess *mgo.Session, databaseName string) (*deviceRepo, error) {
	c := sess.DB(databaseName).C(userCollectioName)
	key := "device_uuid"
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
	return &deviceRepo{
		sess:         sess,
		databaseName: databaseName,
	}, nil
}

func (r *deviceRepo) Create(ctx context.Context, d *domain.SignatureDevice) error {
	panic("implement me")
}
func (r *deviceRepo) GetByID(ctx context.Context, id string) (*domain.SignatureDevice, error) {
	panic("implement me")
}
func (r *deviceRepo) List(ctx context.Context) ([]*domain.SignatureDevice, error) {
	panic("implement me")
}
func (r *deviceRepo) Update(ctx context.Context, d *domain.SignatureDevice) error {
	panic("implement me")
}
