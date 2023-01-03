package controllers

import (
	"fmt"
	"log"
	"net/http"
)

type Heartbeat struct {
	l *log.Logger
}

func NewHeartbeat(l *log.Logger) *Heartbeat {
	return &Heartbeat{l}
}

func (this *Heartbeat) Heartbeat(rw http.ResponseWriter, r *http.Request) {
	this.l.Println("I'm Alive")
	rw.WriteHeader(http.StatusOK)
	fmt.Fprintf(rw, "I'm Alive!!!!!")
}
