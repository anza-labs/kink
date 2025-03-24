// Copyright 2025 anza-labs contributors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controlplane

import (
	"fmt"
	"maps"
	"path"

	controlplanev1alpha1 "github.com/anza-labs/kink/api/controlplane/v1alpha1"
	"github.com/anza-labs/kink/internal/manifests/manifestutils"
	"github.com/anza-labs/kink/internal/naming"
	"github.com/anza-labs/kink/version"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gatewayapiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

const (
	apiServerPKIPath         = "/etc/pki/kube-apiserver"
	apiServerCertificateFile = "tls.crt"
	apiServerKeyFile         = "tls.key"

	etcdPKIPath         = "/etc/pki/etcd"
	etcdCAFile          = "ca.crt"
	etcdCertificateFile = "tls.crt"
	etcdKeyFile         = "tls.key"
)

type APIServer struct {
	KinkControlPlane *controlplanev1alpha1.KinkControlPlane
}

func (b *APIServer) Build() ([]client.Object, error) {
	objects := []client.Object{}

	svc, err := b.Service()
	if err != nil {
		return nil, fmt.Errorf("failed to build Service: %w", err)
	}
	objects = append(objects, svc)

	depl, err := b.Deployment()
	if err != nil {
		return nil, fmt.Errorf("failed to build Deployment: %w", err)
	}
	objects = append(objects, depl)

	if b.KinkControlPlane.Spec.ControlPlaneEndpoint.Gateway != nil {
		gtw, err := b.Gateway()
		if err != nil {
			return nil, fmt.Errorf("failed to build Gateway: %w", err)
		}
		objects = append(objects, gtw)

		rte, err := b.HTTPRoute()
		if err != nil {
			return nil, fmt.Errorf("failed to build HTTPRoute: %w", err)
		}
		objects = append(objects, rte)
	}

	if b.KinkControlPlane.Spec.ControlPlaneEndpoint.Ingress != nil {
		ing, err := b.Ingress()
		if err != nil {
			return nil, fmt.Errorf("failed to build Ingress: %w", err)
		}
		objects = append(objects, ing)
	}

	return objects, nil
}

func (b *APIServer) Deployment() (*appsv1.Deployment, error) {
	name := naming.APIServer(b.KinkControlPlane.Name)

	image, err := manifestutils.Image(
		b.KinkControlPlane.Spec.APIServer.Image,
		b.KinkControlPlane.Spec.Version,
		version.APIServer(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to assess image: %w", err)
	}

	labels := manifestutils.Labels(
		b.KinkControlPlane.ObjectMeta,
		name, image, ComponentAPIServer, ConceptControlPlane,
		nil,
	)
	selectorLabels := manifestutils.SelectorLabels(b.KinkControlPlane.ObjectMeta, ComponentAPIServer, ConceptControlPlane)
	annotations := manifestutils.Annotations(b.KinkControlPlane, nil)
	podAnnotations := manifestutils.PodAnnotations(b.KinkControlPlane, nil)

	replicas := b.KinkControlPlane.Spec.Replicas
	if replicas == nil {
		replicas = ptr.To[int32](1)
	}

	podSpec := corev1.PodSpec{
		Affinity:         manifestutils.Affinity(b.KinkControlPlane),
		Containers:       []corev1.Container{b.container(image)},
		Volumes:          b.volumes(),
		ImagePullSecrets: b.KinkControlPlane.Spec.ImagePullSecrets,
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   b.KinkControlPlane.Namespace,
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: selectorLabels,
			},
			Replicas: replicas,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: podAnnotations,
				},
				Spec: podSpec,
			},
		},
	}, nil
}

func (b *APIServer) Service() (*corev1.Service, error) {
	name := naming.APIServer(b.KinkControlPlane.Name)

	image, err := manifestutils.Image(
		b.KinkControlPlane.Spec.APIServer.Image,
		b.KinkControlPlane.Spec.Version,
		version.APIServer(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to assess image: %w", err)
	}

	labels := manifestutils.Labels(
		b.KinkControlPlane.ObjectMeta,
		name, image, ComponentAPIServer, ConceptControlPlane,
		nil,
	)
	selectorLabels := manifestutils.SelectorLabels(b.KinkControlPlane.ObjectMeta, ComponentAPIServer, ConceptControlPlane)
	annotations := manifestutils.Annotations(b.KinkControlPlane, nil)

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   b.KinkControlPlane.Namespace,
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: corev1.ServiceSpec{
			Selector: selectorLabels,
			Type:     b.KinkControlPlane.Spec.ControlPlaneEndpoint.ServiceType,
			Ports: []corev1.ServicePort{
				{
					Name:       "server",
					Port:       6443,
					TargetPort: intstr.FromString("server"),
					Protocol:   corev1.ProtocolTCP,
				},
			},
		},
	}, nil
}

