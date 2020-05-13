package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/mistikel/api-websocket-go/utils"
)

type Server interface {
	Serve()
	EnableGracefulShutdown()
}

type Service struct {
	Router *mux.Router
}

func New() Server {
	h := NewHandler()
	return &Service{
		Router: h.CreateRouter(),
	}
}

func (s *Service) Serve() {
	ctx := context.Background()
	s.EnableGracefulShutdown()
	middlewares := alice.New(utils.LoggingHandler)
	utils.InfoContext(ctx, "service: Warpin-pubsub is running")
	utils.InfoContext(ctx, "service: Rest Server mounted at [::]:8080")
	err := http.ListenAndServe(":8080", middlewares.Then(s.Router))
	if err != nil {
		log.Fatal(err.Error())
	}
}

func (s *Service) EnableGracefulShutdown() {
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go s.handleShutdown(signalChannel)
}

func (s *Service) handleShutdown(ch chan os.Signal) {
	<-ch
	defer os.Exit(0)
	ctx := context.Background()
	duration := time.Duration(1 * time.Second)
	utils.InfoContext(ctx, "service: Signal termination received. Waiting %v seconds to shutdown.", duration.Seconds())
	IsShuttingDown = true
	time.Sleep(duration)
	utils.InfoContext(ctx, "service: Cleaning up resources...\n")
	utils.InfoContext(ctx, "service: Bye\n")
}
