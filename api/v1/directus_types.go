/*
Copyright 2025.

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

package v1

import (
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DirectusImage defines the container image configuration
type DirectusImage struct {
	// Repository is the container image repository
	Repository string `json:"repository,omitempty"`
	// Tag is the container image tag (defaults to appVersion if not set)
	Tag string `json:"tag,omitempty"`
	// PullPolicy defines the image pull policy
	PullPolicy corev1.PullPolicy `json:"pullPolicy,omitempty"`
}

// DirectusDatabase defines database configuration
type DirectusDatabase struct {
	// Engine defines the database engine (mysql or postgresql)
	Engine string `json:"engine,omitempty"`
	// Host is the database hostname
	Host string `json:"host,omitempty"`
	// Port is the database port
	Port int32 `json:"port,omitempty"`
	// Database is the database name
	Database string `json:"database,omitempty"`
	// Username is the database username
	Username string `json:"username,omitempty"`
	// ExistingSecret refers to an existing secret with database credentials
	ExistingSecret string `json:"existingSecret,omitempty"`
	// EnableInstallation determines if database should be installed (for managed databases)
	EnableInstallation bool `json:"enableInstallation,omitempty"`
}

// DirectusRedis defines Redis configuration
type DirectusRedis struct {
	// Enabled determines if Redis should be used
	Enabled bool `json:"enabled,omitempty"`
	// Host is the Redis hostname
	Host string `json:"host,omitempty"`
	// Port is the Redis port
	Port int32 `json:"port,omitempty"`
	// ExistingSecret refers to an existing secret with Redis credentials
	ExistingSecret string `json:"existingSecret,omitempty"`
	// EnableInstallation determines if Redis should be installed
	EnableInstallation bool `json:"enableInstallation,omitempty"`
}

// DirectusIngress defines ingress configuration
type DirectusIngress struct {
	// Enabled determines if ingress should be created
	Enabled bool `json:"enabled,omitempty"`
	// EnableTLS determines if TLS should be enabled in PUBLIC_URL
	EnableTLS bool `json:"enableTLS,omitempty"`
	// ClassName specifies the ingress class
	ClassName string `json:"className,omitempty"`
	// Annotations contains ingress annotations
	Annotations map[string]string `json:"annotations,omitempty"`
	// Hosts defines the ingress hosts
	Hosts []DirectusIngressHost `json:"hosts,omitempty"`
	// TLS defines the TLS configuration
	TLS []networkingv1.IngressTLS `json:"tls,omitempty"`
}

// DirectusIngressHost defines ingress host configuration
type DirectusIngressHost struct {
	// Host is the hostname
	Host string `json:"host,omitempty"`
	// Paths defines the paths for this host
	Paths []DirectusIngressPath `json:"paths,omitempty"`
}

// DirectusIngressPath defines ingress path configuration
type DirectusIngressPath struct {
	// Path is the URL path
	Path string `json:"path,omitempty"`
	// PathType defines the path type
	PathType string `json:"pathType,omitempty"`
}

// DirectusAutoscaling defines HPA configuration
type DirectusAutoscaling struct {
	// Enabled determines if HPA should be created
	Enabled bool `json:"enabled,omitempty"`
	// MinReplicas is the minimum number of replicas
	MinReplicas int32 `json:"minReplicas,omitempty"`
	// MaxReplicas is the maximum number of replicas
	MaxReplicas int32 `json:"maxReplicas,omitempty"`
	// TargetCPUUtilizationPercentage is the target CPU utilization
	TargetCPUUtilizationPercentage int32 `json:"targetCPUUtilizationPercentage,omitempty"`
	// TargetMemoryUtilizationPercentage is the target memory utilization
	TargetMemoryUtilizationPercentage int32 `json:"targetMemoryUtilizationPercentage,omitempty"`
}

// DirectusProbe defines probe configuration
type DirectusProbe struct {
	// Enabled determines if the probe should be enabled
	Enabled bool `json:"enabled,omitempty"`
	// HTTPGet defines the HTTP probe configuration
	HTTPGet *corev1.HTTPGetAction `json:"httpGet,omitempty"`
}

// DirectusServiceAccount defines service account configuration
type DirectusServiceAccount struct {
	// Create specifies whether a service account should be created
	Create bool `json:"create,omitempty"`
	// Annotations to add to the service account
	Annotations map[string]string `json:"annotations,omitempty"`
	// Name of the service account to use
	Name string `json:"name,omitempty"`
}

// DirectusService defines service configuration
type DirectusService struct {
	// Type defines the service type
	Type corev1.ServiceType `json:"type,omitempty"`
	// Port defines the service port
	Port int32 `json:"port,omitempty"`
}

// DirectusSpec defines the desired state of Directus.
type DirectusSpec struct {
	// ReplicaCount defines the number of Directus replicas
	ReplicaCount int32 `json:"replicaCount,omitempty"`

	// Image defines the container image configuration
	Image DirectusImage `json:"image,omitempty"`

	// ImagePullSecrets defines the image pull secrets
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`

	// AdminEmail defines the admin email address
	AdminEmail string `json:"adminEmail,omitempty"`

	// ServiceAccount defines service account configuration
	ServiceAccount DirectusServiceAccount `json:"serviceAccount,omitempty"`

	// PodAnnotations defines annotations to add to the pod
	PodAnnotations map[string]string `json:"podAnnotations,omitempty"`

	// PodSecurityContext defines the security context for the pod
	PodSecurityContext *corev1.PodSecurityContext `json:"podSecurityContext,omitempty"`

	// SecurityContext defines the security context for the container
	SecurityContext *corev1.SecurityContext `json:"securityContext,omitempty"`

	// Service defines the service configuration
	Service DirectusService `json:"service,omitempty"`

	// Ingress defines the ingress configuration
	Ingress DirectusIngress `json:"ingress,omitempty"`

	// ExtraEnvVars defines additional environment variables
	ExtraEnvVars []corev1.EnvVar `json:"extraEnvVars,omitempty"`

	// AttachExistingSecrets defines existing secrets to attach
	AttachExistingSecrets []string `json:"attachExistingSecrets,omitempty"`

	// Resources defines the resource requirements
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`

	// Autoscaling defines the HPA configuration
	Autoscaling DirectusAutoscaling `json:"autoscaling,omitempty"`

	// EnableLivenessProbe determines if liveness probe should be enabled
	EnableLivenessProbe bool `json:"enableLivenessProbe,omitempty"`

	// EnableReadinessProbe determines if readiness probe should be enabled
	EnableReadinessProbe bool `json:"enableReadinessProbe,omitempty"`

	// EnableStartupProbe determines if startup probe should be enabled
	EnableStartupProbe bool `json:"enableStartupProbe,omitempty"`

	// NodeSelector defines node selection constraints
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// Tolerations defines pod tolerations
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`

	// Affinity defines pod affinity rules
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// ExtraVolumes defines additional volumes
	ExtraVolumes []corev1.Volume `json:"extraVolumes,omitempty"`

	// ExtraVolumeMounts defines additional volume mounts
	ExtraVolumeMounts []corev1.VolumeMount `json:"extraVolumeMounts,omitempty"`

	// CreateApplicationSecret determines if application secrets should be created
	CreateApplicationSecret bool `json:"createApplicationSecret,omitempty"`

	// ApplicationSecretName defines the name of the application secret
	ApplicationSecretName string `json:"applicationSecretName,omitempty"`

	// Database defines the database configuration
	Database DirectusDatabase `json:"database,omitempty"`

	// Redis defines the Redis configuration
	Redis DirectusRedis `json:"redis,omitempty"`

	// InitContainers defines init containers
	InitContainers []corev1.Container `json:"initContainers,omitempty"`

	// Sidecars defines sidecar containers
	Sidecars []corev1.Container `json:"sidecars,omitempty"`
}

// DirectusStatus defines the observed state of Directus.
type DirectusStatus struct {
	// Conditions represent the latest available observations of the Directus state
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// ReadyReplicas indicates how many replicas are ready
	ReadyReplicas int32 `json:"readyReplicas,omitempty"`

	// Replicas indicates the number of replicas
	Replicas int32 `json:"replicas,omitempty"`

	// Phase indicates the current phase of the Directus deployment
	Phase string `json:"phase,omitempty"`

	// Message provides additional information about the current state
	Message string `json:"message,omitempty"`

	// DatabaseReady indicates if the database is ready
	DatabaseReady bool `json:"databaseReady,omitempty"`

	// RedisReady indicates if Redis is ready
	RedisReady bool `json:"redisReady,omitempty"`

	// IngressReady indicates if the ingress is ready
	IngressReady bool `json:"ingressReady,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:subresource:scale:specpath=.spec.replicaCount,statuspath=.status.replicas,selectorpath=.status.selector
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.readyReplicas"
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// Directus is the Schema for the directuses API.
type Directus struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DirectusSpec   `json:"spec,omitempty"`
	Status DirectusStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DirectusList contains a list of Directus.
type DirectusList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Directus `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Directus{}, &DirectusList{})
}