func (b *APIServer) Gateway() (*gatewayapiv1.Gateway, error) {
	name := naming.APIServer(b.KinkControlPlane.Name)
	certName := naming.APIServerCertificate(b.KinkControlPlane.Name)

	image, err := manifestutils.Image(
		b.KinkControlPlane.Spec.APIServer.Image,
		b.KinkControlPlane.Spec.Version,
		version.APIServer(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to assess image: %w", err)
	}

	labels := manifestutils.Labels(
		b.KinkControlPlane.ObjectMeta,
		name, image, ComponentAPIServer, ConceptControlPlane,
		nil,
	)
	annotations := manifestutils.Annotations(b.KinkControlPlane, nil)

	host := gatewayapiv1.Hostname(b.KinkControlPlane.Spec.ControlPlaneEndpoint.Host)
	port := gatewayapiv1.PortNumber(b.KinkControlPlane.Spec.ControlPlaneEndpoint.Port)
	gatewayClassName := gatewayapiv1.ObjectName(b.KinkControlPlane.Spec.ControlPlaneEndpoint.Gateway.GatewayClassName)

	return &gatewayapiv1.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   b.KinkControlPlane.Namespace,
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: gatewayapiv1.GatewaySpec{
			GatewayClassName: gatewayClassName,
			Listeners: []gatewayapiv1.Listener{
				{
					Name:     gatewayapiv1.SectionName(name),
					Hostname: &host,
					Port:     port,
					Protocol: gatewayapiv1.TLSProtocolType,
					TLS: &gatewayapiv1.GatewayTLSConfig{
						Mode: ptr.To(gatewayapiv1.TLSModePassthrough),
						CertificateRefs: []gatewayapiv1.SecretObjectReference{
							{
								Name: gatewayapiv1.ObjectName(certName),
							},
						},
					},
				},
			},
		},
	}, nil
}

func (b *APIServer) HTTPRoute() (*gatewayapiv1.HTTPRoute, error) {
	name := naming.APIServer(b.KinkControlPlane.Name)

	image, err := manifestutils.Image(
		b.KinkControlPlane.Spec.APIServer.Image,
		b.KinkControlPlane.Spec.Version,
		version.APIServer(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to assess image: %w", err)
	}

	labels := manifestutils.Labels(
		b.KinkControlPlane.ObjectMeta,
		name, image, ComponentAPIServer, ConceptControlPlane,
		nil,
	)
	annotations := manifestutils.Annotations(b.KinkControlPlane, nil)

	host := gatewayapiv1.Hostname(b.KinkControlPlane.Spec.ControlPlaneEndpoint.Host)

	return &gatewayapiv1.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   b.KinkControlPlane.Namespace,
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: gatewayapiv1.HTTPRouteSpec{
			CommonRouteSpec: gatewayapiv1.CommonRouteSpec{
				ParentRefs: []gatewayapiv1.ParentReference{
					{
						Name: gatewayapiv1.ObjectName(name),
					},
				},
			},
			Hostnames: []gatewayapiv1.Hostname{host},
			Rules: []gatewayapiv1.HTTPRouteRule{
				{
					BackendRefs: []gatewayapiv1.HTTPBackendRef{
						{
							BackendRef: gatewayapiv1.BackendRef{
								BackendObjectReference: gatewayapiv1.BackendObjectReference{
									Name: gatewayapiv1.ObjectName(name),
									Port: ptr.To(gatewayapiv1.PortNumber(6443)),
								},
							},
						},
					},
				},
			},
		},
	}, nil
}

