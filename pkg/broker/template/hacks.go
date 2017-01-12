package template

import (
	"github.com/jim-minter/origin-template-service-broker/pkg/broker"

	template "github.com/openshift/origin/pkg/template/api"
	"github.com/pborman/uuid"
	"k8s.io/kubernetes/pkg/api/errors"
)

var templateUUIDs = []struct {
	Name      string
	Namespace string
	UUID      uuid.UUID
}{
	{
		Name:      "ruby-helloworld-sample",
		Namespace: "openshift",
		UUID:      uuid.Parse("f5e56db8-81c5-4282-940f-34c025735166"),
	},
}

func (b Broker) uuidFromTemplate(template *template.Template) uuid.UUID {
	for _, templateUUID := range templateUUIDs {
		if templateUUID.Name == template.Name && templateUUID.Namespace == template.Namespace {
			return templateUUID.UUID
		}
	}
	return uuid.NewRandom()
}

func (b Broker) templateFromUUID(u uuid.UUID) (*template.Template, error) {
	for _, templateUUID := range templateUUIDs {
		if uuid.Equal(templateUUID.UUID, u) {
			return b.oc.Templates(templateUUID.Namespace).Get(templateUUID.Name)
		}
	}
	return nil, errors.NewNotFound(template.Resource("templates"), u.String())
}

var plans = []broker.Plan{
	{
		ID:          uuid.Parse("4c10ff42-be89-420a-9bab-27a9bef9aed8"),
		Name:        "default",
		Description: "Default plan",
		Free:        true,
	},
}

// copied from github.com/openshift/origin/pkg/cmd/util/clientcmd/shortcut_restmapper.go
var userResources = []string{
	"buildconfigs", "builds",
	"imagestreams",
	"deploymentconfigs", "replicationcontrollers",
	"routes", "services",
	"pods",
}
