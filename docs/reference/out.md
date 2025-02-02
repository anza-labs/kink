# API Reference

## Packages
- [controlplane.cluster.x-k8s.io/v1alpha1](#controlplaneclusterx-k8siov1alpha1)
- [infrastructure.cluster.x-k8s.io/v1alpha1](#infrastructureclusterx-k8siov1alpha1)


## controlplane.cluster.x-k8s.io/v1alpha1

Package v1alpha1 contains API Schema definitions for the controlplane v1alpha1 API group.

### Resource Types
- [KinkControlPlane](#kinkcontrolplane)
- [KinkControlPlaneList](#kinkcontrolplanelist)
- [KinkControlPlaneTemplate](#kinkcontrolplanetemplate)
- [KinkControlPlaneTemplateList](#kinkcontrolplanetemplatelist)



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
| `replicas` _integer_ | Number of desired pods. Defaults to 1. | 1 | Maximum: 5 <br />Minimum: 1 <br /> |
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
| `replicas` _integer_ | Number of desired pods. Defaults to 1. | 1 | Maximum: 5 <br />Minimum: 1 <br /> |
| `verbosity` _integer_ | Verbosity specifies the log verbosity level for the container. Valid values range from 0 (silent) to 10 (most verbose). | 4 | Maximum: 10 <br />Minimum: 0 <br /> |
| `extraArgs` _object (keys:string, values:string)_ | ExtraArgs defines additional arguments to be passed to the container executable. |  |  |


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
| `dnsName` _string_ | DNSName specifies the cluster endpoint, most likely backed by Ingress. | localhost |  |
| `imagePullSecrets` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#localobjectreference-v1-core) array_ | ImagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling any of the images used by this PodSpec.<br />If specified, these secrets will be passed to individual puller implementations for them to use. |  |  |
| `affinity` _[Affinity](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#affinity-v1-core)_ | Affinity specifies the scheduling constraints for Pods. |  |  |
| `apiServer` _[APIServer](#apiserver)_ | APIServer defines the configuration for the Kubernetes API server. |  |  |
| `kine` _[Kine](#kine)_ | Kine defines the configuration for the Kine component. |  |  |
| `scheduler` _[Scheduler](#scheduler)_ | Scheduler defines the configuration for the Kubernetes scheduler. |  |  |
| `controllerManager` _[ControllerManager](#controllermanager)_ | ControllerManager defines the configuration for the Kubernetes controller manager. |  |  |


#### KinkControlPlaneStatus



KinkControlPlaneStatus defines the observed state of KinkControlPlane.



_Appears in:_
- [KinkControlPlane](#kinkcontrolplane)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `version` _string_ | Version represents the minimum Kubernetes version for the control plane machines<br />in the cluster. |  |  |
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



#### KubeComponent



KubeComponent defines the base configuration for Kink control plane components.



_Appears in:_
- [APIServer](#apiserver)
- [ControllerManager](#controllermanager)
- [Scheduler](#scheduler)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `image` _string_ | Image specifies the container image to use. |  |  |
| `imagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#pullpolicy-v1-core)_ | Image pull policy. One of Always, Never, IfNotPresent. |  |  |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core)_ | Resources describes the compute resource requirements for the container. |  |  |
| `replicas` _integer_ | Number of desired pods. Defaults to 1. | 1 | Maximum: 5 <br />Minimum: 1 <br /> |
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
| `replicas` _integer_ | Number of desired pods. Defaults to 1. | 1 | Maximum: 5 <br />Minimum: 1 <br /> |
| `verbosity` _integer_ | Verbosity specifies the log verbosity level for the container. Valid values range from 0 (silent) to 10 (most verbose). | 4 | Maximum: 10 <br />Minimum: 0 <br /> |
| `extraArgs` _object (keys:string, values:string)_ | ExtraArgs defines additional arguments to be passed to the container executable. |  |  |



## infrastructure.cluster.x-k8s.io/v1alpha1

Package v1alpha1 contains API Schema definitions for the infrastructure v1alpha1 API group.

### Resource Types
- [KinkCluster](#kinkcluster)
- [KinkClusterList](#kinkclusterlist)
- [KinkClusterTemplate](#kinkclustertemplate)
- [KinkClusterTemplateList](#kinkclustertemplatelist)
- [KinkMachine](#kinkmachine)
- [KinkMachineList](#kinkmachinelist)
- [KinkMachineTemplate](#kinkmachinetemplate)
- [KinkMachineTemplateList](#kinkmachinetemplatelist)



#### KinkCluster



KinkCluster is the Schema for the kinkclusters API.



_Appears in:_
- [KinkClusterList](#kinkclusterlist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `infrastructure.cluster.x-k8s.io/v1alpha1` | | |
| `kind` _string_ | `KinkCluster` | | |
| `kind` _string_ | Kind is a string value representing the REST resource this object represents.<br />Servers may infer this from the endpoint the client submits requests to.<br />Cannot be updated.<br />In CamelCase.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds |  |  |
| `apiVersion` _string_ | APIVersion defines the versioned schema of this representation of an object.<br />Servers should convert recognized schemas to the latest internal value, and<br />may reject unrecognized values.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources |  |  |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[KinkClusterSpec](#kinkclusterspec)_ |  |  |  |
| `status` _[KinkClusterStatus](#kinkclusterstatus)_ |  |  |  |


#### KinkClusterList



KinkClusterList contains a list of KinkCluster.





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `infrastructure.cluster.x-k8s.io/v1alpha1` | | |
| `kind` _string_ | `KinkClusterList` | | |
| `kind` _string_ | Kind is a string value representing the REST resource this object represents.<br />Servers may infer this from the endpoint the client submits requests to.<br />Cannot be updated.<br />In CamelCase.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds |  |  |
| `apiVersion` _string_ | APIVersion defines the versioned schema of this representation of an object.<br />Servers should convert recognized schemas to the latest internal value, and<br />may reject unrecognized values.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources |  |  |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[KinkCluster](#kinkcluster) array_ |  |  |  |


#### KinkClusterSpec



KinkClusterSpec defines the desired state of KinkCluster.



_Appears in:_
- [KinkCluster](#kinkcluster)
- [KinkClusterTemplateResource](#kinkclustertemplateresource)



#### KinkClusterStatus



KinkClusterStatus defines the observed state of KinkCluster.



_Appears in:_
- [KinkCluster](#kinkcluster)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `ready` _boolean_ | Ready denotes that the kink cluster infrastructure is fully provisioned. |  |  |


#### KinkClusterTemplate



KinkClusterTemplate is the Schema for the kinkclustertemplates API.



_Appears in:_
- [KinkClusterTemplateList](#kinkclustertemplatelist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `infrastructure.cluster.x-k8s.io/v1alpha1` | | |
| `kind` _string_ | `KinkClusterTemplate` | | |
| `kind` _string_ | Kind is a string value representing the REST resource this object represents.<br />Servers may infer this from the endpoint the client submits requests to.<br />Cannot be updated.<br />In CamelCase.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds |  |  |
| `apiVersion` _string_ | APIVersion defines the versioned schema of this representation of an object.<br />Servers should convert recognized schemas to the latest internal value, and<br />may reject unrecognized values.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources |  |  |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[KinkClusterTemplateSpec](#kinkclustertemplatespec)_ |  |  |  |
| `status` _[KinkClusterTemplateStatus](#kinkclustertemplatestatus)_ |  |  |  |


#### KinkClusterTemplateList



KinkClusterTemplateList contains a list of KinkClusterTemplate.





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `infrastructure.cluster.x-k8s.io/v1alpha1` | | |
| `kind` _string_ | `KinkClusterTemplateList` | | |
| `kind` _string_ | Kind is a string value representing the REST resource this object represents.<br />Servers may infer this from the endpoint the client submits requests to.<br />Cannot be updated.<br />In CamelCase.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds |  |  |
| `apiVersion` _string_ | APIVersion defines the versioned schema of this representation of an object.<br />Servers should convert recognized schemas to the latest internal value, and<br />may reject unrecognized values.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources |  |  |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[KinkClusterTemplate](#kinkclustertemplate) array_ |  |  |  |


#### KinkClusterTemplateResource







_Appears in:_
- [KinkClusterTemplateSpec](#kinkclustertemplatespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[KinkClusterSpec](#kinkclusterspec)_ |  |  |  |


#### KinkClusterTemplateSpec



KinkClusterTemplateSpec defines the desired state of KinkClusterTemplate.



_Appears in:_
- [KinkClusterTemplate](#kinkclustertemplate)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `template` _[KinkClusterTemplateResource](#kinkclustertemplateresource)_ |  |  |  |


#### KinkClusterTemplateStatus



KinkClusterTemplateStatus defines the observed state of KinkClusterTemplate.



_Appears in:_
- [KinkClusterTemplate](#kinkclustertemplate)



#### KinkMachine



KinkMachine is the Schema for the kinkmachines API.



_Appears in:_
- [KinkMachineList](#kinkmachinelist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `infrastructure.cluster.x-k8s.io/v1alpha1` | | |
| `kind` _string_ | `KinkMachine` | | |
| `kind` _string_ | Kind is a string value representing the REST resource this object represents.<br />Servers may infer this from the endpoint the client submits requests to.<br />Cannot be updated.<br />In CamelCase.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds |  |  |
| `apiVersion` _string_ | APIVersion defines the versioned schema of this representation of an object.<br />Servers should convert recognized schemas to the latest internal value, and<br />may reject unrecognized values.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources |  |  |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[KinkMachineSpec](#kinkmachinespec)_ |  |  |  |
| `status` _[KinkMachineStatus](#kinkmachinestatus)_ |  |  |  |


#### KinkMachineList



KinkMachineList contains a list of KinkMachine.





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `infrastructure.cluster.x-k8s.io/v1alpha1` | | |
| `kind` _string_ | `KinkMachineList` | | |
| `kind` _string_ | Kind is a string value representing the REST resource this object represents.<br />Servers may infer this from the endpoint the client submits requests to.<br />Cannot be updated.<br />In CamelCase.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds |  |  |
| `apiVersion` _string_ | APIVersion defines the versioned schema of this representation of an object.<br />Servers should convert recognized schemas to the latest internal value, and<br />may reject unrecognized values.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources |  |  |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[KinkMachine](#kinkmachine) array_ |  |  |  |


#### KinkMachineSpec



KinkMachineSpec defines the desired state of KinkMachine.



_Appears in:_
- [KinkMachine](#kinkmachine)
- [KinkMachineTemplateResource](#kinkmachinetemplateresource)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `image` _string_ | Image specifies the container image to use. |  |  |
| `imagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#pullpolicy-v1-core)_ | Image pull policy. One of Always, Never, IfNotPresent. |  |  |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core)_ | Resources describes the compute resource requirements for the container. |  |  |
| `providerID` _string_ | ProviderID must match the provider ID as seen on the node object corresponding to this machine. |  |  |
| `imagePullSecrets` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#localobjectreference-v1-core) array_ | ImagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling any of the images used by this PodSpec.<br />If specified, these secrets will be passed to individual puller implementations for them to use. |  |  |
| `affinity` _[Affinity](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#affinity-v1-core)_ | Affinity specifies the scheduling constraints for Pods. |  |  |
| `persistence` _[Persistence](#persistence)_ | Persistence specifies volume configuration for Kine data persistence.<br />Defaults to EmptyDir. |  |  |


#### KinkMachineStatus



KinkMachineStatus defines the observed state of KinkMachine.



_Appears in:_
- [KinkMachine](#kinkmachine)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `ready` _boolean_ | Ready denotes that the kink machine infrastructure is fully provisioned. |  |  |


#### KinkMachineTemplate



KinkMachineTemplate is the Schema for the kinkmachinetemplates API.



_Appears in:_
- [KinkMachineTemplateList](#kinkmachinetemplatelist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `infrastructure.cluster.x-k8s.io/v1alpha1` | | |
| `kind` _string_ | `KinkMachineTemplate` | | |
| `kind` _string_ | Kind is a string value representing the REST resource this object represents.<br />Servers may infer this from the endpoint the client submits requests to.<br />Cannot be updated.<br />In CamelCase.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds |  |  |
| `apiVersion` _string_ | APIVersion defines the versioned schema of this representation of an object.<br />Servers should convert recognized schemas to the latest internal value, and<br />may reject unrecognized values.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources |  |  |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[KinkMachineTemplateSpec](#kinkmachinetemplatespec)_ |  |  |  |
| `status` _[KinkMachineTemplateStatus](#kinkmachinetemplatestatus)_ |  |  |  |


#### KinkMachineTemplateList



KinkMachineTemplateList contains a list of KinkMachineTemplate.





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `infrastructure.cluster.x-k8s.io/v1alpha1` | | |
| `kind` _string_ | `KinkMachineTemplateList` | | |
| `kind` _string_ | Kind is a string value representing the REST resource this object represents.<br />Servers may infer this from the endpoint the client submits requests to.<br />Cannot be updated.<br />In CamelCase.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds |  |  |
| `apiVersion` _string_ | APIVersion defines the versioned schema of this representation of an object.<br />Servers should convert recognized schemas to the latest internal value, and<br />may reject unrecognized values.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources |  |  |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[KinkMachineTemplate](#kinkmachinetemplate) array_ |  |  |  |


#### KinkMachineTemplateResource







_Appears in:_
- [KinkMachineTemplateSpec](#kinkmachinetemplatespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[KinkMachineSpec](#kinkmachinespec)_ |  |  |  |


#### KinkMachineTemplateSpec



KinkMachineTemplateSpec defines the desired state of KinkMachineTemplate.



_Appears in:_
- [KinkMachineTemplate](#kinkmachinetemplate)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `template` _[KinkMachineTemplateResource](#kinkmachinetemplateresource)_ |  |  |  |


#### KinkMachineTemplateStatus



KinkMachineTemplateStatus defines the observed state of KinkMachineTemplate.



_Appears in:_
- [KinkMachineTemplate](#kinkmachinetemplate)



