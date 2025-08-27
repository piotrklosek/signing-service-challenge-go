package mongo

import (
	"context"

	"github.com/piotrklosek/signing-service-challenge-go/internal/domain"
	"gopkg.in/mgo.v2"
)

const userCollectioName = "user"

type userRepo struct {
	sess         *mgo.Session
	databaseName string
}

func NewUserRepo(sess *mgo.Session, databaseName string) (*userRepo, error) {
	c := sess.DB(databaseName).C(userCollectioName)
	key := "user_uuid"
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
	return &userRepo{
		sess:         sess,
		databaseName: databaseName,
	}, nil
}

func (r *userRepo) Create(ctx context.Context, u *domain.User) error {
	panic("implement me")
}

func (r *userRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	panic("implement me")
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	panic("implement me")
}

func (r *userRepo) List(ctx context.Context) ([]*domain.User, error) {
	panic("implement me")
}
