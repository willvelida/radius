/*
Copyright 2023 The Radius Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

@description('Specifies the name of the grafana dashboard.')
param name string

@description('Specifies the sku of the grafana dashboard.')
param skuName string = 'Standard'

@description('Specifies the admin object id.')
param adminObjectId string

@description('Specifies the location.')
param location string

@description('Specifies the azure monitor workspace resource id.')
param azureMonitorWorkspaceId string

@description('Specifies the AKS cluster resource Id.')
param clusterResourceId string

@description('Specifies the AKS cluster resource location.')
param clusterLocation string

@description('Specifies comma-separated list of Kubernetes annotation keys that will be used in the resource\'s labels metric (Example: \'namespaces=[kubernetes.io/team,...],pods=[kubernetes.io/team],...\') By default the metric contains only resource name and namespace labels.')
param metricLabelsAllowlist string = ''

@description('Specifies comma-separated list of Kubernetes annotation keys that will be used in the resource\'s labels metric (Example: \'namespaces=[kubernetes.io/team,...],pods=[kubernetes.io/team],...\') By default the metric contains only resource name and namespace labels.')
param metricAnnotationsAllowList string = ''

@description('Specifies the resource tags.')
param tags object

var azureMonitorWorkspaceSubscriptionId = split(azureMonitorWorkspaceId, '/')[2]
var azureMonitorWorkspaceResourceGroupName = split(azureMonitorWorkspaceId, '/')[4]

var clusterSubscriptionId = split(clusterResourceId, '/')[2]
var clusterResourceGroup = split(clusterResourceId, '/')[4]
var clusterName = split(clusterResourceId, '/')[8]

var roleDefinitionId = {
  GrafanaAdmin: {
    id: subscriptionResourceId(azureMonitorWorkspaceSubscriptionId, 'Microsoft.Authorization/roleDefinitions', '22926164-76b3-42b3-bc55-97df8dab3e41')
  }
  MonitoringReader: {
    id: subscriptionResourceId(azureMonitorWorkspaceSubscriptionId, 'Microsoft.Authorization/roleDefinitions', '43d0d8ad-25c7-4714-9337-8ba259a9fe05')
  }
  MonitoringDataReader: {
    id: subscriptionResourceId(azureMonitorWorkspaceSubscriptionId, 'Microsoft.Authorization/roleDefinitions', 'b0d8363b-8ddd-447d-831f-62ca05bff136')
  }
}

resource grafana 'Microsoft.Dashboard/grafana@2022-08-01' = {
  name: name
  location: location
  sku: {
    name: skuName
  }
  identity: {
    type: 'SystemAssigned'
  }
  properties: {
    zoneRedundancy: 'Disabled'
    publicNetworkAccess: 'Enabled'
    autoGeneratedDomainNameLabelScope: 'TenantReuse'
    grafanaIntegrations: {
        azureMonitorWorkspaceIntegrations: [
            {
                azureMonitorWorkspaceResourceId: azureMonitorWorkspaceId
            }
        ]
    }
  }
  tags: tags
}

// Add user's as Grafana Admin for Grafana
resource adminRoleAssignment 'Microsoft.Authorization/roleAssignments@2022-04-01' = if (!empty(adminObjectId)) {
  name: guid(grafana.id, roleDefinitionId.GrafanaAdmin.id)
  scope: grafana
  properties: {
    roleDefinitionId: roleDefinitionId.GrafanaAdmin.id
    principalId: adminObjectId
  }
}

// Add user's as Grafana Admin for Grafana
resource readerRoleAssignment 'Microsoft.Authorization/roleAssignments@2022-04-01' = {
  name: guid(grafana.id, roleDefinitionId.MonitoringReader.id)
  scope: grafana
  properties: {
    roleDefinitionId: roleDefinitionId.MonitoringReader.id
    principalId: grafana.identity.principalId
  }
}

// Assign MonitoringDataReader role to Grafana identity for resourcegroup of azuremonitor workspace resource group.
module grafanaIdenityToAzureMonitor './assign-role.bicep' = {
  name: guid(grafana.id, roleDefinitionId.MonitoringDataReader.id)
  scope: resourceGroup(azureMonitorWorkspaceSubscriptionId, azureMonitorWorkspaceResourceGroupName)
  params: {
    roleNameGuid: guid(grafana.id, adminObjectId, roleDefinitionId.MonitoringDataReader.id)
    roleDefinitionId: roleDefinitionId.MonitoringDataReader.id
    principalId: grafana.identity.principalId
  }
}

module enableMonitorAddon './grafana-onboard-metrics.bicep' = {
  name: 'OnboardMetricsOnCluster-${uniqueString(clusterResourceId)}'
  scope: resourceGroup(clusterSubscriptionId, clusterResourceGroup)
  params: {
    clusterName: clusterName
    clusterLocation: clusterLocation
    metricLabelsAllowlist: metricLabelsAllowlist
    metricAnnotationsAllowList: metricAnnotationsAllowList
  }
  dependsOn: [
    grafanaIdenityToAzureMonitor
  ]
}

output dashboardFQDN string = grafana.properties.endpoint
