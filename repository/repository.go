package repository

import (
	"log"

	"github.com/wilgoz/go-url-shortener/config"
	"github.com/wilgoz/go-url-shortener/repository/mongodb"
	"github.com/wilgoz/go-url-shortener/repository/redis"
	"github.com/wilgoz/go-url-shortener/shortener"
)

func NewRepo() shortener.RedirectRepository {
	cfg := config.GetConfig()
	switch cfg.Backend {
	case "redis":
		conf := cfg.Redis
		repo, err := redis.NewRedisRepo(conf.RedisURL)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	case "mongo":
		conf := cfg.Mongo
		repo, err := mongodb.NewMongoRepo(
			conf.MongoURL, conf.MongoDB, conf.MongoTimeout, conf.CacheEnabled,
		)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	default:
		return nil
	}
}
