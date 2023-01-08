#!/bin/bash

MONGO_TEST_PORT=27019
MONGO_IMAGE="mongo:6"
CONTAINER_NAME="debug_mongo-e2e"

export MONGO_URI=mongodb://localhost:$MONGO_TEST_PORT
export DB_NAME=testdb

CONTAINER_ID=$(docker run --rm -d -p $MONGO_TEST_PORT:27017 --name=$CONTAINER_NAME -e MONGODB_DATABASE=$DB_NAME $MONGO_IMAGE)
go test -v ./tests/
printf "container: %s\n" "$CONTAINER_ID"
read -p "Press enter to remove > "

docker rm -f $CONTAINER_ID