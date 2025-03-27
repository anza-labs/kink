# API Reference

## Packages
- [controlplane.cluster.x-k8s.io/v1alpha1](#controlplaneclusterx-k8siov1alpha1)


## controlplane.cluster.x-k8s.io/v1alpha1

Package v1alpha1 contains API Schema definitions for the controlplane v1alpha1 API group.

### Resource Types
- [KinkControlPlane](#kinkcontrolplane)
- [KinkControlPlaneList](#kinkcontrolplanelist)
- [KinkControlPlaneTemplate](#kinkcontrolplanetemplate)
- [KinkControlPlaneTemplateList](#kinkcontrolplanetemplatelist)



#### APIEndpoint



APIEndpoint represents a reachable Kubernetes API endpoint.



_Appears in:_
- [KinkControlPlaneSpec](#kinkcontrolplanespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `host` _[HostnameOrIP](#hostnameorip)_ | host is the hostname on which the API server is serving. |  |  |
| `port` _integer_ | port is the port on which the API server is serving. |  |  |


#### APIServer



APIServer represents a Kubernetes API server.


Image:
  - If specified image contains tag or sha, those are ignored.
  - Defaults to registry.k8s.io/kube-apiserver



_Appears in:_
- [KinkControlPlaneSpec](#kinkcontrolplanespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `image` _string_ | Image specifies the container image to use. |  |  |
| `imagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#pullpolicy-v1-core)_ | Image pull policy. One of Always, Never, IfNotPresent. |  |  |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core)_ | Resources describes the compute resource requirements for the container. |  |  |
| `verbosity` _integer_ | Verbosity specifies the log verbosity level for the container. Valid values range from 0 (silent) to 10 (most verbose). | 4 | Maximum: 10 <br />Minimum: 0 <br /> |
| `extraArgs` _object (keys:string, values:string)_ | ExtraArgs defines additional arguments to be passed to the container executable. |  |  |


#### ControllerManager



ControllerManager represents a Kubernetes controller manager.


Image:
  - If specified image contains tag or sha, those are ignored.
  - Defaults to registry.k8s.io/kube-controller-manager



_Appears in:_
- [KinkControlPlaneSpec](#kinkcontrolplanespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `image` _string_ | Image specifies the container image to use. |  |  |
| `imagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#pullpolicy-v1-core)_ | Image pull policy. One of Always, Never, IfNotPresent. |  |  |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core)_ | Resources describes the compute resource requirements for the container. |  |  |
| `verbosity` _integer_ | Verbosity specifies the log verbosity level for the container. Valid values range from 0 (silent) to 10 (most verbose). | 4 | Maximum: 10 <br />Minimum: 0 <br /> |
| `extraArgs` _object (keys:string, values:string)_ | ExtraArgs defines additional arguments to be passed to the container executable. |  |  |


#### EndpointsTemplate



EndpointsTemplate.



_Appears in:_
- [KinkControlPlaneSpec](#kinkcontrolplanespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `domain` _string_ | Domain. |  |  |
| `serviceType` _[ServiceType](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#servicetype-v1-core)_ | ServiceType | LoadBalancer | Enum: [LoadBalancer NodePort] <br /> |
| `gateway` _[Gateway](#gateway)_ | Gateway. |  |  |
| `ingress` _[Ingress](#ingress)_ | Ingress. |  |  |


#### Gateway



Gateway.



_Appears in:_
- [EndpointsTemplate](#endpointstemplate)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `gatewayClassName` _string_ | GatewayClassName used for this Gateway. This is the name of a<br />GatewayClass resource. |  |  |


#### HostnameOrIP

_Underlying type:_ _string_

HostnameOrIP.



_Appears in:_
- [APIEndpoint](#apiendpoint)



#### Ingress



Ingress.



_Appears in:_
- [EndpointsTemplate](#endpointstemplate)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `annotations` _object (keys:string, values:string)_ | Annotations is an unstructured key value map stored with a resource that may be<br />set by external tools to store and retrieve arbitrary metadata. They are not<br />queryable and should be preserved when modifying objects. |  |  |
| `ingressClassName` _string_ | GatewayClassName used for this Gateway. This is the name of a<br />GatewayClass resource. |  |  |


#### Kine



Kine represents ETCD-shim container.



_Appears in:_
- [KinkControlPlaneSpec](#kinkcontrolplanespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `image` _string_ | Image specifies the container image to use. |  |  |
| `imagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#pullpolicy-v1-core)_ | Image pull policy. One of Always, Never, IfNotPresent. |  |  |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core)_ | Resources describes the compute resource requirements for the container. |  |  |
| `persistence` _[Persistence](#persistence)_ | Persistence specifies volume configuration for Kine data persistence.<br />Defaults to EmptyDir. |  |  |


#### KinkControlPlane



KinkControlPlane is the Schema for the kinkcontrolplanes API.



_Appears in:_
- [KinkControlPlaneList](#kinkcontrolplanelist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `controlplane.cluster.x-k8s.io/v1alpha1` | | |
| `kind` _string_ | `KinkControlPlane` | | |
| `kind` _string_ | Kind is a string value representing the REST resource this object represents.<br />Servers may infer this from the endpoint the client submits requests to.<br />Cannot be updated.<br />In CamelCase.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds |  |  |
| `apiVersion` _string_ | APIVersion defines the versioned schema of this representation of an object.<br />Servers should convert recognized schemas to the latest internal value, and<br />may reject unrecognized values.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources |  |  |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[KinkControlPlaneSpec](#kinkcontrolplanespec)_ |  |  |  |
| `status` _[KinkControlPlaneStatus](#kinkcontrolplanestatus)_ |  |  |  |


#### KinkControlPlaneList



KinkControlPlaneList contains a list of KinkControlPlane.





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `controlplane.cluster.x-k8s.io/v1alpha1` | | |
| `kind` _string_ | `KinkControlPlaneList` | | |
| `kind` _string_ | Kind is a string value representing the REST resource this object represents.<br />Servers may infer this from the endpoint the client submits requests to.<br />Cannot be updated.<br />In CamelCase.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds |  |  |
| `apiVersion` _string_ | APIVersion defines the versioned schema of this representation of an object.<br />Servers should convert recognized schemas to the latest internal value, and<br />may reject unrecognized values.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources |  |  |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[KinkControlPlane](#kinkcontrolplane) array_ |  |  |  |


#### KinkControlPlaneSpec



KinkControlPlaneSpec defines the desired state of KinkControlPlane.



_Appears in:_
- [KinkControlPlane](#kinkcontrolplane)
- [KinkControlPlaneTemplateResource](#kinkcontrolplanetemplateresource)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `version` _string_ | Version defines the desired Kubernetes version for the control plane.<br />The value must be a valid semantic version; also if the value provided by the user<br />does not start with the v prefix, it must be added. |  |  |
| `endpointsTemplate` _[EndpointsTemplate](#endpointstemplate)_ | EndpointsTemplate. |  |  |
| `controlPlaneEndpoint` _[APIEndpoint](#apiendpoint)_ | ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.<br />Based on it, an ingress will be provisioned.<br />This field is set by the controller. |  |  |
| `replicas` _integer_ | Number of desired ControlPlane replicas. Defaults to 1. | 1 | Maximum: 5 <br />Minimum: 1 <br /> |
| `imagePullSecrets` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#localobjectreference-v1-core) array_ | ImagePullSecrets is an optional list of references to secrets in the same namespace to use<br />for pulling any of the images used by KinkControlPlane. If specified, these secrets will<br />be passed to individual puller implementations for them to use. |  |  |
| `affinity` _[Affinity](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#affinity-v1-core)_ | Affinity specifies the scheduling constraints for Pods. |  |  |
| `apiServer` _[APIServer](#apiserver)_ | APIServer defines the configuration for the Kubernetes API server. |  |  |
| `konnectivityServer` _[KonnectivityServer](#konnectivityserver)_ | KonnectivityServer defines the configuration for the API Proxy server. |  |  |
| `konnectivityAgent` _[KonnectivityAgent](#konnectivityagent)_ | KonnectivityAgent defines the configuration for the API Proxy agent. |  |  |
| `kine` _[Kine](#kine)_ | Kine defines the configuration for the Kine component. |  |  |
| `scheduler` _[Scheduler](#scheduler)_ | Scheduler defines the configuration for the Kubernetes scheduler. |  |  |
| `controllerManager` _[ControllerManager](#controllermanager)_ | ControllerManager defines the configuration for the Kubernetes controller manager. |  |  |


#### KinkControlPlaneStatus



KinkControlPlaneStatus defines the observed state of KinkControlPlane.



_Appears in:_
- [KinkControlPlane](#kinkcontrolplane)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `version` _string_ | Version represents the minimum Kubernetes version for the control plane replicas<br />in the cluster. |  |  |
| `selector` _string_ | Selector is the label selector in string format to avoid introspection<br />by clients, and is used to provide the CRD-based integration for the<br />scale subresource and additional integrations for things like kubectl<br />describe. The string will be in the same format as the query-param syntax.<br />More info about label selectors: http://kubernetes.io/docs/user-guide/labels#label-selectors |  |  |
| `replicas` _integer_ | Replicas is the total number of replicas targeted by this control plane<br />(their labels match the selector). |  |  |
| `updatedReplicas` _integer_ | UpdatedReplicas is the total number of replicas targeted by this control plane<br />that have the desired template spec. |  |  |
| `readyReplicas` _integer_ | ReadyReplicas is the total number of fully running and ready control plane replicas. |  |  |
| `unavailableReplicas` _integer_ | UnavailableReplicas is the total number of unavailable replicas targeted by this control plane.<br />This is the total number of replicas that are still required for the deployment to have 100% available capacity.<br />They may either be replicas that are running but not yet ready or replicas<br />that still have not been created. |  |  |
| `conditions` _[Condition](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#condition-v1-meta) array_ | Conditions defines current service state of the KinkControlPlane. |  |  |
| `initialized` _boolean_ | Initialized denotes that the kink control plane API Server is initialized and thus<br />it can accept requests. |  |  |
| `ready` _boolean_ | Ready denotes that the kink control plane is ready to serve requests. |  |  |


#### KinkControlPlaneTemplate



KinkControlPlaneTemplate is the Schema for the kinkcontrolplanetemplates API.



_Appears in:_
- [KinkControlPlaneTemplateList](#kinkcontrolplanetemplatelist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `controlplane.cluster.x-k8s.io/v1alpha1` | | |
| `kind` _string_ | `KinkControlPlaneTemplate` | | |
| `kind` _string_ | Kind is a string value representing the REST resource this object represents.<br />Servers may infer this from the endpoint the client submits requests to.<br />Cannot be updated.<br />In CamelCase.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds |  |  |
| `apiVersion` _string_ | APIVersion defines the versioned schema of this representation of an object.<br />Servers should convert recognized schemas to the latest internal value, and<br />may reject unrecognized values.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources |  |  |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[KinkControlPlaneTemplateSpec](#kinkcontrolplanetemplatespec)_ |  |  |  |
| `status` _[KinkControlPlaneTemplateStatus](#kinkcontrolplanetemplatestatus)_ |  |  |  |


#### KinkControlPlaneTemplateList



KinkControlPlaneTemplateList contains a list of KinkControlPlaneTemplate.





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `controlplane.cluster.x-k8s.io/v1alpha1` | | |
| `kind` _string_ | `KinkControlPlaneTemplateList` | | |
| `kind` _string_ | Kind is a string value representing the REST resource this object represents.<br />Servers may infer this from the endpoint the client submits requests to.<br />Cannot be updated.<br />In CamelCase.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds |  |  |
| `apiVersion` _string_ | APIVersion defines the versioned schema of this representation of an object.<br />Servers should convert recognized schemas to the latest internal value, and<br />may reject unrecognized values.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources |  |  |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[KinkControlPlaneTemplate](#kinkcontrolplanetemplate) array_ |  |  |  |


#### KinkControlPlaneTemplateResource







_Appears in:_
- [KinkControlPlaneTemplateSpec](#kinkcontrolplanetemplatespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[KinkControlPlaneSpec](#kinkcontrolplanespec)_ |  |  |  |


#### KinkControlPlaneTemplateSpec



KinkControlPlaneTemplateSpec defines the desired state of KinkControlPlaneTemplate.



_Appears in:_
- [KinkControlPlaneTemplate](#kinkcontrolplanetemplate)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `template` _[KinkControlPlaneTemplateResource](#kinkcontrolplanetemplateresource)_ |  |  |  |


#### KinkControlPlaneTemplateStatus



KinkControlPlaneTemplateStatus defines the observed state of KinkControlPlaneTemplate.



_Appears in:_
- [KinkControlPlaneTemplate](#kinkcontrolplanetemplate)



#### KonnectivityAgent



KonnectivityAgent represents a API Proxy agent.


Image:
  - Defaults to registry.k8s.io/kas-network-proxy/proxy-agent



_Appears in:_
- [KinkControlPlaneSpec](#kinkcontrolplanespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `image` _string_ | Image specifies the container image to use. |  |  |
| `imagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#pullpolicy-v1-core)_ | Image pull policy. One of Always, Never, IfNotPresent. |  |  |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core)_ | Resources describes the compute resource requirements for the container. |  |  |
| `verbosity` _integer_ | Verbosity specifies the log verbosity level for the container. Valid values range from 0 (silent) to 10 (most verbose). | 4 | Maximum: 10 <br />Minimum: 0 <br /> |
| `extraArgs` _object (keys:string, values:string)_ | ExtraArgs defines additional arguments to be passed to the container executable. |  |  |


#### KonnectivityServer



KonnectivityServer represents a API Proxy server.


Image:
  - Defaults to registry.k8s.io/kas-network-proxy/proxy-server



_Appears in:_
- [KinkControlPlaneSpec](#kinkcontrolplanespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `image` _string_ | Image specifies the container image to use. |  |  |
| `imagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#pullpolicy-v1-core)_ | Image pull policy. One of Always, Never, IfNotPresent. |  |  |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core)_ | Resources describes the compute resource requirements for the container. |  |  |
| `verbosity` _integer_ | Verbosity specifies the log verbosity level for the container. Valid values range from 0 (silent) to 10 (most verbose). | 4 | Maximum: 10 <br />Minimum: 0 <br /> |
| `extraArgs` _object (keys:string, values:string)_ | ExtraArgs defines additional arguments to be passed to the container executable. |  |  |


#### KubeComponent



KubeComponent defines the base configuration for Kink control plane components.



_Appears in:_
- [APIServer](#apiserver)
- [ControllerManager](#controllermanager)
- [KonnectivityAgent](#konnectivityagent)
- [KonnectivityServer](#konnectivityserver)
- [Scheduler](#scheduler)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `image` _string_ | Image specifies the container image to use. |  |  |
| `imagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#pullpolicy-v1-core)_ | Image pull policy. One of Always, Never, IfNotPresent. |  |  |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core)_ | Resources describes the compute resource requirements for the container. |  |  |
| `verbosity` _integer_ | Verbosity specifies the log verbosity level for the container. Valid values range from 0 (silent) to 10 (most verbose). | 4 | Maximum: 10 <br />Minimum: 0 <br /> |
| `extraArgs` _object (keys:string, values:string)_ | ExtraArgs defines additional arguments to be passed to the container executable. |  |  |


#### Scheduler



Scheduler represents a Kubernetes scheduler.


Image:
  - If specified image contains tag or sha, those are ignored.
  - Defaults to registry.k8s.io/kube-scheduler



_Appears in:_
- [KinkControlPlaneSpec](#kinkcontrolplanespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `image` _string_ | Image specifies the container image to use. |  |  |
| `imagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#pullpolicy-v1-core)_ | Image pull policy. One of Always, Never, IfNotPresent. |  |  |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core)_ | Resources describes the compute resource requirements for the container. |  |  |
| `verbosity` _integer_ | Verbosity specifies the log verbosity level for the container. Valid values range from 0 (silent) to 10 (most verbose). | 4 | Maximum: 10 <br />Minimum: 0 <br /> |
| `extraArgs` _object (keys:string, values:string)_ | ExtraArgs defines additional arguments to be passed to the container executable. |  |  |


