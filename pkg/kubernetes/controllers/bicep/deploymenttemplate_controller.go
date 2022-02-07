// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------

package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/project-radius/radius/pkg/azure/azresources"
	"github.com/project-radius/radius/pkg/cli/armtemplate"
	"github.com/project-radius/radius/pkg/cli/armtemplate/providers"
	"github.com/project-radius/radius/pkg/kubernetes"
	bicepv1alpha3 "github.com/project-radius/radius/pkg/kubernetes/api/bicep/v1alpha3"
	radiusv1alpha3 "github.com/project-radius/radius/pkg/kubernetes/api/radius/v1alpha3"
	"github.com/project-radius/radius/pkg/kubernetes/converters"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	ConditionReady = "Ready"
)

// DeploymentTemplateReconciler reconciles a Arm object
type DeploymentTemplateReconciler struct {
	client.Client
	meta.RESTMapper
	DynamicClient dynamic.Interface
	Log           logr.Logger
	Scheme        *runtime.Scheme
	Recorder      record.EventRecorder
}

//+kubebuilder:rbac:groups=bicep.dev,resources=deploymenttemplates,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=bicep.dev,resources=deploymenttemplates/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=bicep.dev,resources=deploymenttemplates/finalizers,verbs=update

func (r *DeploymentTemplateReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("deploymenttemplate", req.NamespacedName)

	arm := &bicepv1alpha3.DeploymentTemplate{}
	err := r.Get(ctx, req.NamespacedName, arm)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Update observed generation
	// We know for sure that the resource is currently being provisioned.
	if arm.Generation != arm.Status.ObservedGeneration {
		r.StatusProvisioned(ctx, arm, ConditionReady)
	}

	templateCondition := meta.FindStatusCondition(arm.Status.Conditions, ConditionReady)

	if len(arm.Status.Conditions) > 0 && templateCondition != nil && templateCondition.Status == metav1.ConditionTrue {
		// Template has already deployed, don't do anything
		r.Log.Info("template is already deployed")
		return ctrl.Result{}, nil
	}

	result, err := r.ApplyState(ctx, req, arm)

	// Always try to update status even if there was a failure.
	_ = r.Status().Update(ctx, arm)
	if err != nil {
		return ctrl.Result{}, err
	}

	return result, err
}

func (r *DeploymentTemplateReconciler) evaluateNestedDeployment(req ctrl.Request, evaluator *armtemplate.DeploymentEvaluator, resource armtemplate.Resource, parentName string) (*unstructured.Unstructured, error) {
	body, err := evaluator.VisitResourceBody(resource)
	if err != nil {
		return nil, err
	}
	resource.Body = body

	return armtemplate.ConvertToK8sDeploymentTemplate(resource, req.NamespacedName.Namespace, parentName)
}

func (r *DeploymentTemplateReconciler) evaluateResource(req ctrl.Request, evaluator *armtemplate.DeploymentEvaluator, resource armtemplate.Resource) (*unstructured.Unstructured, map[string]string, error) {
	body, err := evaluator.VisitResourceBody(resource)
	if err != nil {
		return nil, nil, err
	}
	resource.Body = body
	return armtemplate.ConvertToK8s(resource, req.NamespacedName.Namespace)
}

