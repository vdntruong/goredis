# Simple job queues aka workers based on asynq lib

## Pages

- Job monitoring: http://localhost:8080/monitoring

- Redis Insight: http://localhost:8001

## Setup 

### Start a single Redis cluster go along with RedisInsight

```bash
  docker run -d \
    -e REDIS_ARGS="--requirepass mypassword" \
    --name redis-stack \
    -p 6379:6379 -p 8001:8001 \
    redis/redis-stack:latest
```

### Redis Stack packaging

There are two distinct Redis Stack packages to choose from:
[refer](https://redis.io/docs/about/about-stack/)

- Redis Stack Server: This package contains Redis OSS and module extensions only. It does not 
contain RedisInsight, the developer desktop application. This package is best for production deployment and is intended to be a drop-in replacement (for example, if you're already deploying Redis OSS as a cache). You can still download RedisInsight separately.

- Redis Stack: This package contains everything a developer needs in a single bundle. This 
includes Redis Stack Server (Redis OSS and module extensions) along with the RedisInsight desktop application (or part of the docker container). If you want to create an application locally and explore how it interacts with Redis, this is the package for you.

