package loadbalancer

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type Strategy string

const (
	StrategyRoundRobin Strategy = "round-robin"
)

type Service struct {
	Name              string
	URL               string
	ActiveConnections uint32
	Healthy           bool
}

type LoadBalancer struct {
	Algorithm         Balancer
	Pool              []Service
	ProxyRoundTripper http.RoundTripper
}

func New(strategy Strategy, pool []Service) *LoadBalancer {
	var algorithm Balancer
	switch strategy {
	case StrategyRoundRobin:
		algorithm = &RoundRobin{}
	}

	return &LoadBalancer{
		Pool:              pool,
		Algorithm:         algorithm,
		ProxyRoundTripper: NewTimedRoundTripper(),
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
	var serv Service = lb.Algorithm.ChooseInstance(lb.Pool)
	for !serv.Healthy {
		serv = lb.Algorithm.ChooseInstance(lb.Pool)
	}

	targetURL, err := url.Parse(serv.URL)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.Transport = lb.ProxyRoundTripper
	proxy.ServeHTTP(w, r)
}
