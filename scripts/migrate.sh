#!/usr/bin/env bash

up() {
 migrate -database "${POSTGRESQL_URL}" -path postgres/migrations up

}

down() {
 migrate -database "${POSTGRESQL_URL}" -path postgres/migrations down

}

create() {
  migrate create -ext sql -dir postgres/migrations -seq "$1"
}

"$@"