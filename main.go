package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/impzero/golb/loadbalancer"
	"github.com/impzero/golb/server"
)

func main() {
	pool := []loadbalancer.Service{}
	for i := 0; i < 10; i++ {
		bs := server.New(fmt.Sprintf("be-%d", i), fmt.Sprintf(":%d", 8001+i))
		go func() {
			log.Fatal(bs.Serve())
		}()

		pool = append(pool, loadbalancer.Service{
			Name:    bs.Name,
			URL:     bs.URL(false),
			Healthy: true,
		})
	}

	lb := loadbalancer.New(loadbalancer.StrategyRoundRobin, pool)
	go lb.CheckHealth()
	log.Fatal(http.ListenAndServe(":8000", lb))
}
