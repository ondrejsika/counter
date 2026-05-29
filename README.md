# ondrejsika/counter

## Configuration

- `BACKEND` - Storage enginge for counter, default `redis`. Values can be `redis`, `inmemory`, `postgres`, `mongodb`, or `kafka`.
- `PORT` - port to listen on (default: `8000`)
- `API_ONLY` - Disable homepage / index page (`/`) when set to `1`, default is `0`
- `REDIS` - Redis host (default: `127.0.0.1`)
- `REDIS_PASSWORD` - Redis password (default: `''` - empty string)
- `REDIS_SCHEMA` - Fake schema version used for migration testing (default: ``)
- `POSTGRES_HOST` - Postgres host (default: `127.0.0.1`)
- `POSTGRES_USER` - Postgres user (default: `postgres`)
- `POSTGRES_PASSWORD` - Postgres password (default: `pg`)
- `POSTGRES_DATABASE` - Postgres database (default: `postgres`)
- `POSTGRES_SSLMODE` - Postgres SSL mode (default: `disable`, values can be `disable`, `require`)
- `MONGODB_URI` - MongoDB host (default: `mongodb://127.0.0.1:27017`)
- `KAFKA_PEERS` - Comma-separated list of Kafka bootstrap servers (default: `127.0.0.1:9092`)
- `KAFKA_TOPIC` - Kafka topic to store counts (default: `counter`)
- `SLOW_START` - Time in seconds to wait before start (default: `0`)
- `EXTRA_TEXT` -  Extra text to display (default: `''`)
- `EXTRA_TEXT_SUFFIX` -  Extra text suffix to display after the text value (default: `''`)

## Vault Configuration

- `VAULT_ENV_SECRET_PATH` - Vault secret path to load all environment variables from (e.g. `secret/data/counter`), applied before any other config is read
- `VAULT_ADDR` - Vault server address (default: `http://vault.vault:8200`)
- `VAULT_KUBERNETES_AUTH_ROLE` - Vault Kubernetes auth role (required when using `VAULT_ENV_SECRET_PATH`)
- `VAULT_KUBERNETES_AUTH_PATH` - Vault Kubernetes auth mount path (default: `kubernetes`)

## API Endpoints

### Counter Endpoints

- **GET /** - Get counter value with UI rendering (HTML for browsers, plain text for curl)
- **GET /api/counter** - Get and increase counter value in JSON format
- **GET /api/counter-txt** - Get and increase counter value in plain text format
- **GET /api/read-counter** - Read counter value in JSON format, NOT increase the counter
- **GET /api/read-counter-txt** - Read counter value in plain text format, NOT increase the counter

### Health Check Endpoints

- **GET /api/livez** or **GET /livez** - for liveness probe
- **GET /api/readyz** or **GET /readyz** - for readiness probe

### Information Endpoints

- **GET /api/version** or **GET /version** - Get application version
- **GET /api/status** or **GET /status** - Get application status and runtime information

### Monitoring Endpoints

- **GET /metrics** - Prometheus metrics endpoint

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

Kafka

```
docker run --name kafka -d -p 9092:9092 apache/kafka
```

Get counter value from Kafka (number of messages in topic)

```
docker exec kafka /opt/kafka/bin/kafka-get-offsets.sh --bootstrap-server localhost:9092 --topic counter
```
