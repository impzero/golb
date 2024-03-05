package main

import (
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/impzero/golb/backend"
)

type Strategy string

const (
	StrategyRoundRobin Strategy = "round-robin"
)

type Service struct {
	Hostname string
	Healthy  bool
}

type LoadBalancer struct {
	Pool     map[string][]Service
	Strategy Strategy
}

func NewLoadBalancer(strategy Strategy) *LoadBalancer {
	return &LoadBalancer{
		Strategy: strategy,
	}
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	proxy := httputil.NewSingleHostReverseProxy(r.URL)
	proxy.ServeHTTP(w, r)
}

func main() {
	be1 := backend.NewService("be1")
	be2 := backend.NewService("be2")

	go func() {
		log.Fatal(be1.Serve(":8080"))
	}()

	go func() {
		log.Fatal(be2.Serve(":8081"))
	}()

	lb := NewLoadBalancer(StrategyRoundRobin)

	log.Fatal(http.ListenAndServe(":8000", lb))
}
