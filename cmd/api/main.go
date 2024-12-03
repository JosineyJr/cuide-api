package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"cuide/api/router"
	"cuide/config"
	"cuide/util/logger"
	"cuide/util/validator"

	_ "github.com/lib/pq"
)

const fmtDBString = "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable"

//	@title			CUIDE API
//	@version		1.0
//	@description	This is a sample RESTful API with a CRUD

// @host		localhost:8080
// @basePath	/v1
func main() {
	c := config.New()
	l := logger.New(c.Server.Debug)
	v := validator.New()

	dbString := fmt.Sprintf(
		fmtDBString,
		c.DB.Host,
		c.DB.Username,
		c.DB.Password,
		c.DB.DBName,
		c.DB.Port,
	)
	db, err := sql.Open("postgres", dbString)
	if err != nil {
		l.Fatal().Err(err).Msg("DB connection start failure")
		return
	}

	r := router.New(l, v, db)

	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", c.Server.Port),
		Handler:      r,
		ReadTimeout:  c.Server.TimeoutRead,
		WriteTimeout: c.Server.TimeoutWrite,
		IdleTimeout:  c.Server.TimeoutIdle,
	}

	closed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		l.Info().Msgf("Shutting down server %v", s.Addr)

		ctx, cancel := context.WithTimeout(context.Background(), c.Server.TimeoutIdle)
		defer cancel()

		if err := s.Shutdown(ctx); err != nil {
			l.Error().Err(err).Msg("Server shutdown failure")
		}

		if err = db.Close(); err != nil {
			l.Error().Err(err).Msg("DB connection closing failure")
		}

		close(closed)
	}()

	l.Info().Msgf("Starting server %v", s.Addr)
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		l.Fatal().Err(err).Msg("Server startup failure")
	}

	<-closed
	l.Info().Msgf("Server shutdown successfully")
}
