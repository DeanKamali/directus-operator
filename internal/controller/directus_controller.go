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

package controller

import (
	"context"
	"fmt"
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	directusv1 "github.com/example/directus-operator/api/v1"
)

// DirectusReconciler reconciles a Directus object
type DirectusReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=directus.example.com,resources=directuses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=directus.example.com,resources=directuses/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=directus.example.com,resources=directuses/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=serviceaccounts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=autoscaling,resources=horizontalpodautoscalers,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *DirectusReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// Fetch the Directus instance
	var directus directusv1.Directus
	if err := r.Get(ctx, req.NamespacedName, &directus); err != nil {
		if errors.IsNotFound(err) {
			log.Info("Directus resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get Directus")
		return ctrl.Result{}, err
	}

	// Apply defaults if not specified
	if directus.Spec.ReplicaCount == 0 {
		directus.Spec.ReplicaCount = 1
	}
	if directus.Spec.Image.Repository == "" {
		directus.Spec.Image.Repository = "directus/directus"
	}
	if directus.Spec.Image.Tag == "" {
		directus.Spec.Image.Tag = "latest"
	}
	if directus.Spec.Service.Port == 0 {
		directus.Spec.Service.Port = 80
	}
	if directus.Spec.Service.Type == "" {
		directus.Spec.Service.Type = corev1.ServiceTypeClusterIP
	}
	if directus.Spec.AdminEmail == "" {
		directus.Spec.AdminEmail = "directus-admin@example.com"
	}

	// Create or update resources
	if err := r.reconcileServiceAccount(ctx, &directus); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.reconcileSecrets(ctx, &directus); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.reconcileConfigMap(ctx, &directus); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.reconcileService(ctx, &directus); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.reconcileDeployment(ctx, &directus); err != nil {
		return ctrl.Result{}, err
	}

	if directus.Spec.Ingress.Enabled {
		if err := r.reconcileIngress(ctx, &directus); err != nil {
			return ctrl.Result{}, err
		}
	}

	if directus.Spec.Autoscaling.Enabled {
		if err := r.reconcileHPA(ctx, &directus); err != nil {
			return ctrl.Result{}, err
		}
	}

	// Update status
	if err := r.updateStatus(ctx, &directus); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *DirectusReconciler) reconcileServiceAccount(ctx context.Context, directus *directusv1.Directus) error {
	if !directus.Spec.ServiceAccount.Create {
		return nil
	}

	sa := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:        r.getServiceAccountName(directus),
			Namespace:   directus.Namespace,
			Annotations: directus.Spec.ServiceAccount.Annotations,
		},
	}

	if err := controllerutil.SetControllerReference(directus, sa, r.Scheme); err != nil {
		return err
	}

	found := &corev1.ServiceAccount{}
	err := r.Get(ctx, types.NamespacedName{Name: sa.Name, Namespace: sa.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		return r.Create(ctx, sa)
	} else if err != nil {
		return err
	}

	// Update if needed
	found.Annotations = sa.Annotations
	return r.Update(ctx, found)
}

func (r *DirectusReconciler) reconcileSecrets(ctx context.Context, directus *directusv1.Directus) error {
	if !directus.Spec.CreateApplicationSecret {
		return nil
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.getApplicationSecretName(directus),
			Namespace: directus.Namespace,
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			"ADMIN_PASSWORD": []byte("admin123"),    // In production, generate random password
			"KEY":            []byte("your-key"),    // In production, generate random key
			"SECRET":         []byte("your-secret"), // In production, generate random secret
		},
	}

	if err := controllerutil.SetControllerReference(directus, secret, r.Scheme); err != nil {
		return err
	}

	found := &corev1.Secret{}
	err := r.Get(ctx, types.NamespacedName{Name: secret.Name, Namespace: secret.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		return r.Create(ctx, secret)
	} else if err != nil {
		return err
	}

	return nil // Don't update existing secrets to avoid overwriting passwords
}

func (r *DirectusReconciler) reconcileConfigMap(ctx context.Context, directus *directusv1.Directus) error {
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      directus.Name + "-configmap",
			Namespace: directus.Namespace,
		},
		Data: r.buildConfigMapData(directus),
	}

	if err := controllerutil.SetControllerReference(directus, configMap, r.Scheme); err != nil {
		return err
	}

	found := &corev1.ConfigMap{}
	err := r.Get(ctx, types.NamespacedName{Name: configMap.Name, Namespace: configMap.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		return r.Create(ctx, configMap)
	} else if err != nil {
		return err
	}

	// Update data
	found.Data = configMap.Data
	return r.Update(ctx, found)
}

