package template

import (
	"github.com/jim-minter/origin-template-service-broker/pkg/broker"
	"strings"

	template "github.com/openshift/origin/pkg/template/api"
	kapi "k8s.io/kubernetes/pkg/api"
)

func (b Broker) serviceFromTemplate(template *template.Template) *broker.Service {
	return &broker.Service{
		Name:        template.Name,
		ID:          b.uuidFromTemplate(template), // TODO: have the apiserver generate a fixed UUID upon template admission
		Description: template.Annotations["description"],
		Tags:        strings.Split(template.Annotations["tags"], ","),
		Bindable:    false, // TODO
		Metadata: map[string]interface{}{
			"displayName": template.Annotations["openshift.io/display-name"],
			// TODO: "imageUrl":            "",
			// TODO: "longDescription":     "",
			// TODO: "providerDisplayName": "",
			// TODO: "documentationUrl":    "",
			// TODO: "supportUrl":          "",
		},
		PlanUpdatable: false, // TODO
		Plans:         plans, // TODO
	}
}

func (b Broker) Catalog() (*broker.CatalogResponse, error) {
	templates, err := b.oc.Templates("openshift").List(kapi.ListOptions{})
	if err != nil {
		return nil, err
	}

	services := []broker.Service{}
	for _, template := range templates.Items {
		services = append(services, *b.serviceFromTemplate(&template))
	}

	return &broker.CatalogResponse{Services: services}, nil
}
