package loadbalancer

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
	return service
}