func (r *DirectusReconciler) reconcileService(ctx context.Context, directus *directusv1.Directus) error {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      directus.Name,
			Namespace: directus.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Type: directus.Spec.Service.Type,
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Port:       directus.Spec.Service.Port,
					TargetPort: intstr.FromInt(8055),
					Protocol:   corev1.ProtocolTCP,
				},
			},
			Selector: r.getLabels(directus),
		},
	}

	if err := controllerutil.SetControllerReference(directus, service, r.Scheme); err != nil {
		return err
	}

	found := &corev1.Service{}
	err := r.Get(ctx, types.NamespacedName{Name: service.Name, Namespace: service.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		return r.Create(ctx, service)
	} else if err != nil {
		return err
	}

	// Update spec (excluding ClusterIP which is immutable)
	found.Spec.Type = service.Spec.Type
	found.Spec.Ports = service.Spec.Ports
	found.Spec.Selector = service.Spec.Selector
	return r.Update(ctx, found)
}

func (r *DirectusReconciler) reconcileDeployment(ctx context.Context, directus *directusv1.Directus) error {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      directus.Name,
			Namespace: directus.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &directus.Spec.ReplicaCount,
			Selector: &metav1.LabelSelector{
				MatchLabels: r.getLabels(directus),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      r.getLabels(directus),
					Annotations: directus.Spec.PodAnnotations,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: r.getServiceAccountName(directus),
					SecurityContext:    directus.Spec.PodSecurityContext,
					ImagePullSecrets:   directus.Spec.ImagePullSecrets,
					InitContainers:     directus.Spec.InitContainers,
					Containers:         r.buildContainers(directus),
					Volumes:            directus.Spec.ExtraVolumes,
					NodeSelector:       directus.Spec.NodeSelector,
					Tolerations:        directus.Spec.Tolerations,
					Affinity:           directus.Spec.Affinity,
				},
			},
		},
	}

	// Add sidecar containers
	deployment.Spec.Template.Spec.Containers = append(deployment.Spec.Template.Spec.Containers, directus.Spec.Sidecars...)

	if err := controllerutil.SetControllerReference(directus, deployment, r.Scheme); err != nil {
		return err
	}

	found := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		return r.Create(ctx, deployment)
	} else if err != nil {
		return err
	}

	// Update deployment
	found.Spec = deployment.Spec
	return r.Update(ctx, found)
}

func (r *DirectusReconciler) reconcileIngress(ctx context.Context, directus *directusv1.Directus) error {
	if len(directus.Spec.Ingress.Hosts) == 0 {
		return nil
	}

	pathType := networkingv1.PathTypePrefix
	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        directus.Name,
			Namespace:   directus.Namespace,
			Annotations: directus.Spec.Ingress.Annotations,
		},
		Spec: networkingv1.IngressSpec{
			IngressClassName: &directus.Spec.Ingress.ClassName,
			TLS:              directus.Spec.Ingress.TLS,
			Rules:            []networkingv1.IngressRule{},
		},
	}

	// Convert hosts
	for _, host := range directus.Spec.Ingress.Hosts {
		rule := networkingv1.IngressRule{
			Host: host.Host,
			IngressRuleValue: networkingv1.IngressRuleValue{
				HTTP: &networkingv1.HTTPIngressRuleValue{
					Paths: []networkingv1.HTTPIngressPath{},
				},
			},
		}

		for _, path := range host.Paths {
			rule.HTTP.Paths = append(rule.HTTP.Paths, networkingv1.HTTPIngressPath{
				Path:     path.Path,
				PathType: &pathType,
				Backend: networkingv1.IngressBackend{
					Service: &networkingv1.IngressServiceBackend{
						Name: directus.Name,
						Port: networkingv1.ServiceBackendPort{
							Number: directus.Spec.Service.Port,
						},
					},
				},
			})
		}

		ingress.Spec.Rules = append(ingress.Spec.Rules, rule)
	}

	if err := controllerutil.SetControllerReference(directus, ingress, r.Scheme); err != nil {
		return err
	}

	found := &networkingv1.Ingress{}
	err := r.Get(ctx, types.NamespacedName{Name: ingress.Name, Namespace: ingress.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		return r.Create(ctx, ingress)
	} else if err != nil {
		return err
	}

	// Update ingress
	found.Spec = ingress.Spec
	found.Annotations = ingress.Annotations
	return r.Update(ctx, found)
}

