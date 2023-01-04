#!/bin/bash

docker-compose -f ./docker/development/docker-compose.hot-reload.yml --env-file .env up --build