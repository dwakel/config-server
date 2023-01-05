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
	owner, repo, auth := setConfigVars(*logger)

	hb := controllers.NewHeartbeat(logger)
	c := controllers.NewConfig(logger, owner, repo, auth)

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

//todo: increase coverage of method to catch other invalid inputs like whitespaces, etc.
func setConfigVars(l log.Logger) (string, string, string) {
	owner := os.Getenv("OWNER")
	if owner == "" {
		l.Fatalln("Failed to start service, OWNER not provided")
	}
	repo := os.Getenv("REPO_NAME")
	if repo == "" {
		l.Fatalln("Failed to start service, REPO_NAME not provided")
	}
	auth := os.Getenv("AUTH_TOKEN")
	if auth == "" {
		l.Fatalln("Failed to start service, AUTH_TOKEN not Provided")
	}
	return owner, repo, auth
}