// Parses the arm template and deploys individual resources to the cluster
// TODO: Can we avoid parsing resources multiple times by caching?
func (r *DeploymentTemplateReconciler) ApplyState(ctx context.Context, req ctrl.Request, arm *bicepv1alpha3.DeploymentTemplate) (ctrl.Result, error) {
	template, err := armtemplate.Parse(string(arm.Spec.Content.Raw))
	if err != nil {
		return ctrl.Result{}, err
	}

	parameters := map[string]map[string]interface{}{}
	if arm.Spec.Parameters != nil {
		err = json.Unmarshal(arm.Spec.Parameters.Raw, &parameters)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	options := armtemplate.TemplateOptions{
		SubscriptionID: "kubernetes",
		ResourceGroup:  req.Namespace,
		Parameters:     parameters,
	}
	resources, err := armtemplate.Eval(ctx, template, options)
	if err != nil {
		return ctrl.Result{}, err
	}
	// All previously deployed resources to be used by other resources
	// to fill in variables ex: ([reference(...)])
	deployed := map[string]map[string]interface{}{}
	evaluator := &armtemplate.DeploymentEvaluator{
		Context:   ctx,
		Template:  template,
		Options:   options,
		Deployed:  deployed,
		Variables: map[string]interface{}{},
		Outputs:   map[string]map[string]interface{}{},

		CustomActionCallback: func(id string, apiVersion string, action string, payload interface{}) (interface{}, error) {
			return r.InvokeCustomAction(ctx, req.Namespace, id, apiVersion, action, payload)
		},
		Providers: map[string]providers.Provider{
			// TODO(rynowak): right now this is only used for GETs on existing resources
			// it only supports Radius types.
			providers.KubernetesProviderImport: providers.NewK8sProvider(r.Log, r.DynamicClient, r.RESTMapper),
			providers.RadiusProviderImport:     providers.NewK8sProvider(r.Log, r.DynamicClient, r.RESTMapper),
		},
	}

	for name, variable := range template.Variables {
		value, err := evaluator.VisitValue(variable)
		if err != nil {
			return ctrl.Result{}, err
		}

		evaluator.Variables[name] = value
	}

	for _, resource := range resources {
		var k8sInfo *unstructured.Unstructured
		var scrapedSecrets map[string]string
		var err error
		if resource.Type == armtemplate.DeploymentResourceType {
			k8sInfo, err = r.evaluateNestedDeployment(req, evaluator, resource, arm.Name)
		} else {
			k8sInfo, scrapedSecrets, err = r.evaluateResource(req, evaluator, resource)
		}
		if err != nil {
			return ctrl.Result{}, err
		}

		// Set name of deployment template
		annotations := k8sInfo.GetAnnotations()
		if annotations == nil {
			annotations = make(map[string]string)
		}
		annotations[kubernetes.LabelRadiusDeployment] = arm.Name
		k8sInfo.SetAnnotations(annotations)

		// All radius.dev resources, except for Application, require an Application to be created first.
		//
		// Other resources like K8s extension resources will not need the same.
		gvk := k8sInfo.GroupVersionKind()
		if gvk.Group == radiusv1alpha3.GroupVersion.Group && gvk.Kind != "Application" {
			application := &radiusv1alpha3.Application{}

			err = r.Client.Get(ctx, client.ObjectKey{
				Namespace: k8sInfo.GetNamespace(),
				Name:      k8sInfo.GetAnnotations()[kubernetes.LabelRadiusApplication],
			}, application)
			if err != nil {
				return ctrl.Result{}, err
			}

			err := controllerutil.SetControllerReference(application, k8sInfo, r.Scheme)
			_, ok := err.(*controllerutil.AlreadyOwnedError) // Ignore already owned error as if the resource is already created, it will be owned
			if err != nil && !ok {
				return ctrl.Result{}, err
			}
		}

		r.Recorder.Eventf(k8sInfo, "Normal", "Provisioned", "Resource %s has been provisioned", k8sInfo.GetName())
		r.StatusProvisionedResource(ctx, resource.ID, arm, k8sInfo)

		// Always patch the resource, even if it already exists.
		err = r.Patch(ctx, k8sInfo, client.Apply, &client.PatchOptions{FieldManager: kubernetes.FieldManager})
		if err != nil {
			return ctrl.Result{}, err
		}

		err = r.Get(ctx, client.ObjectKey{
			Namespace: k8sInfo.GetNamespace(),
			Name:      k8sInfo.GetName(),
		}, k8sInfo)

		if err != nil {
			return ctrl.Result{}, err
		}
		// Now store secret we scraped from the rendered template.
		if len(scrapedSecrets) > 0 {
			secret := kubernetes.MakeScrapedSecret(k8sInfo, scrapedSecrets)
			err := controllerutil.SetControllerReference(k8sInfo, secret, r.Scheme)
			_, ok := err.(*controllerutil.AlreadyOwnedError) // Ignore already owned error as if the resource is already created, it will be owned
			if err != nil && !ok {
				return ctrl.Result{}, err
			}
			err = r.Client.Patch(ctx, secret, client.Apply, &client.PatchOptions{FieldManager: kubernetes.FieldManager})
			if err != nil {
				return ctrl.Result{}, err
			}
		}

		// TODO could remove this dependecy on radiusv1alpha3
		k8sResource := &radiusv1alpha3.Resource{}
		err = runtime.DefaultUnstructuredConverter.FromUnstructured(k8sInfo.Object, k8sResource)
		if err != nil {
			return ctrl.Result{}, err
		}
		// Transform from k8s representation to arm representation
		//
		// We need to overlay stateful properties over the original definition.
		//
		// For now we just modify the body in place.
		err = converters.ConvertToARMResource(k8sResource, resource.Body)
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to convert to ARM representation: %w", err)
		}

		deployed[resource.ID] = resource.Body

		switch gvk.Group {
		case radiusv1alpha3.GroupVersion.Group, bicepv1alpha3.GroupVersion.Group:
			// If this resource is a Radius resource, we have a complete understanding of its
			// Readiness semantic, and we will use that to wait for readiness here.
			resourceStatus := meta.FindStatusCondition(k8sResource.Status.Conditions, ConditionReady)
			if len(k8sResource.Status.Conditions) == 0 || (resourceStatus != nil && resourceStatus.Status != metav1.ConditionTrue) {
				// Resource is not ready, wait until they are.
				return ctrl.Result{}, nil
			}
		default:
			// In the future, other resources' Readiness semantic (like Deployment, Service) may
			// be added here.
		}

		if resource.Type == armtemplate.DeploymentResourceType {
			outputMap := map[string]interface{}{}
			deployed[resource.ID]["properties"].(map[string]interface{})["outputs"] = outputMap

			deResource := &bicepv1alpha3.DeploymentTemplate{}
			err = runtime.DefaultUnstructuredConverter.FromUnstructured(k8sInfo.Object, deResource)
			if err != nil {
				return ctrl.Result{}, err
			}

			for k, v := range deResource.Status.Outputs {
				outputMap[k] = map[string]interface{}{
					"value": v.Value,
					"type":  v.Type,
				}
			}
		}

		r.Recorder.Eventf(k8sInfo, "Normal", "Deployed", "Resource %s has been deployed", k8sInfo.GetName())
		r.StatusDeployedResource(ctx, resource.ID, arm, k8sInfo)
	}

	outputs, err := evaluator.EvaluateOutputs()
	if err != nil {
		return ctrl.Result{}, err
	}

	arm.Status.Outputs = make(map[string]bicepv1alpha3.DeploymentOutput)
	for k, t := range outputs {
		if err != nil {
			return ctrl.Result{}, err
		}

		arm.Status.Outputs[k] = bicepv1alpha3.DeploymentOutput{
			Type:  t["type"].(string),
			Value: t["value"].(string),
		}
	}

	// All resources have been deployed, update status to be Deployed
	r.Recorder.Eventf(arm, "Normal", "Deployed", "Deployment Template %s has been deployed", arm.GetName())
	r.StatusDeployed(ctx, arm, ConditionReady)

	return ctrl.Result{}, nil
}

