package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mraiyuu/M-Pesa/internal/handlers"
	"github.com/mraiyuu/M-Pesa/internal/services"
	// mpesaexpress "github.com/mraiyuu/M-Pesa/internal/mpesa_express"
)

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/"))

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	mpesaExpressService := services.NewService()
	mpesExpressHandler := handlers.NewHandler(mpesaExpressService)
	r.Post("/initiateMpesaExpress", mpesExpressHandler.InitiateMpesaExpress)

	return r
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:              app.config.addr,
		Handler:           h,
		WriteTimeout:      time.Second * 30,
		ReadHeaderTimeout: time.Second * 10,
		IdleTimeout:       time.Minute,
	}

	log.Printf("server has started at addr %s", app.config.addr)
	return srv.ListenAndServe()
}

type application struct {
	config config
	//logger
	//db driver
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	dsn string
}
