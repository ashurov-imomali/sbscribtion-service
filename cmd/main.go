package main

import (
	"fmt"
	"github.com/ashurov-imomali/sbscribtion-service/config"
	"github.com/ashurov-imomali/sbscribtion-service/internal/api"
	"github.com/ashurov-imomali/sbscribtion-service/internal/repository"
	"github.com/ashurov-imomali/sbscribtion-service/internal/server"
	"github.com/ashurov-imomali/sbscribtion-service/internal/usecase"
	"github.com/ashurov-imomali/sbscribtion-service/pkg/db"
	"github.com/ashurov-imomali/sbscribtion-service/pkg/logger"
	"golang.org/x/net/context"
	"os"
	"os/signal"
	"time"
)

func main() {
	log := logger.New()
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Error get configs. Error: %v", err)
		return
	}
	dbSettings := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		cfg.Pg.Host, cfg.Pg.Port, cfg.Pg.Username, cfg.Pg.DbName, cfg.Pg.Password)

	pgConnection, err := db.New(dbSettings)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	r := repository.New(pgConnection)

	service := usecase.New(r, log)

	h := api.New(service)

	srv := server.NewServer(":"+cfg.Srv.Port, h)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Infof("Server starting on port %s", cfg.Srv.Port)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-stop
	log.Infof("%s", "Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Errorf("Graceful shutdown failed: %v", err)
	} else {
		log.Infof("%s", "Server stopped gracefully")
	}
}
