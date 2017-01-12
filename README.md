# OpenShift template service broker R&D

Work in progress.

## Prerequisites

Ensure you are logged in to an OpenShift environment before running the broker.

```bash
oc create -n openshift https://raw.githubusercontent.com/openshift/origin/master/examples/sample-app/application-template-stibuild.json
```

## Usage

In terminal 1:

```bash
go get -u github.com/jim-minter/origin-template-service-broker
$GOPATH/bin/broker
```

In terminal 2:

```bash
cd $GOPATH/src/github.com/jim-minter/origin-template-service-broker
test/catalog.sh
test/provision.sh
test/unprovision.sh
```

## Links

- [OpenShift Origin](https://github.com/openshift/origin)
- [Open Service Broker API](https://github.com/openservicebrokerapi/servicebroker)
- [Template Service Broker R&D](https://trello.com/c/lJf7723w/1117-5-template-service-broker-r-d-templates)
