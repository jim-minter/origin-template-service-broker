package template

import (
	"strings"

	"github.com/jim-minter/origin-template-service-broker/pkg/broker"
	"github.com/pborman/uuid"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/kubelet/envvars"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/selection"
	"k8s.io/kubernetes/pkg/util/sets"
)

func (b Broker) Bind(instanceUUID uuid.UUID, bindingUUID uuid.UUID, req *broker.BindRequest) (*broker.BindResponse, error) {
	namespace := b.namespace
	if b.createProjects {
		namespace = instanceUUID.String()
	}

	credentials := map[string]interface{}{}

	configmap, err := b.kc.ConfigMaps(namespace).Get(templateConfigMapPrefix + instanceUUID.String())
	if err != nil {
		return nil, err
	}

	for k, v := range configmap.Data {
		kl := strings.ToLower(k)
		// TODO: annotating "exportable" parameters in templates may be a good idea
		if strings.Contains(kl, "user") || strings.Contains(kl, "password") {
			credentials[k] = v
		}
	}

	requirement, _ := labels.NewRequirement(serviceInstanceLabel, selection.Equals, sets.NewString(instanceUUID.String()))
	selector := labels.NewSelector().Add(*requirement)

	serviceList, err := b.kc.Services(namespace).List(kapi.ListOptions{LabelSelector: selector})
	if err != nil {
		return nil, err
	}

	for _, envvar := range envvars.FromServices(serviceList) {
		// TODO: we're somewhat overloading 'credentials' here
		credentials[envvar.Name] = envvar.Value
	}

	// TODO: handle identical pre-existing bind (return http.StatusOK)
	return &broker.BindResponse{Credentials: credentials}, nil
}

func (b Broker) Unbind(instanceUUID uuid.UUID, bindingUUID uuid.UUID) error {
	return nil
}