func (r *DirectusReconciler) reconcileHPA(ctx context.Context, directus *directusv1.Directus) error {
	hpa := &autoscalingv2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      directus.Name,
			Namespace: directus.Namespace,
		},
		Spec: autoscalingv2.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv2.CrossVersionObjectReference{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
				Name:       directus.Name,
			},
			MinReplicas: &directus.Spec.Autoscaling.MinReplicas,
			MaxReplicas: directus.Spec.Autoscaling.MaxReplicas,
			Metrics:     []autoscalingv2.MetricSpec{},
		},
	}

	// Add CPU metrics if specified
	if directus.Spec.Autoscaling.TargetCPUUtilizationPercentage > 0 {
		hpa.Spec.Metrics = append(hpa.Spec.Metrics, autoscalingv2.MetricSpec{
			Type: autoscalingv2.ResourceMetricSourceType,
			Resource: &autoscalingv2.ResourceMetricSource{
				Name: corev1.ResourceCPU,
				Target: autoscalingv2.MetricTarget{
					Type:               autoscalingv2.UtilizationMetricType,
					AverageUtilization: &directus.Spec.Autoscaling.TargetCPUUtilizationPercentage,
				},
			},
		})
	}

	// Add Memory metrics if specified
	if directus.Spec.Autoscaling.TargetMemoryUtilizationPercentage > 0 {
		hpa.Spec.Metrics = append(hpa.Spec.Metrics, autoscalingv2.MetricSpec{
			Type: autoscalingv2.ResourceMetricSourceType,
			Resource: &autoscalingv2.ResourceMetricSource{
				Name: corev1.ResourceMemory,
				Target: autoscalingv2.MetricTarget{
					Type:               autoscalingv2.UtilizationMetricType,
					AverageUtilization: &directus.Spec.Autoscaling.TargetMemoryUtilizationPercentage,
				},
			},
		})
	}

	if err := controllerutil.SetControllerReference(directus, hpa, r.Scheme); err != nil {
		return err
	}

	found := &autoscalingv2.HorizontalPodAutoscaler{}
	err := r.Get(ctx, types.NamespacedName{Name: hpa.Name, Namespace: hpa.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		return r.Create(ctx, hpa)
	} else if err != nil {
		return err
	}

	// Update HPA
	found.Spec = hpa.Spec
	return r.Update(ctx, found)
}

func (r *DirectusReconciler) updateStatus(ctx context.Context, directus *directusv1.Directus) error {
	// Get deployment status
	deployment := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{Name: directus.Name, Namespace: directus.Namespace}, deployment)
	if err != nil {
		directus.Status.Phase = "Failed"
		directus.Status.Message = fmt.Sprintf("Failed to get deployment: %v", err)
	} else {
		directus.Status.Replicas = deployment.Status.Replicas
		directus.Status.ReadyReplicas = deployment.Status.ReadyReplicas

		if deployment.Status.ReadyReplicas == directus.Spec.ReplicaCount {
			directus.Status.Phase = "Running"
			directus.Status.Message = "All replicas are ready"
		} else {
			directus.Status.Phase = "Pending"
			directus.Status.Message = fmt.Sprintf("Waiting for replicas: %d/%d ready", deployment.Status.ReadyReplicas, directus.Spec.ReplicaCount)
		}
	}

	// Update conditions
	now := metav1.Now()
	conditions := []metav1.Condition{
		{
			Type:               "Ready",
			Status:             metav1.ConditionUnknown,
			Reason:             "Reconciling",
			Message:            "Reconciling Directus resources",
			LastTransitionTime: now,
		},
	}

	if directus.Status.Phase == "Running" {
		conditions[0].Status = metav1.ConditionTrue
		conditions[0].Reason = "Ready"
		conditions[0].Message = "Directus is running"
	}

	directus.Status.Conditions = conditions
	directus.Status.DatabaseReady = true // Simplified for now
	directus.Status.RedisReady = directus.Spec.Redis.Enabled
	directus.Status.IngressReady = directus.Spec.Ingress.Enabled

	return r.Status().Update(ctx, directus)
}

// Helper methods
func (r *DirectusReconciler) getLabels(directus *directusv1.Directus) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       "directus",
		"app.kubernetes.io/instance":   directus.Name,
		"app.kubernetes.io/version":    directus.Spec.Image.Tag,
		"app.kubernetes.io/component":  "directus",
		"app.kubernetes.io/managed-by": "directus-operator",
	}
}

func (r *DirectusReconciler) getServiceAccountName(directus *directusv1.Directus) string {
	if directus.Spec.ServiceAccount.Name != "" {
		return directus.Spec.ServiceAccount.Name
	}
	return directus.Name + "-sa"
}

func (r *DirectusReconciler) getApplicationSecretName(directus *directusv1.Directus) string {
	if directus.Spec.ApplicationSecretName != "" {
		return directus.Spec.ApplicationSecretName
	}
	return directus.Name + "-application-secret"
}

