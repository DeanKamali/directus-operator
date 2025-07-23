# Directus Operator

A Kubernetes operator for managing Directus deployments, converted from the original Helm chart.

## Overview

The Directus Operator provides a declarative way to deploy and manage Directus instances on Kubernetes. It automatically handles:

- Directus application deployment
- Service and ingress configuration
- Database connectivity
- Redis integration
- Secrets management
- Horizontal Pod Autoscaling
- Health checks and probes

## Features

✅ **Converted from Helm Chart**: All functionality from the original Directus Helm chart
✅ **Declarative Management**: Define desired state with Custom Resources
✅ **Automatic Reconciliation**: Ensures actual state matches desired state
✅ **Status Reporting**: Real-time status and health information
✅ **Scaling Support**: Built-in horizontal pod autoscaling
✅ **Security**: Service accounts, security contexts, and secrets management
✅ **Flexibility**: Supports various deployment scenarios

## Installation

### Prerequisites

- Kubernetes cluster (v1.25+)
- kubectl configured to access your cluster
- Admin permissions to install CRDs and RBAC

### Deploy the Operator

1. **Install CRDs and deploy the operator:**
   ```bash
   # Apply CRDs
   kubectl apply -f config/crd/bases/

   # Deploy the operator
   kubectl apply -f config/manager/
   ```

2. **Verify the operator is running:**
   ```bash
   kubectl get pods -n directus-operator-system
   ```

## Usage

### Basic Example

Create a simple Directus instance:

```yaml
apiVersion: directus.example.com/v1
kind: Directus
metadata:
  name: my-directus
  namespace: default
spec:
  replicaCount: 1
  image:
    repository: directus/directus
    tag: "11.8.0"
  adminEmail: "admin@example.com"
  createApplicationSecret: true
  enableLivenessProbe: true
  enableReadinessProbe: true
```

Apply it:
```bash
kubectl apply -f my-directus.yaml
```

### Advanced Example with Database and Ingress

```yaml
apiVersion: directus.example.com/v1
kind: Directus
metadata:
  name: production-directus
  namespace: directus
spec:
  replicaCount: 3
  
  image:
    repository: directus/directus
    tag: "11.8.0"
    pullPolicy: IfNotPresent
  
  adminEmail: "admin@company.com"
  
  # Database configuration
  database:
    engine: mysql
    host: mysql-service.database.svc.cluster.local
    port: 3306
    database: directus_prod
    username: directus_user
    existingSecret: mysql-credentials
  
  # Redis for caching and sessions
  redis:
    enabled: true
    host: redis-service.cache.svc.cluster.local
    port: 6379
  
  # Ingress for external access
  ingress:
    enabled: true
    className: "nginx"
    annotations:
      cert-manager.io/cluster-issuer: "letsencrypt-prod"
      nginx.ingress.kubernetes.io/proxy-body-size: "100m"
    hosts:
      - host: directus.company.com
        paths:
          - path: /
            pathType: Prefix
    tls:
      - secretName: directus-tls
        hosts:
          - directus.company.com
  
  # Autoscaling
  autoscaling:
    enabled: true
    minReplicas: 3
    maxReplicas: 10
    targetCPUUtilizationPercentage: 70
  
  # Resources
  resources:
    limits:
      cpu: 1000m
      memory: 1Gi
    requests:
      cpu: 500m
      memory: 512Mi
  
  # Security
  podSecurityContext:
    fsGroup: 2000
  securityContext:
    runAsNonRoot: true
    runAsUser: 1000
```

## Configuration Reference

### Image Configuration
```yaml
spec:
  image:
    repository: directus/directus  # Container image repository
    tag: "11.8.0"                 # Image tag
    pullPolicy: IfNotPresent      # Pull policy: Always, IfNotPresent, Never
```

### Database Configuration
```yaml
spec:
  database:
    engine: mysql                 # Database engine: mysql, postgresql, sqlite
    host: mysql-service          # Database hostname
    port: 3306                   # Database port
    database: directus_db        # Database name
    username: directus_user      # Database username
    existingSecret: db-secret    # Secret containing database credentials
```

