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

type LeastConnections struct {
}

func (ls LeastConnections) ChooseInstance(services []Service) Service {
	min := int(^uint(0) >> 1)
	var serv Service
	for _, s := range services {
		if s.ActiveConnections < min {
			min = s.ActiveConnections
			serv = s
		}

	}
	return serv
}
