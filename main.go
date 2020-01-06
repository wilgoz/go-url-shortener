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
	if err := config.InitConfigFromFile(); err != nil {
		log.Fatal(err)
	}
	repo, err := repository.NewRepo()
	if err != nil {
		log.Fatal(err)
	}
	service := shortener.NewRedirectService(repo)
	handler := http.NewHandler(service)

	quitChan := make(chan error, 2)
	go func() {
		quitChan <- handler.Listen(config.GetConfig().Port)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		quitChan <- fmt.Errorf("%s", <-c)
	}()

	quit := <-quitChan
	handler.Shutdown()
	fmt.Printf("Terminated: %s\n", quit)
}
