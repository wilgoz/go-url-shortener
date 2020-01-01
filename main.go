package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/wilgoz/go-url-shortener/http"
	"github.com/wilgoz/go-url-shortener/repository/mongodb"
	"github.com/wilgoz/go-url-shortener/repository/redis"
	"github.com/wilgoz/go-url-shortener/shortener"
)

func main() {
	repo := getRepo()
	service := shortener.NewRedirectService(repo)
	handler := http.NewHandler(service)

	e := echo.New()
	e.Use(
		middleware.Logger(),
		middleware.Recover(),
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{
				echo.GET,
				echo.POST,
			},
		}),
	)

	e.GET("/:code", handler.Get)
	e.POST("/", handler.Post)

	errs := make(chan error, 2)
	go func() {
		errs <- e.Start(getPort())
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	fmt.Printf("Terminated %s", <-errs)
}

func getPort() string {
	port := "8080"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	return fmt.Sprintf(":%s", port)
}

// Choose either redis or mongoDB with the option of a redis cache layer
// A URL shortener realistically has more writes than reads, thus adding a cache layer could actually result in a bottleneck
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
