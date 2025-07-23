# Production Directus Deployment

This example provides a production-ready Directus deployment with PostgreSQL database and Redis for caching and session management.

## Features

- **PostgreSQL Database**: Reliable, scalable database backend
- **Redis Caching**: Enhanced performance with Redis for caching and sessions
- **Production Resources**: Appropriate CPU and memory allocation
- **Security Hardening**: Enhanced security contexts and configurations
- **Persistent Storage**: Data persistence using persistent volumes
- **Health Checks**: Comprehensive startup, liveness, and readiness probes

## Architecture

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Directus  │────│ PostgreSQL  │    │    Redis    │
│             │    │  Database   │    │   Cache     │
└─────────────┘    └─────────────┘    └─────────────┘
```

## Prerequisites

1. **PostgreSQL Database**: You need a PostgreSQL instance accessible from your cluster
2. **Redis Instance**: Optional but recommended for production performance
3. **Persistent Volumes**: Ensure your cluster supports persistent volume claims
4. **Secrets**: Database and Redis credentials should be stored as Kubernetes secrets

## Deployment

1. **Create database credentials secret**:
   ```bash
   kubectl create secret generic postgresql-credentials \
     --from-literal=username=directus_user \
     --from-literal=password=your-secure-password \
     --from-literal=database=directus_db \
     --from-literal=host=your-postgres-host \
     --from-literal=port=5432
   ```

2. **Create Redis credentials secret** (if using external Redis):
   ```bash
   kubectl create secret generic redis-credentials \
     --from-literal=host=your-redis-host \
     --from-literal=port=6379 \
     --from-literal=password=your-redis-password
   ```

3. **Deploy dependencies** (if using the included PostgreSQL and Redis):
   ```bash
   kubectl apply -f dependencies/
   ```

4. **Deploy the Directus instance**:
   ```bash
   kubectl apply -f directus.yaml
   ```

5. **Check the deployment status**:
   ```bash
   kubectl get directus production-directus
   kubectl get pods -l app.kubernetes.io/instance=production-directus
   ```

## Access Methods

### Port Forward (Development)
```bash
kubectl port-forward svc/production-directus 8080:80
```

### Ingress (Production)
See the [with-ingress example](../with-ingress/) for ingress configuration.

## Configuration

### Database Configuration
- **Engine**: PostgreSQL
- **Connection Pooling**: Configured for production workloads
- **SSL**: Recommended for production environments
- **Backup**: Ensure regular database backups are configured

### Redis Configuration
- **Caching**: Enabled for improved performance
- **Sessions**: Redis-backed session storage
- **Rate Limiting**: Redis-based rate limiting

### Security
- **Service Account**: Dedicated service account with minimal permissions
- **Pod Security Context**: Non-root user with restricted permissions
- **Security Context**: Dropped capabilities and read-only root filesystem where possible
- **Network Policies**: Consider implementing network policies for additional security

## Monitoring

The deployment includes annotations for Prometheus monitoring:
- Metrics endpoint: `/metrics`
- Port: `8055`

## Scaling

This deployment supports horizontal scaling:
```bash
kubectl patch directus production-directus -p '{"spec":{"replicaCount":3}}'
```

## Customization

Review and customize these settings for your environment:

- **Database Connection**: Update database host, credentials, and SSL settings
- **Redis Configuration**: Adjust Redis host and authentication
- **Resource Limits**: Scale CPU and memory based on your workload
- **Environment Variables**: Configure PUBLIC_URL, email settings, etc.
- **Storage**: Configure persistent volumes for file uploads

## Backup and Recovery

1. **Database Backups**: Implement regular PostgreSQL backups
2. **File Storage**: Backup uploaded files if using persistent volumes
3. **Configuration**: Backup your Directus configuration and custom code

## Troubleshooting

### Common Issues

1. **Database Connection Failed**:
   - Verify database credentials in the secret
   - Check network connectivity to PostgreSQL
   - Ensure database exists and user has proper permissions

2. **Redis Connection Issues**:
   - Verify Redis host and port configuration
   - Check Redis authentication settings
   - Test Redis connectivity from a debug pod

3. **Performance Issues**:
   - Monitor resource usage with `kubectl top pods`
   - Check database query performance
   - Verify Redis cache hit rates

### Debug Commands
```bash
# Check pod logs
kubectl logs -l app.kubernetes.io/instance=production-directus

# Check database connectivity
kubectl exec -it <directus-pod> -- env | grep DB_

# Check Redis connectivity
kubectl exec -it <directus-pod> -- env | grep REDIS
``` 