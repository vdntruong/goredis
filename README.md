# Goredis

## Pitfalls & Best Practices

### Best practices

- Cache Monitoring and Alerts
- Deploy highly available cache clusters with redundancy
- Cache prewarming

### Read Stragegy

### Write Stragegy

### Pitfalls

- Cache Avalanche aka Cache Stampede
  - Description: A large mumber of requests hits the databases all at once.
  - Causes / Scenarios:
    - When a massive chunk of cached data expires all at once. Or a large number of cache misses happen simultaneously for the same cache key or set of keys.
    - When the cache restarts and its cold and empty (redis cluster crash or restart by any reason)
  - Impact: Cache avalanche can overwhelm the backend systems, leading to increased latency, resource contention, and potential service degradation.
  - Solutions:
    - Cache locking
    - Cache prewarming
    - Randomized cache expiration times
    - Deploy highly available cache clusters
- asdf

## init

```bash
go mod init goredis
```

```bash
docker init
```

## tech

## Workers and [asynq lib](https://github.com/hibiken/asynq)

```bash
  swag init -g /cmd/read/main.go
```