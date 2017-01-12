package template

import (
	"github.com/jim-minter/origin-template-service-broker/pkg/broker"

	"github.com/pborman/uuid"
)

func (b Broker) Bind(instanceUUID uuid.UUID, bindingUUID uuid.UUID, req *broker.BindRequest) (*broker.BindResponse, error) {
	// TODO: handle identical pre-existing bind (return http.StatusOK)
	return nil, notImplemented // TODO
}

func (b Broker) Unbind(instanceUUID uuid.UUID, bindingUUID uuid.UUID) error {
	return notImplemented // TODO
}
