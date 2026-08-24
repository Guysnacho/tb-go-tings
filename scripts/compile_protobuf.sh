#!/bin/bash
protoc --proto_path=../protobuf --go_out=../ test.proto