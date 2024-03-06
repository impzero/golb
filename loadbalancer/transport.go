package loadbalancer

import (
	"net/http"
	"time"
)

type TimedRoundTripper struct {
	next        http.RoundTripper
	ReqDuration time.Duration
}

func NewTimedRoundTripper() *TimedRoundTripper {
	return &TimedRoundTripper{
		next: http.DefaultTransport,
	}
}

func (c *TimedRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()
	resp, err := c.next.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	c.ReqDuration = time.Since(start)
	return resp, nil
}
