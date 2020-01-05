package redis

import (
	"strconv"

	"github.com/go-redis/redis/v7"
	"github.com/pkg/errors"

	"github.com/wilgoz/go-url-shortener/shortener"
)

type redisRepository struct {
	client *redis.Client
}

func newRedisClient(redisURL string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opts)
	_, err = client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

// NewRedisRepo initializes redis repo
func NewRedisRepo(redisURL string) (shortener.RedirectRepository, error) {
	repo := &redisRepository{}
	client, err := newRedisClient(redisURL)
	if err != nil {
		return nil, errors.Wrap(
			err, "repository.redis.NewRedisRepo",
		)
	}
	repo.client = client
	return repo, nil
}

func (r *redisRepository) Find(shortened string) (*shortener.Redirect, error) {
	data, err := r.client.HGetAll(shortened).Result()
	if err != nil {
		return nil, errors.Wrap(
			err, "repository.redis.Find",
		)
	}
	if len(data) == 0 {
		return nil, errors.Wrap(
			shortener.ErrRedirectNotFound, "repository.redis.Find",
		)
	}
	createdAt, err := strconv.ParseInt(data["created_at"], 10, 64)
	if err != nil {
		return nil, errors.Wrap(
			err, "repository.redis.Find",
		)
	}
	redirect := &shortener.Redirect{
		Original:  data["original"],
		Shortened: data["shortened"],
		CreatedAt: createdAt,
	}
	return redirect, nil
}

func (r *redisRepository) Store(model *shortener.Redirect) error {
	data := map[string]interface{}{
		"original":   model.Original,
		"shortened":  model.Shortened,
		"created_at": model.CreatedAt,
	}
	_, err := r.client.HMSet(model.Shortened, data).Result()
	if err != nil {
		return errors.Wrap(
			err, "repository.redis.Store",
		)
	}
	return nil
}