func (r *DeploymentTemplateReconciler) InvokeCustomAction(ctx context.Context, namespace string, id string, apiVersion string, action string, payload interface{}) (interface{}, error) {
	if action != "listSecrets" {
		return nil, fmt.Errorf("only %q is supported", "listSecrets")
	}

	// We can ignore ID in this case because it reference to the Radius Custom RP name ('radiusv3')
	// The resource ID we actually want is inside the payload.
	type ListSecretsInput = struct {
		TargetID string `json:"targetID"`
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.New("failed to read listSecrets payload")
	}

	input := ListSecretsInput{}
	err = json.Unmarshal(b, &input)
	if err != nil {
		return nil, errors.New("failed to read listSecrets payload")
	}

	targetID, err := azresources.Parse(input.TargetID)
	if err != nil {
		return nil, fmt.Errorf("resource id %q is invalid: %w", id, err)
	}

	if len(targetID.Types) != 3 {
		return nil, fmt.Errorf("resource id must refer to a Radius resource, was: %q", id)
	}

	unst := unstructured.Unstructured{}
	kind, ok := armtemplate.GetKindFromArmType(targetID.Types[2].Type)
	if !ok {
		return nil, fmt.Errorf("resource id does not have matching kind was: %q", targetID.Types[2].Type)
	}

	unst.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "radius.dev",
		Version: "v1alpha3",
		Kind:    kind,
	})

	err = r.Client.Get(ctx, types.NamespacedName{Namespace: namespace, Name: kubernetes.MakeResourceName(targetID.Types[1].Name, targetID.Types[2].Name)}, &unst)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve resource matching id %q: %w", id, err)
	}

	resource := radiusv1alpha3.Resource{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unst.Object, &resource)
	if err != nil {
		return nil, err
	}

	secretValues, err := converters.GetSecretValues(resource.Status)
	if err != nil {
		return nil, err
	}

	secretClient := converters.SecretClient{Client: r.Client}
	values := map[string]interface{}{}
	for key, reference := range secretValues {
		value, err := secretClient.LookupSecretValue(ctx, resource.Status, reference)
		if err != nil {
			return nil, err
		}

		values[key] = value
	}

	return values, nil
}

