package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/Picus-Security-Golang-Backend-Bootcamp/homework-4-Burak-Atak/handlers"
	"github.com/gorilla/mux"
)

// main is the entrypoint of the application
func main() {

	router := mux.NewRouter()

	router.Use(loggingMiddleware)
	router.Use(authenticationMiddleware)

	router.HandleFunc("/buy/{id}/{count}", handlers.Buy).Methods("PUT")
	router.HandleFunc("/delete/{id}", handlers.Delete).Methods("DELETE")
	router.HandleFunc("/search/{query}", handlers.Search).Methods("GET")
	router.HandleFunc("/create", handlers.Create).Methods("POST")
	router.HandleFunc("/list", handlers.List).Methods("GET")

	srv := &http.Server{
		Addr:         "0.0.0.0:8090",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
	ShutdownServer(srv, time.Second*10)
}

// loggingMiddleware is a simple middleware to log the request
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.URL)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)

	})
}

// authenticationMiddleware is a simple middleware to check if the request has a valid token
func authenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if strings.HasPrefix(r.URL.Path, "/delete") || strings.HasPrefix(r.URL.Path, "/create") {
			if !strings.HasPrefix(token, "Bearer ") {
				http.Error(w, "You are not allowed for this command", http.StatusUnauthorized)
				return
			}
		} else {
			next.ServeHTTP(w, r)
		}

	})

}

// ShutdownServer is a function to shutdown the server
func ShutdownServer(srv *http.Server, timeout time.Duration) {
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
