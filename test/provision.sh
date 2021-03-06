#!/bin/bash -e

. shared.sh

planUUID=${planUUID-$defaultplanUUID}
serviceUUID=${serviceUUID-$(oc get template ruby-helloworld-sample -n openshift -o template --template '{{.metadata.uid}}')}

req="{
  \"plan_id\": \"$planUUID\",
  \"service_id\": \"$serviceUUID\",
  \"parameters\": {
    \"MYSQL_USER\": \"username\"
  }
}"

curl \
  -X PUT \
  -H 'X-Broker-API-Version: 2.9' \
  -H 'Content-Type: application/json' \
  -d "$req" \
  -v \
  http://localhost:8000/v2/service_instances/$instanceUUID
