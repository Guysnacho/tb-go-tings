# Infra

Local Kafka broker for development, defined two ways:

- `docker-compose.yml` — single-node Kafka (KRaft mode, no ZooKeeper) using the stock `apache/kafka:latest` image, wired up with the env vars needed to run broker+controller in one process.
- `Dockerfile` — the same image and env vars baked in, for when you want to build/run a standalone image instead of using compose.

## Dockerfile explained

```
FROM apache/kafka:latest
```
Starts from the official Apache Kafka image (bundles Kafka + a JRE, KRaft-capable).

```
ENV KAFKA_NODE_ID=1
ENV KAFKA_PROCESS_ROLES=broker,controller
```
Single node acting as both broker and controller (KRaft, no separate ZooKeeper/controller cluster).

```
ENV KAFKA_LISTENERS=PLAINTEXT://broker:9092,CONTROLLER://broker:9093,EXTERNAL://0.0.0.0:29092
ENV KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://broker:9092,EXTERNAL://localhost:29092
ENV KAFKA_CONTROLLER_LISTENER_NAMES=CONTROLLER
ENV KAFKA_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT,EXTERNAL:PLAINTEXT
ENV KAFKA_INTER_BROKER_LISTENER_NAME=PLAINTEXT
```
Three listeners:
- `PLAINTEXT` (9092) — used by other containers on the same network, advertised as `broker:9092`.
- `CONTROLLER` (9093) — KRaft controller traffic, internal only.
- `EXTERNAL` (29092) — for clients on the host machine, advertised as `localhost:29092`.

```
ENV KAFKA_CONTROLLER_QUORUM_VOTERS=1@broker:9093
```
Tells the single node it's the one and only voter in the controller quorum (node id `1`, reachable at `broker:9093`).

```
ENV KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1
ENV KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR=1
ENV KAFKA_TRANSACTION_STATE_LOG_MIN_ISR=1
```
Internal topic replication pinned to 1, since there's only one broker.

```
ENV KAFKA_GROUP_INITIAL_REBALANCE_DELAY_MS=0
```
Skip the default rebalance delay so consumer groups form instantly — fine for local dev, not for prod.

```
ENV KAFKA_NUM_PARTITIONS=3
```
Default partition count for auto-created topics.

```
EXPOSE 9092 29092
```
Documents the ports the image listens on (doesn't publish them — that still needs `-p` at run time).

**Note:** the listener config above hard-codes the hostname `broker`. A Dockerfile can't set a container's hostname at build time, so you must pass `--hostname broker` (or `--network-alias broker` on a user-defined network) when running the image standalone, or it won't bind/advertise correctly.

## docker-compose.yml explained

Same image and env vars as the Dockerfile, plus what compose adds on top:
- `container_name: broker` and the `kafka_net` bridge network give the container the hostname `broker`, which is what the listener env vars above expect.
- `ports: 9092:9092, 29092:29092` publishes both listeners to the host.

## Commands

### Docker Compose

```sh
docker compose up -d      # start
docker compose logs -f    # tail logs
docker compose down       # stop and remove
```

### Podman Compose

```sh
podman compose up -d
podman compose logs -f
podman compose down
```

### Dockerfile — build and run standalone

Docker:
```sh
docker build -t kafka-broker .
docker run -d --name broker --hostname broker \
  -p 9092:9092 -p 29092:29092 \
  kafka-broker
```

Podman:
```sh
podman build -t kafka-broker .
podman run -d --name broker --hostname broker \
  -p 9092:9092 -p 29092:29092 \
  kafka-broker
```
