# ondrejsika/counter

## Configuration

- `BACKEND` - Storage enginge for counter, default `redis`. Values can be `redis`, `inmemory`, `postgres`, or `mongodb`.
- `PORT` - port to listen on (default: `80`)
- `REDIS` - Redis host (default: `127.0.0.1`)
- `POSTGRES_HOST` - Postgres host (default: `127.0.0.1`)
- `POSTGRES_USER` - Postgres user (default: `postgres`)
- `POSTGRES_PASSWORD` - Postgres password (default: `pg`)
- `POSTGRES_DATABASE` - Postgres database (default: `postgres`)
- `MONGODB_URI` - MongoDB host (default: `mongodb://127.0.0.1:27017`)
- `SLOW_START` - Time in seconds to wait before start (default: `0`)
- `EXTRA_TEXT` -  Extra text to display (default: `''`)

## Images

- `ondrejsika/counter`
- `ghcr.io/ondrejsika/counter`
