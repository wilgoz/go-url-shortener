package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/wilgoz/go-url-shortener/config"
	"github.com/wilgoz/go-url-shortener/http"
	"github.com/wilgoz/go-url-shortener/repository"
	"github.com/wilgoz/go-url-shortener/shortener"
)

func main() {
	if err := config.SetConfigFromFile(); err != nil {
		log.Fatal(err)
	}

	repo := repository.NewRepo()
	service := shortener.NewRedirectService(repo)
	handler := http.NewHandler(service)

	errs := make(chan error, 2)
	go func() {
		errs <- handler.Listen(config.GetConfig().Port)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err := <-errs
	handler.Shutdown()
	fmt.Printf("Terminated: %s", err)
}
