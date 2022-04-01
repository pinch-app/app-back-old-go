package repo

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"

	"github.com/EQUISEED-WEALTH/pinch/backend/domain/interfaces"
)

type AugmontInMemRepo struct {
	db *redis.Client
}

// NewAugmontInMemRepo returns new AugmontInMemRepo
func NewAugmontInMemRepo(db *redis.Client) interfaces.AugmontInMemRepo {
	return &AugmontInMemRepo{db}
}

// GetToken returns augmont auth token
func (r *AugmontInMemRepo) GetToken() (string, error) {
	token, err := r.db.Get(context.TODO(), "augmont-token").Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		return "", err
	}
	return token, nil
}

//SetToken  Sets auth token with expiry time
func (r *AugmontInMemRepo) SetToken(token string, expireAt time.Time) error {
	log.WithField("expireAt", expireAt).Info("new augmont token added")
	_, err := r.db.Set(
		context.TODO(),
		"augmont-token",
		token,
		time.Until(expireAt),
	).Result()
	if err != nil {
		return err
	}
	return nil
}
