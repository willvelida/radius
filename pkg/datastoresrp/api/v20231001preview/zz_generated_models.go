// Licensed under the Apache License, Version 2.0 . See LICENSE in the repository root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package v20231001preview

import "time"

// AzureContainerInstanceCompute - The Azure container instance compute configuration
type AzureContainerInstanceCompute struct {
// REQUIRED; Discriminator property for EnvironmentCompute.
	Kind *string

// Configuration for supported external identity providers
	Identity *IdentitySettings

// The resource group to use for the environment.
	ResourceGroup *string

// The resource id of the compute resource for application environment.
	ResourceID *string
}

// GetEnvironmentCompute implements the EnvironmentComputeClassification interface for type AzureContainerInstanceCompute.
func (a *AzureContainerInstanceCompute) GetEnvironmentCompute() *EnvironmentCompute {
	return &EnvironmentCompute{
		Identity: a.Identity,
		Kind: a.Kind,
		ResourceID: a.ResourceID,
	}
}

// AzureResourceManagerCommonTypesTrackedResourceUpdate - The resource model definition for an Azure Resource Manager tracked
// top level resource which has 'tags' and a 'location'
type AzureResourceManagerCommonTypesTrackedResourceUpdate struct {
// Resource tags.
	Tags map[string]*string

// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

// READ-ONLY; The name of the resource
	Name *string

// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// EnvironmentCompute - Represents backing compute resource
type EnvironmentCompute struct {
// REQUIRED; Discriminator property for EnvironmentCompute.
	Kind *string

// Configuration for supported external identity providers
	Identity *IdentitySettings

// The resource id of the compute resource for application environment.
	ResourceID *string
}

// GetEnvironmentCompute implements the EnvironmentComputeClassification interface for type EnvironmentCompute.
func (e *EnvironmentCompute) GetEnvironmentCompute() *EnvironmentCompute { return e }

// ErrorAdditionalInfo - The resource management error additional info.
type ErrorAdditionalInfo struct {
// READ-ONLY; The additional info.
	Info map[string]any

// READ-ONLY; The additional info type.
	Type *string
}

// ErrorDetail - The error detail.
type ErrorDetail struct {
// READ-ONLY; The error additional info.
	AdditionalInfo []*ErrorAdditionalInfo

// READ-ONLY; The error code.
	Code *string

// READ-ONLY; The error details.
	Details []*ErrorDetail

// READ-ONLY; The error message.
	Message *string

// READ-ONLY; The error target.
	Target *string
}

// ErrorResponse - Common error response for all Azure Resource Manager APIs to return error details for failed operations.
// (This also follows the OData error response format.).
type ErrorResponse struct {
// The error object.
	Error *ErrorDetail
}

// IdentitySettings is the external identity setting.
type IdentitySettings struct {
// REQUIRED; kind of identity setting
	Kind *IdentitySettingKind

// The list of user assigned managed identities
	ManagedIdentity []*string

// The URI for your compute platform's OIDC issuer
	OidcIssuer *string

// The resource ID of the provisioned identity
	Resource *string
}

// KubernetesCompute - The Kubernetes compute configuration
type KubernetesCompute struct {
// REQUIRED; Discriminator property for EnvironmentCompute.
	Kind *string

// REQUIRED; The namespace to use for the environment.
	Namespace *string

// Configuration for supported external identity providers
	Identity *IdentitySettings

// The resource id of the compute resource for application environment.
	ResourceID *string
}

// GetEnvironmentCompute implements the EnvironmentComputeClassification interface for type KubernetesCompute.
func (k *KubernetesCompute) GetEnvironmentCompute() *EnvironmentCompute {
	return &EnvironmentCompute{
		Identity: k.Identity,
		Kind: k.Kind,
		ResourceID: k.ResourceID,
	}
}

// MongoDatabaseListSecretsResult - The secret values for the given MongoDatabase resource
type MongoDatabaseListSecretsResult struct {
// Connection string used to connect to the target Mongo database
	ConnectionString *string

// Password to use when connecting to the target Mongo database
	Password *string
}

// MongoDatabaseProperties - MongoDatabase portable resource properties
type MongoDatabaseProperties struct {
// REQUIRED; Fully qualified resource ID for the environment that the portable resource is linked to
	Environment *string

// Fully qualified resource ID for the application that the portable resource is consumed by (if applicable)
	Application *string

// Database name of the target Mongo database
	Database *string

// Host name of the target Mongo database
	Host *string

// Port value of the target Mongo database
	Port *int32

// The recipe used to automatically deploy underlying infrastructure for the resource
	Recipe *Recipe

// Specifies how the underlying service/resource is provisioned and managed.
	ResourceProvisioning *ResourceProvisioning

// List of the resource IDs that support the MongoDB resource
	Resources []*ResourceReference

// Secret values provided for the resource
	Secrets *MongoDatabaseSecrets

// Username to use when connecting to the target Mongo database
	Username *string

// READ-ONLY; The status of the asynchronous operation.
	ProvisioningState *ProvisioningState

// READ-ONLY; Status of a resource.
	Status *ResourceStatus
}

// MongoDatabaseResource - MongoDatabase portable resource
type MongoDatabaseResource struct {
// REQUIRED; The geo-location where the resource lives
	Location *string

// REQUIRED; The resource-specific properties for this resource.
	Properties *MongoDatabaseProperties

// Resource tags.
	Tags map[string]*string

// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

// READ-ONLY; The name of the resource
	Name *string

// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// MongoDatabaseResourceListResult - The response of a MongoDatabaseResource list operation.
type MongoDatabaseResourceListResult struct {
// REQUIRED; The MongoDatabaseResource items on this page
	Value []*MongoDatabaseResource

// The link to the next page of items
	NextLink *string
}

// MongoDatabaseResourceUpdate - MongoDatabase portable resource
type MongoDatabaseResourceUpdate struct {
// Resource tags.
	Tags map[string]*string

// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

// READ-ONLY; The name of the resource
	Name *string

// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// MongoDatabaseSecrets - The secret values for the given MongoDatabase resource
type MongoDatabaseSecrets struct {
// Connection string used to connect to the target Mongo database
	ConnectionString *string

// Password to use when connecting to the target Mongo database
	Password *string
}

// Operation - Details of a REST API operation, returned from the Resource Provider Operations API
type Operation struct {
// Localized display information for this particular operation.
	Display *OperationDisplay

// READ-ONLY; Enum. Indicates the action type. "Internal" refers to actions that are for internal only APIs.
	ActionType *ActionType

// READ-ONLY; Whether the operation applies to data-plane. This is "true" for data-plane operations and "false" for ARM/control-plane
// operations.
	IsDataAction *bool

// READ-ONLY; The name of the operation, as per Resource-Based Access Control (RBAC). Examples: "Microsoft.Compute/virtualMachines/write",
// "Microsoft.Compute/virtualMachines/capture/action"
	Name *string

// READ-ONLY; The intended executor of the operation; as in Resource Based Access Control (RBAC) and audit logs UX. Default
// value is "user,system"
	Origin *Origin
}

// OperationDisplay - Localized display information for this particular operation.
type OperationDisplay struct {
// READ-ONLY; The short, localized friendly description of the operation; suitable for tool tips and detailed views.
	Description *string

// READ-ONLY; The concise, localized friendly name for the operation; suitable for dropdowns. E.g. "Create or Update Virtual
// Machine", "Restart Virtual Machine".
	Operation *string

// READ-ONLY; The localized friendly form of the resource provider name, e.g. "Microsoft Monitoring Insights" or "Microsoft
// Compute".
	Provider *string

// READ-ONLY; The localized friendly name of the resource type related to this operation. E.g. "Virtual Machines" or "Job
// Schedule Collections".
	Resource *string
}

// OperationListResult - A list of REST API operations supported by an Azure Resource Provider. It contains an URL link to
// get the next set of results.
type OperationListResult struct {
// READ-ONLY; URL to get the next set of operation list results (if there are any).
	NextLink *string

// READ-ONLY; List of operations supported by the resource provider
	Value []*Operation
}

// OutputResource - Properties of an output resource.
type OutputResource struct {
// The UCP resource ID of the underlying resource.
	ID *string

// The logical identifier scoped to the owning Radius resource. This is only needed or used when a resource has a dependency
// relationship. LocalIDs do not have any particular format or meaning beyond
// being compared to determine dependency relationships.
	LocalID *string

// Determines whether Radius manages the lifecycle of the underlying resource.
	RadiusManaged *bool
}

// Recipe - The recipe used to automatically deploy underlying infrastructure for a portable resource
type Recipe struct {
// REQUIRED; The name of the recipe within the environment to use
	Name *string

// Key/value parameters to pass into the recipe at deployment
	Parameters map[string]any
}

// RecipeStatus - Recipe status at deployment time for a resource.
type RecipeStatus struct {
// REQUIRED; TemplateKind is the kind of the recipe template used by the portable resource upon deployment.
	TemplateKind *string

// REQUIRED; TemplatePath is the path of the recipe consumed by the portable resource upon deployment.
	TemplatePath *string

// TemplateVersion is the version number of the template.
	TemplateVersion *string
}

// RedisCacheListSecretsResult - The secret values for the given RedisCache resource
type RedisCacheListSecretsResult struct {
// The connection string used to connect to the Redis cache
	ConnectionString *string

// The password for this Redis cache instance
	Password *string

// The URL used to connect to the Redis cache
	URL *string
}

// RedisCacheProperties - RedisCache portable resource properties
type RedisCacheProperties struct {
// REQUIRED; Fully qualified resource ID for the environment that the portable resource is linked to
	Environment *string

// Fully qualified resource ID for the application that the portable resource is consumed by (if applicable)
	Application *string

// The host name of the target Redis cache
	Host *string

// The port value of the target Redis cache
	Port *int32

// The recipe used to automatically deploy underlying infrastructure for the resource
	Recipe *Recipe

// Specifies how the underlying service/resource is provisioned and managed.
	ResourceProvisioning *ResourceProvisioning

// List of the resource IDs that support the Redis resource
	Resources []*ResourceReference

// Secrets provided by resource
	Secrets *RedisCacheSecrets

// Specifies whether to enable SSL connections to the Redis cache
	TLS *bool

// The username for Redis cache
	Username *string

// READ-ONLY; The status of the asynchronous operation.
	ProvisioningState *ProvisioningState

// READ-ONLY; Status of a resource.
	Status *ResourceStatus
}

// RedisCacheResource - RedisCache portable resource
type RedisCacheResource struct {
// REQUIRED; The geo-location where the resource lives
	Location *string

// REQUIRED; The resource-specific properties for this resource.
	Properties *RedisCacheProperties

// Resource tags.
	Tags map[string]*string

// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

// READ-ONLY; The name of the resource
	Name *string

// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// RedisCacheResourceListResult - The response of a RedisCacheResource list operation.
type RedisCacheResourceListResult struct {
// REQUIRED; The RedisCacheResource items on this page
	Value []*RedisCacheResource

// The link to the next page of items
	NextLink *string
}

// RedisCacheResourceUpdate - RedisCache portable resource
type RedisCacheResourceUpdate struct {
// Resource tags.
	Tags map[string]*string

// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

// READ-ONLY; The name of the resource
	Name *string

// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// RedisCacheSecrets - The secret values for the given RedisCache resource
type RedisCacheSecrets struct {
// The connection string used to connect to the Redis cache
	ConnectionString *string

// The password for this Redis cache instance
	Password *string

// The URL used to connect to the Redis cache
	URL *string
}

// Resource - Common fields that are returned in the response for all Azure Resource Manager resources
type Resource struct {
// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

// READ-ONLY; The name of the resource
	Name *string

// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// ResourceReference - Describes a reference to an existing resource
type ResourceReference struct {
// REQUIRED; Resource id of an existing resource
	ID *string
}

// ResourceStatus - Status of a resource.
type ResourceStatus struct {
// The compute resource associated with the resource.
	Compute EnvironmentComputeClassification

// Properties of an output resource
	OutputResources []*OutputResource

// READ-ONLY; The recipe data at the time of deployment
	Recipe *RecipeStatus
}

// SQLDatabaseListSecretsResult - The secret values for the given SqlDatabase resource
type SQLDatabaseListSecretsResult struct {
// Connection string used to connect to the target Sql database
	ConnectionString *string

// Password to use when connecting to the target Sql database
	Password *string
}

// SQLDatabaseProperties - SqlDatabase properties
type SQLDatabaseProperties struct {
// REQUIRED; Fully qualified resource ID for the environment that the portable resource is linked to
	Environment *string

// Fully qualified resource ID for the application that the portable resource is consumed by (if applicable)
	Application *string

// The name of the Sql database.
	Database *string

// Port value of the target Sql database
	Port *int32

// The recipe used to automatically deploy underlying infrastructure for the resource
	Recipe *Recipe

// Specifies how the underlying service/resource is provisioned and managed.
	ResourceProvisioning *ResourceProvisioning

// List of the resource IDs that support the SqlDatabase resource
	Resources []*ResourceReference

// Secret values provided for the resource
	Secrets *SQLDatabaseSecrets

// The fully qualified domain name of the Sql database.
	Server *string

// Username to use when connecting to the target Sql database
	Username *string

// READ-ONLY; The status of the asynchronous operation.
	ProvisioningState *ProvisioningState

// READ-ONLY; Status of a resource.
	Status *ResourceStatus
}

// SQLDatabaseResource - SqlDatabase portable resource
type SQLDatabaseResource struct {
// REQUIRED; The geo-location where the resource lives
	Location *string

// REQUIRED; The resource-specific properties for this resource.
	Properties *SQLDatabaseProperties

// Resource tags.
	Tags map[string]*string

// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

// READ-ONLY; The name of the resource
	Name *string

// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// SQLDatabaseResourceListResult - The response of a SqlDatabaseResource list operation.
type SQLDatabaseResourceListResult struct {
// REQUIRED; The SqlDatabaseResource items on this page
	Value []*SQLDatabaseResource

// The link to the next page of items
	NextLink *string
}

// SQLDatabaseResourceUpdate - SqlDatabase portable resource
type SQLDatabaseResourceUpdate struct {
// Resource tags.
	Tags map[string]*string

// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

// READ-ONLY; The name of the resource
	Name *string

// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// SQLDatabaseSecrets - The secret values for the given SqlDatabase resource
type SQLDatabaseSecrets struct {
// Connection string used to connect to the target Sql database
	ConnectionString *string

// Password to use when connecting to the target Sql database
	Password *string
}

// SystemData - Metadata pertaining to creation and last modification of the resource.
type SystemData struct {
// The timestamp of resource creation (UTC).
	CreatedAt *time.Time

// The identity that created the resource.
	CreatedBy *string

// The type of identity that created the resource.
	CreatedByType *CreatedByType

// The timestamp of resource last modification (UTC)
	LastModifiedAt *time.Time

// The identity that last modified the resource.
	LastModifiedBy *string

// The type of identity that last modified the resource.
	LastModifiedByType *CreatedByType
}

// TrackedResource - The resource model definition for an Azure Resource Manager tracked top level resource which has 'tags'
// and a 'location'
type TrackedResource struct {
// REQUIRED; The geo-location where the resource lives
	Location *string

// Resource tags.
	Tags map[string]*string

// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

// READ-ONLY; The name of the resource
	Name *string

// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

