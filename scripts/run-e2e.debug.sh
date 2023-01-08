#!/bin/bash

MONGO_TEST_PORT=27019
MONGO_IMAGE="mongo:6"
CONTAINER_NAME="debug_mongo-e2e"
APP_SRC=$(cat .env | grep "APP_SRC" | cut -d "=" -f2)
DB_NAME="testdb"

export MONGO_URI=mongodb://localhost:$MONGO_TEST_PORT
export DB_NAME=$DB_NAME
# run mongo
CONTAINER_ID=$(docker run --rm -d -p $MONGO_TEST_PORT:27017 --name=$CONTAINER_NAME -e MONGODB_DATABASE=$DB_NAME $MONGO_IMAGE)
# run migrations
docker run -v $APP_SRC/migrations:/migrations --network host --rm migrate/migrate -path=/migrations/ -database $MONGO_URI/$DB_NAME up
# run tests
go test -count=1 -v ./tests/

printf "container: %s\n" "$CONTAINER_ID"
read -p "Press enter to remove > "

docker rm -f $CONTAINER_ID