func (r *DirectusReconciler) buildConfigMapData(directus *directusv1.Directus) map[string]string {
	data := map[string]string{
		"ADMIN_EMAIL": directus.Spec.AdminEmail,
	}

	// Database configuration
	if directus.Spec.Database.Engine != "" {
		data["DB_CLIENT"] = directus.Spec.Database.Engine
	}
	if directus.Spec.Database.Host != "" {
		data["DB_HOST"] = directus.Spec.Database.Host
	}
	if directus.Spec.Database.Port > 0 {
		data["DB_PORT"] = strconv.Itoa(int(directus.Spec.Database.Port))
	}
	if directus.Spec.Database.Database != "" {
		data["DB_DATABASE"] = directus.Spec.Database.Database
	}
	if directus.Spec.Database.Username != "" {
		data["DB_USER"] = directus.Spec.Database.Username
	}

	// Redis configuration
	data["REDIS_ENABLED"] = strconv.FormatBool(directus.Spec.Redis.Enabled)
	if directus.Spec.Redis.Enabled {
		if directus.Spec.Redis.Host != "" {
			data["REDIS_HOST"] = directus.Spec.Redis.Host
		}
		if directus.Spec.Redis.Port > 0 {
			data["REDIS_PORT"] = strconv.Itoa(int(directus.Spec.Redis.Port))
		}
	}

	return data
}

func (r *DirectusReconciler) buildContainers(directus *directusv1.Directus) []corev1.Container {
	container := corev1.Container{
		Name:            "directus",
		Image:           directus.Spec.Image.Repository + ":" + directus.Spec.Image.Tag,
		ImagePullPolicy: directus.Spec.Image.PullPolicy,
		Ports: []corev1.ContainerPort{
			{
				Name:          "http",
				ContainerPort: 8055,
				Protocol:      corev1.ProtocolTCP,
			},
		},
		EnvFrom: []corev1.EnvFromSource{
			{
				ConfigMapRef: &corev1.ConfigMapEnvSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: directus.Name + "-configmap",
					},
				},
			},
		},
		Env:             directus.Spec.ExtraEnvVars,
		Resources:       directus.Spec.Resources,
		VolumeMounts:    directus.Spec.ExtraVolumeMounts,
		SecurityContext: directus.Spec.SecurityContext,
	}

	// Add environment variables from secrets
	for _, secretName := range directus.Spec.AttachExistingSecrets {
		container.EnvFrom = append(container.EnvFrom, corev1.EnvFromSource{
			SecretRef: &corev1.SecretEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: secretName,
				},
			},
		})
	}

	// Add database password if using external database
	if directus.Spec.Database.ExistingSecret != "" {
		container.Env = append(container.Env, corev1.EnvVar{
			Name: "DB_PASSWORD",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: directus.Spec.Database.ExistingSecret,
					},
					Key: "password",
				},
			},
		})
	}

	// Add application secret if created
	if directus.Spec.CreateApplicationSecret {
		container.EnvFrom = append(container.EnvFrom, corev1.EnvFromSource{
			SecretRef: &corev1.SecretEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: r.getApplicationSecretName(directus),
				},
			},
		})
	}

	// Add probes
	if directus.Spec.EnableLivenessProbe {
		container.LivenessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: "/",
					Port: intstr.FromString("http"),
				},
			},
			InitialDelaySeconds: 60, // Wait 60 seconds before starting probes
			PeriodSeconds:       10, // Check every 10 seconds
			TimeoutSeconds:      5,  // 5 second timeout
			FailureThreshold:    5,  // Allow 5 failures before restart
		}
	}

	if directus.Spec.EnableReadinessProbe {
		container.ReadinessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: "/",
					Port: intstr.FromString("http"),
				},
			},
			InitialDelaySeconds: 30, // Wait 30 seconds before starting readiness checks
			PeriodSeconds:       5,  // Check every 5 seconds
			TimeoutSeconds:      3,  // 3 second timeout
			FailureThreshold:    3,  // Allow 3 failures
		}
	}

	if directus.Spec.EnableStartupProbe {
		container.StartupProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: "/",
					Port: intstr.FromString("http"),
				},
			},
			InitialDelaySeconds: 10, // Start checking after 10 seconds
			PeriodSeconds:       10, // Check every 10 seconds
			TimeoutSeconds:      3,  // 3 second timeout
			FailureThreshold:    30, // Allow up to 5 minutes for startup (30 * 10s)
		}
	}

	return []corev1.Container{container}
}

// SetupWithManager sets up the controller with the Manager.
func (r *DirectusReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&directusv1.Directus{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&corev1.Secret{}).
		Owns(&corev1.ServiceAccount{}).
		Owns(&networkingv1.Ingress{}).
		Owns(&autoscalingv2.HorizontalPodAutoscaler{}).
		Named("directus").
		Complete(r)
}
