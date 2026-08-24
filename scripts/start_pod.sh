#!/bin/bash
set -e
podman build -t kafka-broker -f Dockerfile.kafka ../infra/.
podman build -t tigerbeetle -f Dockerfile.tigerbeetle ../infra/.
podman run -d --name broker --hostname broker -p 9092:9092 --replace -p 29092:29092 kafka-broker
podman run -d --name tigerbeetle --hostname tigerbeetle --replace -p 3000:3000 --security-opt seccomp=unconfined tigerbeetle
