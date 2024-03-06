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
		bs, err := server.New(fmt.Sprintf(":%d", 8001+i))
		if err != nil {
			panic(err)
		}

		go func() {
			log.Fatal(bs.Serve())
		}()

		pool = append(pool, loadbalancer.Service{
			ID:      bs.ID,
			URL:     bs.URL(false),
			Healthy: true,
		})
	}

	lb := loadbalancer.New(loadbalancer.StrategyRoundRobin, pool)
	go lb.CheckHealth()
	log.Fatal(http.ListenAndServe(":8000", lb))
}
