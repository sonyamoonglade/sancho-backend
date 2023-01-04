#!/bin/bash

MONGO_TEST_PORT=27019
MONGO_IMAGE="mongo:6"

export MONGO_URI=mongodb://localhost:$MONGO_TEST_PORT
export DB_NAME=testdb

docker run --rm -d -p $MONGO_TEST_PORT:27017 --name=$DB_NAME -e MONGODB_DATABASE=$DB_NAME $MONGO_IMAGE
go test -v ./tests/
docker stop $DB_NAME
