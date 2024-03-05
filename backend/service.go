package backend

import (
	"fmt"
	"net/http"
)

type Service struct {
	Name string
}

func NewService(name string) Service {
	return Service{
		Name: name,
	}
}

func (b *Service) handle(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello from %s", b.Name)
}

func (b *Service) health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "%s is healthy", b.Name)
}

func (b *Service) Serve(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", b.handle)
	mux.HandleFunc("/health", b.health)

	s := http.Server{
		Addr:    addr,
		Handler: mux,
	}
	if err := s.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
