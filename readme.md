<h1 align="center">⚖️ golb</h1>

<h4 align="center">
a simple Go Load-Balancer written from scratch
</h4>

### Round-Robin Demo
![](https://github.com/impzero/golb/assets/35530157/3ecfd141-7ec8-4ae2-bffc-ab305503880b)


#### Features

- [x] Round-Robin Load-Balancing
- [x] Health-Check
- [ ] Weighted Round-Robin Load-Balancing
- [ ] Random Load-Balancing
- [ ] Weighted Random Load-Balancing
- [ ] Dynamic Weighted Round-Robin Load-Balancing
- [ ] Least-Connection Load-Balancing
- [ ] Peak Exponentially Weighted Load-Balancing
- [ ] Keep-Alive Connection
- [ ] Logging
- [ ] Metrics
- [ ] Circuit-Breaker

### Why

I found [this cool website](https://codingchallenges.fyi/challenges/challenge-load-balancer/) and decided to give it a try.

Also previously inspired by Sam's awesome load-balancing visualization blog article [here](https://samwho.dev/load-balancing/)

## Run locally

```bash
git clone http://github.com/impzero/golb
```

**or**

```bash
git clone git@github.com:impzero/golb
```

**then**

```bash
cd golb
go run main.go
```