func (b *APIServer) Ingress() (*netv1.Ingress, error) {
	name := naming.APIServer(b.KinkControlPlane.Name)

	image, err := manifestutils.Image(
		b.KinkControlPlane.Spec.APIServer.Image,
		b.KinkControlPlane.Spec.Version,
		version.APIServer(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to assess image: %w", err)
	}

	labels := manifestutils.Labels(
		b.KinkControlPlane.ObjectMeta,
		name, image, ComponentAPIServer, ConceptControlPlane,
		nil,
	)
	annotations := manifestutils.Annotations(b.KinkControlPlane, nil)
	maps.Insert(annotations, maps.All(b.KinkControlPlane.Spec.ControlPlaneEndpoint.Ingress.Annotations))

	return &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   b.KinkControlPlane.Namespace,
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: netv1.IngressSpec{
			IngressClassName: &b.KinkControlPlane.Spec.ControlPlaneEndpoint.Ingress.IngressClassName,
			TLS: []netv1.IngressTLS{
				{
					Hosts: []string{
						string(b.KinkControlPlane.Spec.ControlPlaneEndpoint.Host),
					},
					SecretName: naming.APIServerCertificate(b.KinkControlPlane.Name),
				},
			},
			Rules: []netv1.IngressRule{
				{
					Host: string(b.KinkControlPlane.Spec.ControlPlaneEndpoint.Host),
					IngressRuleValue: netv1.IngressRuleValue{
						HTTP: &netv1.HTTPIngressRuleValue{
							Paths: []netv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: ptr.To(netv1.PathTypePrefix),
									Backend: netv1.IngressBackend{
										Service: &netv1.IngressServiceBackend{
											Name: name,
											Port: netv1.ServiceBackendPort{
												Name: "server",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}, nil
}

func (b *APIServer) volumes() []corev1.Volume {
	return []corev1.Volume{
		{
			Name: "etcd",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  naming.KineAPIServerClientCertificate(b.KinkControlPlane.Name),
					DefaultMode: ptr.To[int32](420),
				},
			},
		},
		{
			Name: "apiserver-tls",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  naming.APIServerCertificate(b.KinkControlPlane.Name),
					DefaultMode: ptr.To[int32](420),
				},
			},
		},
		{
			Name: "root-ca",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  naming.ClusterCA(b.KinkControlPlane.Name),
					DefaultMode: ptr.To[int32](420),
				},
			},
		},
		{
			Name: "service-accounts-cert",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  naming.ServiceAccountCertificate(b.KinkControlPlane.Name),
					DefaultMode: ptr.To[int32](420),
				},
			},
		},
	}
}

func (b *APIServer) volumeMounts() []corev1.VolumeMount {
	return []corev1.VolumeMount{
		{
			Name:      "etcd",
			ReadOnly:  true,
			MountPath: etcdPKIPath,
		},
		{
			Name:      "apiserver-tls",
			ReadOnly:  true,
			MountPath: apiServerPKIPath,
		},
		{
			Name:      "root-ca",
			ReadOnly:  true,
			MountPath: rootPKIPath,
		},
		{
			Name:      "service-accounts-cert",
			ReadOnly:  true,
			MountPath: serviceAccountsPKIPath,
		},
	}
}

func (b *APIServer) container(image string) corev1.Container {
	cfg := b.KinkControlPlane.Spec.APIServer
	resources := cfg.Resources
	verbosity := cfg.Verbosity

	args := map[string]string{ // TODO: add default flags
		"v":                                fmt.Sprint(verbosity),
		"client-ca-file":                   path.Join(rootPKIPath, rootCAFile),
		"tls-cert-file":                    path.Join(apiServerPKIPath, apiServerCertificateFile),
		"tls-private-key-file":             path.Join(apiServerPKIPath, apiServerKeyFile),
		"service-account-key-file":         path.Join(serviceAccountsPKIPath, serviceAccountsCertificateFile),
		"service-account-signing-key-file": path.Join(serviceAccountsPKIPath, serviceAccountsKeyFile),
		"service-account-issuer":           "https://kubernetes.default.svc.cluster.local",
		"etcd-cafile":                      path.Join(etcdPKIPath, etcdCAFile),
		"etcd-certfile":                    path.Join(etcdPKIPath, etcdCertificateFile),
		"etcd-keyfile":                     path.Join(etcdPKIPath, etcdKeyFile),
		"etcd-servers":                     naming.KineEndpoint(b.KinkControlPlane.Name, b.KinkControlPlane.Namespace),
		"authorization-mode":               "Node,RBAC",
		"service-cluster-ip-range":         "10.32.0.0/24",
	}
	for arg, value := range cfg.ExtraArgs {
		if _, ok := args[arg]; !ok {
			args[arg] = value
		}
	}

	return corev1.Container{
		Name:      naming.APIServerContainer(),
		Image:     image,
		Command:   []string{"kube-apiserver"},
		Args:      buildArgs(args),
		Resources: resources,
		Ports: []corev1.ContainerPort{
			{
				Name:          "server",
				ContainerPort: 6443,
				Protocol:      corev1.ProtocolTCP,
			},
		},
		VolumeMounts:    b.volumeMounts(),
		ImagePullPolicy: b.KinkControlPlane.Spec.APIServer.ImagePullPolicy,
		SecurityContext: nil,
		StartupProbe: &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path:   "/livez",
					Port:   intstr.FromString("server"),
					Scheme: corev1.URISchemeHTTPS,
				},
			},
			InitialDelaySeconds: 10,
		},
		LivenessProbe: &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path:   "/livez",
					Port:   intstr.FromString("server"),
					Scheme: corev1.URISchemeHTTPS,
				},
			},
			InitialDelaySeconds: 10,
		},
		ReadinessProbe: &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path:   "/readyz",
					Port:   intstr.FromString("server"),
					Scheme: corev1.URISchemeHTTPS,
				},
			},
			InitialDelaySeconds: 10,
		},
	}
}
