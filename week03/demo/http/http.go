package main

import (
	"net"
	"net/http"
)

type Handler struct {
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	panic("implement me")
}

func main() {
	l, _ := net.Listen("tcp", ":8080")
	svr := http.Server{Handler: &Handler{}}
	svr.Serve(l)
}
