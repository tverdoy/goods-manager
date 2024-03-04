# Goods manager
this is a test task. The task itself is:
## Task
### 1. Migrate database
Create a database with the following tables:
- project
- goods

Create index on goods name and insert one row to project

### 2. Create API
Create GRUD for goods. Also create handler `Reprioritize`.
Goods has priority field. When add new row, it should be the highest priority.
When call `Reprioritize`, then set a new priority for the given row and update other rows.

### 3. Create cache
Using `redis` was created cache system. TTL for cache is 60 seconds.

### 4. Create logger
Create async logger, that send/receive log event from `NATS` and save to `ClickHouse`.

# How run
For configuration see `.env`.

## Docker
```shell
docker-compose up
```

# Documentation
API has documentation at address http://localhost:8080/swagger/index.html

Or see [swagger.yaml](internal/docs/swagger.yaml)