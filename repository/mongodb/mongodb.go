package mongodb

import (
	"context"
	"log"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/wilgoz/go-url-shortener/shortener"
)

type mongoRepository struct {
	client   *mongo.Client
	database string
	cache    shortener.RedirectRepository
	timeout  time.Duration
}

func newMongoClient(mongoURL string, mongoTimeout int) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(mongoTimeout)*time.Second,
	)
	defer cancel()
	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(mongoURL),
	)
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	return client, err
}

// NewMongoRepo sets up a new mongoDB repo
func NewMongoRepo(
	mongoURL, mongoDB string, mongoTimeout int, cache shortener.RedirectRepository,
) (shortener.RedirectRepository, error) {
	repo := &mongoRepository{
		database: mongoDB,
		timeout:  time.Duration(mongoTimeout) * time.Second,
		cache:    cache,
	}
	client, err := newMongoClient(mongoURL, mongoTimeout)
	if err != nil {
		return nil, errors.Wrap(err, "repository.mongo.NewMongoRepo")
	}
	repo.client = client
	return repo, nil
}

func (m *mongoRepository) findInCache(shortened string) (*shortener.Redirect, error) {
	// Attempts to find in cache
	redirect, err := m.cache.Find(shortened)
	if err == nil {
		log.Println("cache hit")
		return redirect, nil
	}
	// Finds from DB and updates the cache on cache misses
	if errors.Cause(err) == shortener.ErrRedirectNotFound {
		log.Println("cache miss")
		redirect, err = m.findInDB(shortened)
		if err == nil {
			if err = m.cache.Store(redirect); err != nil {
				log.Println("failed to update cache")
			}
			return redirect, nil
		}
	}
	return nil, errors.Wrap(err, "repository.mongo.findInCache")
}

func (m *mongoRepository) findInDB(shortened string) (*shortener.Redirect, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		m.timeout,
	)
	defer cancel()
	collection := m.client.Database(m.database).Collection("redirects")
	filter := bson.M{"shortened": shortened}
	redirect := &shortener.Redirect{}
	err := collection.FindOne(ctx, filter).Decode(redirect)
	// Valid redirect found
	if err == nil {
		return redirect, nil
	}
	// No valid redirects
	if err == mongo.ErrNoDocuments {
		return nil, errors.Wrap(shortener.ErrRedirectNotFound, "repository.mongo.findInDB")
	}
	return nil, errors.Wrap(err, "repository.mongo.findInDB")
}

func (m *mongoRepository) Find(shortened string) (*shortener.Redirect, error) {
	if m.cache != nil {
		return m.findInCache(shortened)
	}
	return m.findInDB(shortened)
}

func (m *mongoRepository) Store(model *shortener.Redirect) error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		m.timeout,
	)
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
		return errors.Wrap(err, "repository.mongo.Store")
	}
	return nil
}
