#!/bin/bash -e

. shared.sh

curl \
  -H 'X-Broker-API-Version: 2.9' \
  -v \
  http://localhost:8000/v2/catalog
