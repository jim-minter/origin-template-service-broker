package template

import (
	"github.com/openshift/origin/pkg/client"
	"github.com/openshift/origin/pkg/cmd/util/clientcmd"
	"k8s.io/kubernetes/pkg/client/transport"
)

type Broker struct {
	factory        *clientcmd.Factory
	oc             *client.Client
	namespace      string
	createProjects bool
}

func NewBroker(factory *clientcmd.Factory, createProjects bool) (*Broker, error) {
	namespace, _, err := factory.DefaultNamespace()
	if err != nil {
		return nil, err
	}

	oc, _, _, err := factory.Clients()
	if err != nil {
		return nil, err
	}

	oc.Client.Transport = transport.DebugWrappers(oc.Client.Transport)

	return &Broker{
		factory:        factory,
		oc:             oc,
		namespace:      namespace,
		createProjects: createProjects,
	}, nil
}
