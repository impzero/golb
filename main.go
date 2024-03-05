package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/impzero/golb/backend"
)

type Strategy string

const (
	StrategyRoundRobin Strategy = "round-robin"
)

type Service struct {
	Name    string
	URL     string
	Healthy bool
}

type Balancer interface {
	ChooseInstance(services []Service) Service
}

type RoundRobin struct {
	current int
}

func (rb *RoundRobin) ChooseInstance(services []Service) Service {
	if rb.current == len(services) {
		rb.current = 0
	}

	service := services[rb.current]
	rb.current++
	if !service.Healthy {
		return rb.ChooseInstance(services)
	}
	return service
}

type LoadBalancer struct {
	Strategy  Strategy
	Algorithm Balancer
	Pool      []Service
}

func NewLoadBalancer(strategy Strategy, pool []Service) *LoadBalancer {
	var algorithm Balancer
	switch strategy {
	case StrategyRoundRobin:
		algorithm = &RoundRobin{}
	}

	return &LoadBalancer{
		Strategy:  strategy,
		Pool:      pool,
		Algorithm: algorithm,
	}
}

func (lb *LoadBalancer) CheckHealth() {
	for {
		for i, serv := range lb.Pool {
			i := i
			serv := serv
			go func() {
				resp, err := http.Get(fmt.Sprintf("%s/health", serv.URL))
				if err != nil {
					lb.Pool[i].Healthy = false
					return
				}
				lb.Pool[i].Healthy = resp.StatusCode == http.StatusOK
			}()
		}
		time.Sleep(5 * time.Second)
	}
}

func (lb *LoadBalancer) HealthStatus() map[string]string {
	status := map[string]string{}
	for _, serv := range lb.Pool {
		if serv.Healthy {
			status[serv.Name] = "healthy"
		} else {
			status[serv.Name] = "unhealthy"
		}
	}
	return status
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	service := lb.Algorithm.ChooseInstance(lb.Pool)
	targetURL, err := url.Parse(service.URL)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.ServeHTTP(w, r)
}

func main() {
	pool := []Service{}
	for i := 0; i < 10; i++ {
		bs := backend.NewService(fmt.Sprintf("be-%d", i), fmt.Sprintf(":%d", 8001+i))
		go func() {
			log.Fatal(bs.Serve())
		}()

		pool = append(pool, Service{
			Name:    bs.Name,
			URL:     bs.URL(false),
			Healthy: true,
		})
	}

	lb := NewLoadBalancer(StrategyRoundRobin, pool)

	go lb.CheckHealth()
	log.Fatal(http.ListenAndServe(":8000", lb))
}
