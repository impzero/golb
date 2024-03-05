package backend

import (
	"fmt"
	"net/http"
)

type Service struct {
	Address string
	Name    string
}

func NewService(name string, addr string) Service {
	return Service{
		Name:    name,
		Address: addr,
	}
}

func (b *Service) handle(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello from %s", b.Name)
}

func (b *Service) health(w http.ResponseWriter, _ *http.Request) {
	if b.Name == "be-8" {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (b *Service) URL(tls bool) string {
	if tls {
		return fmt.Sprintf("https://%s", b.Address)
	}
	return fmt.Sprintf("http://%s", b.Address)
}

func (b *Service) Serve() error {
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
