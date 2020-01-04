package repository

import (
	"github.com/pkg/errors"

	"github.com/wilgoz/go-url-shortener/config"
	"github.com/wilgoz/go-url-shortener/repository/mongodb"
	"github.com/wilgoz/go-url-shortener/repository/redis"
	"github.com/wilgoz/go-url-shortener/shortener"
)

func NewRepo() (shortener.RedirectRepository, error) {
	cfg := config.GetConfig()
	switch cfg.Backend {
	case "redis":
		conf := cfg.Redis
		repo, err := redis.NewRedisRepo(conf.RedisURL)
		if err != nil {
			return nil, err
		}
		return repo, nil
	case "mongo":
		conf := cfg.Mongo
		repo, err := mongodb.NewMongoRepo(
			conf.MongoURL, conf.MongoDB, conf.MongoTimeout, conf.CacheURL,
		)
		if err != nil {
			return nil, err
		}
		return repo, nil
	default:
		return nil, errors.New("Unrecognizable backend")
	}
}
