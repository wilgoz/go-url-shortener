package mongodb

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/wilgoz/go-url-shortener/repository/redis"
	"github.com/wilgoz/go-url-shortener/shortener"
)

type mongoRepository struct {
	client   *mongo.Client
	database string
	cache    shortener.RedirectRepository // redis
	timeout  time.Duration
}

func (m *mongoRepository) findInCache(shortened string) (*shortener.Redirect, error) {
	redirect, err := m.cache.Find(shortened)
	if err != nil {
		// Add to the cache on cache misses
		if errors.Cause(err) == shortener.ErrRedirectNotFound {
			log.Println("cache miss")
			redirect, err = m.findInDB(shortened)
			if err == nil {
				_ = m.cache.Store(redirect)
				return redirect, nil
			}
		}
		return nil, errors.Wrap(err, "repository.Find")
	}
	return redirect, nil
}

func (m *mongoRepository) findInDB(shortened string) (*shortener.Redirect, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	redirect := &shortener.Redirect{}
	collection := m.client.Database(m.database).Collection("redirects")

	filter := bson.M{"shortened": shortened}
	err := collection.FindOne(ctx, filter).Decode(redirect)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Wrap(shortener.ErrRedirectNotFound, "repository.Find")
		}
		return nil, errors.Wrap(err, "repository.Find")
	}
	return redirect, nil
}

func (m *mongoRepository) Find(shortened string) (*shortener.Redirect, error) {
	if m.cache != nil {
		return m.findInCache(shortened)
	}
	return m.findInDB(shortened)
}

func (m *mongoRepository) Store(model *shortener.Redirect) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()
	collection := m.client.Database(m.database).Collection("redirects")
	_, err := collection.InsertOne(
		ctx,
		bson.M{
			"shortened":  model.Shortened,
			"original":   model.Original,
			"created_at": model.CreatedAt,
		},
	)
	if err != nil {
		return errors.Wrap(err, "repository.Store")
	}
	return nil
}

func newMongoClient(mongoURL string, mongoTimeout int) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mongoTimeout)*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	return client, err
}

func initRedisCache() shortener.RedirectRepository {
	redisURL := os.Getenv("REDIS_URL")
	repo, err := redis.NewRedisRepo(redisURL)
	if err != nil {
		log.Fatal(err)
	}
	return repo
}

func NewMongoRepo(mongoURL, mongoDB string, mongoTimeout int, cacheEnabled bool) (shortener.RedirectRepository, error) {
	repo := &mongoRepository{
		database: mongoDB,
		timeout:  time.Duration(mongoTimeout) * time.Second,
	}
	if cacheEnabled {
		repo.cache = initRedisCache()
	}
	client, err := newMongoClient(mongoURL, mongoTimeout)
	if err != nil {
		return nil, errors.Wrap(err, "repository.NewMongoRepo")
	}
	repo.client = client
	return repo, nil
}