func (r *DeploymentTemplateReconciler) StatusProvisioned(ctx context.Context, arm *bicepv1alpha3.DeploymentTemplate, conditionType string) {
	r.Log.Info("updating status to provisioned deployment template")
	arm.Status.Conditions = []metav1.Condition{}
	arm.Status.ObservedGeneration = arm.Generation
	arm.Status.Phrase = "Provisioned"

	newCondition := metav1.Condition{
		Status:             metav1.ConditionUnknown,
		Reason:             "Provisioned",
		Type:               conditionType,
		Message:            "provisioned deployment template",
		ObservedGeneration: arm.Generation,
	}

	meta.SetStatusCondition(&arm.Status.Conditions, newCondition)

}

func (r *DeploymentTemplateReconciler) StatusDeployed(ctx context.Context, arm *bicepv1alpha3.DeploymentTemplate, conditionType string) {
	r.Log.Info("updating status to deployed deployment template")
	newCondition := metav1.Condition{
		Status:             metav1.ConditionTrue,
		Type:               conditionType,
		Reason:             "Deployed",
		Message:            "deployed deployment template",
		ObservedGeneration: arm.Generation,
	}

	meta.SetStatusCondition(&arm.Status.Conditions, newCondition)
	arm.Status.Phrase = "Deployed"
}

func (r *DeploymentTemplateReconciler) StatusProvisionedResource(ctx context.Context, id string, arm *bicepv1alpha3.DeploymentTemplate, unst *unstructured.Unstructured) {
	newResourceStatus := bicepv1alpha3.ResourceStatus{
		Status:     metav1.ConditionUnknown,
		Name:       unst.GetName(),
		Kind:       unst.GetKind(),
		ResourceID: id,
	}
	r.setResourceStatus(&arm.Status.ResourceStatuses, newResourceStatus)
}

func (r *DeploymentTemplateReconciler) StatusDeployedResource(ctx context.Context, id string, arm *bicepv1alpha3.DeploymentTemplate, unst *unstructured.Unstructured) {
	newResourceStatus := bicepv1alpha3.ResourceStatus{
		Status:     metav1.ConditionTrue,
		Name:       unst.GetName(),
		Kind:       unst.GetKind(),
		ResourceID: id,
	}
	r.setResourceStatus(&arm.Status.ResourceStatuses, newResourceStatus)
}

func (r *DeploymentTemplateReconciler) setResourceStatus(resourceStatuses *[]bicepv1alpha3.ResourceStatus, status bicepv1alpha3.ResourceStatus) {
	if resourceStatuses == nil {
		return
	}
	existingStatus := r.findResourceStatus(*resourceStatuses, status)
	if existingStatus == nil {
		*resourceStatuses = append(*resourceStatuses, status)
		return
	}

	if existingStatus.Status != status.Status {
		existingStatus.Status = status.Status
	}
}

func (r *DeploymentTemplateReconciler) findResourceStatus(resourceStatuses []bicepv1alpha3.ResourceStatus, status bicepv1alpha3.ResourceStatus) *bicepv1alpha3.ResourceStatus {
	for i := range resourceStatuses {
		if resourceStatuses[i].Name == status.Name && resourceStatuses[i].Kind == status.Kind {
			return &resourceStatuses[i]
		}
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DeploymentTemplateReconciler) SetupWithManager(mgr ctrl.Manager, objs []struct {
	client.Object
	client.ObjectList
}) error {
	// Watch for the application & deployment changes on top of other resources
	watchTypes := []struct {
		client.Object
		client.ObjectList
	}{{
		&radiusv1alpha3.Application{},
		&radiusv1alpha3.ApplicationList{},
	}, {
		&bicepv1alpha3.DeploymentTemplate{},
		&bicepv1alpha3.DeploymentTemplateList{},
	}}

	objs = append(objs, watchTypes...)

	c := ctrl.NewControllerManagedBy(mgr).
		For(&bicepv1alpha3.DeploymentTemplate{})
	for _, obj := range objs {
		resourceSource := &source.Kind{Type: obj.Object}
		handler := handler.EnqueueRequestsFromMapFunc(func(clientObj client.Object) []ctrl.Request {
			annotations := clientObj.GetAnnotations()
			template := annotations[kubernetes.LabelRadiusDeployment]
			if template == "" {
				return nil
			}

			return []ctrl.Request{
				{NamespacedName: types.NamespacedName{Namespace: clientObj.GetNamespace(), Name: template}},
			}
		})

		c = c.Watches(resourceSource, handler)
	}

	return c.Complete(r)
}
