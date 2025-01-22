# API Reference

## Packages
- [infrastructure.cluster.x-k8s.io/v1alpha1](#infrastructureclusterx-k8siov1alpha1)


## infrastructure.cluster.x-k8s.io/v1alpha1

Package v1alpha1 contains API Schema definitions for the infrastructure v1alpha1 API group.

### Resource Types
- [KinkCluster](#kinkcluster)
- [KinkClusterList](#kinkclusterlist)
- [KinkMachine](#kinkmachine)
- [KinkMachineList](#kinkmachinelist)



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

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `foo` _string_ | Foo is an example field of KinkCluster. Edit kinkcluster_types.go to remove/update |  |  |


#### KinkClusterStatus



KinkClusterStatus defines the observed state of KinkCluster.



_Appears in:_
- [KinkCluster](#kinkcluster)



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

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `foo` _string_ | Foo is an example field of KinkMachine. Edit kinkmachine_types.go to remove/update |  |  |


#### KinkMachineStatus



KinkMachineStatus defines the observed state of KinkMachine.



_Appears in:_
- [KinkMachine](#kinkmachine)



