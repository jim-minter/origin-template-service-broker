package template

import (
	"github.com/jim-minter/origin-template-service-broker/pkg/broker"
	"github.com/jim-minter/origin-template-service-broker/pkg/errors"
	"github.com/openshift/origin/pkg/cmd/util/clientcmd"
	"github.com/openshift/origin/pkg/config/cmd"
	projectapi "github.com/openshift/origin/pkg/project/api"
	templateapi "github.com/openshift/origin/pkg/template/api"
	"github.com/pborman/uuid"
	kapi "k8s.io/kubernetes/pkg/api"
	kerrors "k8s.io/kubernetes/pkg/api/errors"
	"k8s.io/kubernetes/pkg/api/meta"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/selection"
	"k8s.io/kubernetes/pkg/util/sets"
)

const serviceInstanceLabel = "service-instance"
const templateConfigMapPrefix = "service-instance-"

var accessor = meta.NewAccessor()

func projectRequestFromUUID(u uuid.UUID) *projectapi.ProjectRequest {
	return &projectapi.ProjectRequest{
		ObjectMeta: kapi.ObjectMeta{Name: u.String()},
	}
}

func (b Broker) templateFromUUID(u uuid.UUID) (*templateapi.Template, error) {
	// TODO: cache this
	templateList, err := b.oc.Templates("openshift").List(kapi.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, template := range templateList.Items {
		if uuid.Equal(uuid.Parse(string(template.GetUID())), u) {
			return &template, nil
		}
	}

	return nil, kerrors.NewNotFound(templateapi.Resource("templates"), u.String())
}

func (b Broker) Provision(instanceUUID uuid.UUID, req *broker.ProvisionRequest) (*broker.ProvisionResponse, error) {
	if !uuid.Equal(req.PlanID, plans[0].ID) { // TODO: think about plans and implement them
		return nil, kerrors.NewBadRequest("invalid plan_id")
	}

	template, err := b.templateFromUUID(req.ServiceID)
	if err != nil {
		return nil, err
	}

	for i, param := range template.Parameters {
		if value, ok := req.Parameters[param.Name]; ok {
			template.Parameters[i].Value = value
			template.Parameters[i].Generate = ""
		}
	}

	template.ObjectLabels[serviceInstanceLabel] = instanceUUID.String()

	template, err = b.oc.TemplateConfigs(b.namespace).Create(template)
	if err != nil {
		return nil, err
	}

	// TODO: discuss what the appropriate place for this state is
	configmap := &kapi.ConfigMap{
		ObjectMeta: kapi.ObjectMeta{
			Name: templateConfigMapPrefix + instanceUUID.String(),
			Labels: map[string]string{
				serviceInstanceLabel: instanceUUID.String(),
			},
		},
		Data: map[string]string{},
	}

	for _, param := range template.Parameters {
		// TODO: check that any valid param.Name is a valid ConfigMap data key
		configmap.Data[param.Name] = param.Value
	}

	errs := runtime.DecodeList(template.Objects, kapi.Codecs.UniversalDecoder())
	if len(errs) > 0 {
		return nil, errors.Errors(errs)
	}

	// TODO: if !b.createProjects, need to signal a namespace and presumably a
	// user for impersonation
	namespace := b.namespace
	if b.createProjects {
		project, err := b.oc.ProjectRequests().Create(projectRequestFromUUID(instanceUUID))
		if err != nil {
			// TODO: handle kerrors.IsAlreadyExists(err) and identical pre-existing project (return http.StatusOK)
			return nil, err
		}
		namespace = project.Name
	}

	bulk := &cmd.Bulk{
		Mapper: clientcmd.ResourceMapper(b.factory),
		Op:     cmd.Create,
	}
	errs = bulk.Run(&kapi.List{Items: template.Objects}, namespace)
	if len(errs) > 0 {
		return nil, errors.Errors(errs)
	}

	configmap, err = b.kc.ConfigMaps(namespace).Create(configmap)
	if err != nil {
		return nil, err
	}

	// TODO: rollback on error?
	// TODO: wait for the template to finish deploying?

	return &broker.ProvisionResponse{}, nil
}

func (b Broker) Update(instanceUUID uuid.UUID, req *broker.UpdateRequest) (*broker.UpdateResponse, error) {
	return nil, notImplemented // TODO
}

func (b Broker) Deprovision(instanceUUID uuid.UUID) (*broker.DeprovisionResponse, error) {
	if b.createProjects {
		namespace := instanceUUID.String()

		err := b.oc.Projects().Delete(namespace)
		if err != nil {
			return nil, err
		}

	} else {
		// TODO: what follows is horrible

		requirement, _ := labels.NewRequirement(serviceInstanceLabel, selection.Equals, sets.NewString(instanceUUID.String()))
		selector := labels.NewSelector().Add(*requirement)

		resourcemapper := clientcmd.ResourceMapper(b.factory)

		var errs []error
		for _, resource := range userResources {
			gvk, _ := resourcemapper.KindFor(unversioned.GroupVersionResource{Resource: resource})
			restmapping, _ := resourcemapper.RESTMapping(unversioned.GroupKind{Group: gvk.Group, Kind: gvk.Kind})
			cli, _ := b.factory.ClientForMapping(restmapping)

			obj, err := cli.Get().Namespace(b.namespace).Resource(resource).LabelsSelectorParam(selector).Do().Get()
			if err != nil {
				errs = append(errs, err)
				continue
			}
			objs, _ := meta.ExtractList(obj)

			for _, obj = range objs {
				name, _ := accessor.Name(obj)

				reaper, _ := b.factory.Reaper(restmapping)
				if reaper != nil {
					// TODO: DeploymentControllers don't transfer their labels to their
					// -deploy / -hook-* pods, so this doesn't clear them up.
					reaper.Stop(b.namespace, name, 0, kapi.NewDeleteOptions(0))

				} else {
					cli.Delete().Namespace(b.namespace).Resource(resource).Name(name).Do()
				}
			}
		}
	}

	return &broker.DeprovisionResponse{}, nil
}
