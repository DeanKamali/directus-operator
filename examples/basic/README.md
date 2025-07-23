# Basic Directus Deployment

This example provides a minimal Directus deployment using SQLite for quick testing and development purposes.

## Features

- **SQLite Database**: No external database required
- **No Redis**: Simplified caching configuration
- **Minimal Resources**: Low CPU and memory requirements
- **Basic Security**: Default security contexts and probes

## Use Cases

- Development and testing environments
- Quick demos and proof-of-concepts
- Learning the Directus Operator
- Local Kubernetes testing (minikube, kind, etc.)

## Deployment

1. **Deploy the Directus instance**:
   ```bash
   kubectl apply -f directus.yaml
   ```

2. **Check the deployment status**:
   ```bash
   kubectl get directus basic-directus
   kubectl get pods -l app.kubernetes.io/name=basic-directus
   ```

3. **Access Directus**:
   ```bash
   # Port forward to access locally
   kubectl port-forward svc/basic-directus 8080:80
   ```

   Then open http://localhost:8080 in your browser.

## Default Credentials

- **Email**: `admin@example.com`
- **Password**: Generated automatically and stored in secret `basic-directus-secrets`

To retrieve the admin password:
```bash
kubectl get secret basic-directus-secrets -o jsonpath='{.data.ADMIN_PASSWORD}' | base64 -d
```

## Customization

Before deploying, you may want to customize:

- **Admin Email**: Change `spec.adminEmail` in `directus.yaml`
- **Namespace**: Add/change `metadata.namespace` to deploy in a specific namespace
- **Resource Limits**: Adjust `spec.resources` based on your cluster capacity

## Clean Up

To remove the deployment:
```bash
kubectl delete -f directus.yaml
```

## Limitations

- **Data Persistence**: SQLite data is stored in the container and will be lost if the pod restarts
- **Single Instance**: Not suitable for production use with multiple replicas
- **No Caching**: Redis caching is disabled for simplicity 