### Redis Configuration
```yaml
spec:
  redis:
    enabled: true               # Enable Redis
    host: redis-service        # Redis hostname
    port: 6379                 # Redis port
    existingSecret: redis-secret # Secret containing Redis credentials
```

### Ingress Configuration
```yaml
spec:
  ingress:
    enabled: true              # Enable ingress
    enableTLS: true           # Enable TLS in PUBLIC_URL
    className: "nginx"        # Ingress class
    annotations:              # Ingress annotations
      cert-manager.io/cluster-issuer: "letsencrypt-prod"
    hosts:
      - host: directus.example.com
        paths:
          - path: /
            pathType: Prefix
    tls:
      - secretName: directus-tls
        hosts:
          - directus.example.com
```

### Autoscaling Configuration
```yaml
spec:
  autoscaling:
    enabled: true                          # Enable HPA
    minReplicas: 2                        # Minimum replicas
    maxReplicas: 10                       # Maximum replicas
    targetCPUUtilizationPercentage: 70    # Target CPU utilization
    targetMemoryUtilizationPercentage: 80 # Target memory utilization
```

## Secret Management

The operator can create and manage secrets for you:

### Application Secrets
```yaml
spec:
  createApplicationSecret: true           # Create ADMIN_PASSWORD, KEY, SECRET
  applicationSecretName: my-app-secrets  # Optional: custom secret name
```

### Existing Secrets
```yaml
spec:
  attachExistingSecrets:                 # Attach existing secrets as env vars
    - database-credentials
    - external-api-keys
```

## Status and Monitoring

Check the status of your Directus instance:

```bash
kubectl get directus my-directus -o yaml
```

The status section provides information about:
- Current phase (Pending, Running, Failed)
- Ready replicas vs desired replicas
- Database and Redis connectivity
- Ingress status
- Conditions and events

Example status:
```yaml
status:
  conditions:
    - type: Ready
      status: "True"
      reason: Ready
      message: Directus is running
  readyReplicas: 3
  replicas: 3
  phase: Running
  message: All replicas are ready
  databaseReady: true
  redisReady: true
  ingressReady: true
```

## Comparison with Helm Chart

| Feature | Helm Chart | Operator |
|---------|------------|----------|
| **Deployment** | `helm install` | `kubectl apply` |
| **Updates** | `helm upgrade` | Automatic reconciliation |
| **Rollbacks** | `helm rollback` | Edit CR or use kubectl |
| **Status** | `helm status` | `kubectl get directus` |
| **Customization** | values.yaml | Custom Resource spec |
| **Lifecycle** | Manual commands | Declarative management |

## Migration from Helm

To migrate from the Helm chart to the operator:

1. **Export current values:**
   ```bash
   helm get values my-directus > current-values.yaml
   ```

2. **Create equivalent Custom Resource:**
   Use the values to create a Directus CR with equivalent configuration.

3. **Deploy with operator:**
   ```bash
   kubectl apply -f my-directus.yaml
   ```

4. **Remove Helm release:**
   ```bash
   helm uninstall my-directus
   ```

## Troubleshooting

### Check Operator Logs
```bash
kubectl logs -n directus-operator-system deployment/directus-operator-controller-manager
```

### Check Resource Events
```bash
kubectl describe directus my-directus
```

### Common Issues

1. **Pods not starting**: Check database connectivity and secrets
2. **Ingress not working**: Verify ingress controller and DNS
3. **Database connection failed**: Check database credentials and network policies
4. **Operator not reconciling**: Check operator logs and RBAC permissions

## Development

### Building and Testing

```bash
# Generate manifests
make manifests

# Run tests
make test

# Build and push image
make docker-build docker-push IMG=your-registry/directus-operator:tag

# Deploy to cluster
make deploy IMG=your-registry/directus-operator:tag
```

### Local Development

```bash
# Install CRDs
make install

# Run controller locally
make run
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

Licensed under the Apache License, Version 2.0. See LICENSE for details. 