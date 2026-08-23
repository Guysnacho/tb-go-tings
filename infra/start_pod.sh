#!/bin/bash
podman run -d --name broker --hostname broker -p 9092:9092 --replace -p 29092:29092 kafka-broker
