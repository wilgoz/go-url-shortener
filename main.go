package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/wilgoz/go-url-shortener/http"
	"github.com/wilgoz/go-url-shortener/repository/mongodb"
	"github.com/wilgoz/go-url-shortener/repository/redis"
	"github.com/wilgoz/go-url-shortener/shortener"
)

func main() {
	repo := getRepo()
	service := shortener.NewRedirectService(repo)
	handler := http.NewHandler(service)

	errs := make(chan error, 2)
	go func() {
		errs <- handler.Listen()
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	fmt.Printf("Terminated %s", <-errs)
}

// Choose either redis or mongoDB with the option of a redis cache layer
// A URL shortener realistically has more writes than reads, thus adding a cache layer could actually be a bottleneck
func getRepo() shortener.RedirectRepository {
	switch os.Getenv("DB") {
	case "redis":
		redisURL := os.Getenv("REDIS_URL")
		repo, err := redis.NewRedisRepo(redisURL)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	case "mongo":
		mongoURL := os.Getenv("MONGO_URL")
		mongoDB := os.Getenv("MONGO_DB")

		cacheEnabled, _ := strconv.ParseBool(os.Getenv("CACHE_ENABLED"))
		mongoTimeout, _ := strconv.Atoi(os.Getenv("MONGO_TIMEOUT"))

		repo, err := mongodb.NewMongoRepo(mongoURL, mongoDB, mongoTimeout, cacheEnabled)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	}
	return nil
}
