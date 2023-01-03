package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/config-server/controllers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("./config.env")

	logger := log.New(os.Stdout, "config server", log.LstdFlags)
	hb := controllers.NewHeartbeat(logger)
	c := controllers.NewConfig(logger)

	sm := mux.NewRouter()
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/healthcheck", hb.Heartbeat)
	getRouter.HandleFunc("/{path}/{branch}", c.ServeConfig)

	server := &http.Server{
		Addr:         ":8080", //fmt.Sprintf(":%s", 8080),
		Handler:      sm,
		IdleTimeout:  2 * time.Second,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		logger.Println(fmt.Sprintf("Starting Server on port: %s", 8080))
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	logger.Println("Received terminate, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(tc)

}
