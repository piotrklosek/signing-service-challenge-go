package mongo

import (
	"fmt"

	"github.com/piotrklosek/signing-service-challenge-go/internal/persistence"
	"gopkg.in/mgo.v2"
)

type MongoStore struct {
	session *mgo.Session
}

func NewStore(dbUri string) (*MongoStore, error) {
	session, err := mgo.Dial(dbUri)
	if err != nil {
		return nil, fmt.Errorf("failed to open mongo connection: %w", err)
	}

	if err := session.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping mongo: %w", err)
	}

	return &MongoStore{
		session: session,
	}, nil
}

func (s *MongoStore) Close() error {
	return s.Close()
}

func NewRepositories(dbUri, databaseName string) (
	persistence.DeviceRepository,
	persistence.SignatureRepository,
	persistence.UserRepository,
	error) {

	store, err := NewStore(dbUri)
	if err != nil {
		return nil, nil, nil, err
	}

	device, err := NewDeviceRepo(store.session, databaseName)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error while creating device repository for mongo driver :%v\n", err)
	}

	user, err := NewUserRepo(store.session, databaseName)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error while creating user repository for mongo driver :%v\n", err)
	}

	signature, err := NewSignatureRepo(store.session, databaseName)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error while creating signature repository for mongo driver :%v\n", err)
	}

	return device, signature, user, nil
}
