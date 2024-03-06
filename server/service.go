package server

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type Server struct {
	Address string
	Name    string
}

func New(name string, addr string) Server {
	return Server{
		Name:    name,
		Address: addr,
	}
}

func (b *Server) handle(w http.ResponseWriter, _ *http.Request) {
	time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello from %s", b.Name)
}

func (b *Server) health(w http.ResponseWriter, _ *http.Request) {
	if b.Name == "be-2" {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (b *Server) URL(tls bool) string {
	if tls {
		return fmt.Sprintf("https://%s", b.Address)
	}
	return fmt.Sprintf("http://%s", b.Address)
}

func (b *Server) Serve() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", b.handle)
	mux.HandleFunc("/health", b.health)

	s := http.Server{
		Addr:    b.Address,
		Handler: mux,
	}
	if err := s.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
