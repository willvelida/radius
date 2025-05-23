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

import "@typespec/rest";
import "@typespec/versioning";
import "@typespec/openapi";
import "@azure-tools/typespec-autorest";
import "@azure-tools/typespec-azure-core";
import "@azure-tools/typespec-azure-resource-manager";

import "../radius/v1/ucprootscope.tsp";
import "../radius/v1/resources.tsp";
import "./common.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;
using TypeSpec.Versioning;
using Autorest;
using Azure.Core;
using Azure.ResourceManager;
using Azure.ResourceManager.Foundations;
using OpenAPI;

@doc("The UCP HTTP request base parameters.")
model CredentialBaseParameters<TResource, TPlaneName = {}> {
  ...ApiVersionParameter;
  ...TPlaneName;
  ...KeysOf<TResource>;
}

@doc("Credential store kinds supported.")
enum CredentialStorageKind {
  @doc("Internal credential storage")
  Internal,

  // This was added because typespec doesn't generate necessary functions and other variables for single value enums.
  // This can be deleted if we add more values to the enum.
  string,
}

@doc("The base credential storage properties")
@discriminator("kind")
model CredentialStorageProperties {
  @doc("The kind of credential storage")
  kind: CredentialStorageKind;
}

@doc("Internal credential storage properties")
model InternalCredentialStorageProperties extends CredentialStorageProperties {
  @doc("Internal credential storage kind")
  kind: CredentialStorageKind.Internal;

  @doc("The name of secret stored.")
  @visibility("read")
  secretName: string;
}
