# Directus with Ingress

This example demonstrates how to deploy Directus with ingress configuration for web accessibility, including TLS/SSL setup.

## Features

- **Ingress Controller**: NGINX ingress configuration
- **TLS/SSL**: Automatic certificate management with cert-manager
- **Custom Domain**: Configure your own domain name
- **Path-based Routing**: Support for path-based routing
- **Security Headers**: Enhanced security with proper headers

## Prerequisites

1. **Ingress Controller**: Ensure you have an ingress controller installed (e.g., NGINX, Traefik)
2. **cert-manager**: For automatic SSL certificate provisioning (optional)
3. **DNS Configuration**: Your domain should point to your cluster's ingress IP
4. **TLS Certificate**: Either use cert-manager or provide your own certificate

## Ingress Controllers

### NGINX Ingress Controller
```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.2/deploy/static/provider/cloud/deploy.yaml
```

### Traefik (Alternative)
```bash
helm repo add traefik https://helm.traefik.io/traefik
helm install traefik traefik/traefik
```

## Certificate Management

### Option 1: cert-manager (Recommended)
Install cert-manager for automatic SSL certificates:

```bash
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.2/cert-manager.yaml
```

Create a ClusterIssuer for Let's Encrypt:
```bash
kubectl apply -f cert-issuer.yaml
```

### Option 2: Manual Certificate
If you have your own certificate:

```bash
kubectl create secret tls directus-tls \
  --cert=path/to/tls.crt \
  --key=path/to/tls.key
```

## Deployment

1. **Update domain configuration**:
   Edit `directus.yaml` and replace `directus.example.com` with your actual domain.

2. **Deploy cert-manager ClusterIssuer** (if using cert-manager):
   ```bash
   kubectl apply -f cert-issuer.yaml
   ```

3. **Deploy Directus with ingress**:
   ```bash
   kubectl apply -f directus.yaml
   ```

4. **Verify ingress**:
   ```bash
   kubectl get ingress directus-ingress
   kubectl describe ingress directus-ingress
   ```

## Configuration

### Domain Setup
Update these values in `directus.yaml`:
- `spec.ingress.hosts[0].host`: Your domain name
- `spec.extraEnvVars` → `PUBLIC_URL`: Your full URL

### TLS Configuration
- **Automatic (cert-manager)**: Certificates are automatically provisioned
- **Manual**: Update `spec.ingress.tls` with your certificate secret name

### Security Headers
The ingress includes security headers:
- `X-Frame-Options`: Prevents clickjacking
- `X-Content-Type-Options`: Prevents MIME sniffing
- `X-XSS-Protection`: XSS protection
- `Strict-Transport-Security`: Forces HTTPS

## Access Methods

Once deployed, access Directus at:
- **HTTPS**: `https://your-domain.com`
- **HTTP**: Automatically redirects to HTTPS

## Customization

### Path-based Routing
To deploy under a subpath (e.g., `/directus`):

1. Update ingress path in `directus.yaml`
2. Set `PUBLIC_URL` to include the subpath
3. Configure Directus `ROOT_PATH` environment variable

### Multiple Domains
To support multiple domains, add additional hosts to the ingress configuration.

### Load Balancer Annotations
Add cloud-specific annotations for your load balancer:

```yaml
# AWS ALB
service.beta.kubernetes.io/aws-load-balancer-type: nlb

# GCP
cloud.google.com/load-balancer-type: External

# Azure
service.beta.kubernetes.io/azure-load-balancer-internal: "false"
```

## Troubleshooting

### Common Issues

1. **503 Service Unavailable**:
   - Check if Directus pods are running and ready
   - Verify service endpoints: `kubectl get endpoints`

2. **Certificate Issues**:
   - Check cert-manager logs: `kubectl logs -n cert-manager deployment/cert-manager`
   - Verify certificate status: `kubectl describe certificate directus-tls`

3. **DNS Resolution**:
   - Verify DNS points to your ingress IP
   - Check ingress controller external IP: `kubectl get svc -n ingress-nginx`

4. **Redirect Loops**:
   - Ensure `PUBLIC_URL` matches your ingress configuration
   - Check `CORS_ORIGIN` environment variable

### Debug Commands

```bash
# Check ingress configuration
kubectl describe ingress directus-ingress

# Test DNS resolution
nslookup your-domain.com

# Check certificate
kubectl get certificate
kubectl describe certificate directus-tls

# Test from inside cluster
kubectl run debug --image=curlimages/curl -it --rm -- sh
curl -H "Host: your-domain.com" http://directus-service
```

## Security Considerations

1. **TLS Configuration**: Always use HTTPS in production
2. **CORS Settings**: Configure appropriate CORS origins
3. **Rate Limiting**: Consider implementing rate limiting at ingress level
4. **WAF**: Use a Web Application Firewall for additional protection
5. **Network Policies**: Implement Kubernetes network policies for pod isolation 