# ondrejsika/counter

## Configuration

- `BACKEND` - Storage enginge for counter, default `redis`. Values can be `redis`, `inmemory`, `postgres`, or `mongodb`.
- `PORT` - port to listen on (default: `80`)
- `REDIS` - Redis host (default: `127.0.0.1`)
- `REDIS_PASSWORD` - Redis password (default: `''` - empty string)
- `POSTGRES_HOST` - Postgres host (default: `127.0.0.1`)
- `POSTGRES_USER` - Postgres user (default: `postgres`)
- `POSTGRES_PASSWORD` - Postgres password (default: `pg`)
- `POSTGRES_DATABASE` - Postgres database (default: `postgres`)
- `POSTGRES_SSLMODE` - Postgres SSL mode (default: `disable`, values can be `disable`, `require`)
- `MONGODB_URI` - MongoDB host (default: `mongodb://127.0.0.1:27017`)
- `SLOW_START` - Time in seconds to wait before start (default: `0`)
- `EXTRA_TEXT` -  Extra text to display (default: `''`)

## Images

- `ondrejsika/counter`
- `ghcr.io/ondrejsika/counter`

## Run Dependencies

Redis

```
docker run --name redis -d -p 6379:6379 redis
```

Get counter value from Redis

```
docker exec redis redis-cli get counter
```

MongoDB

```
docker run --name mongodb -d -p 27017:27017 mongo
```

Get counter value from MongoDB

```
docker exec -it mongodb mongosh counter --eval 'db.counter.find()'
```

Postgres

```
docker run --name postgres -d -p 5432:5432 -e POSTGRES_PASSWORD=pg postgres
```

Get counter value from Postgres

```
docker exec postgres psql -U postgres -c 'SELECT * FROM counters'
```
