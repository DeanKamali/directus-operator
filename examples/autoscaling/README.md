# Directus with Autoscaling

This example demonstrates deploying Directus with Horizontal Pod Autoscaler (HPA) for automatic scaling based on CPU and memory usage.

## Features

- **Horizontal Pod Autoscaler**: Automatic scaling based on metrics
- **CPU-based Scaling**: Scales based on CPU utilization
- **Memory-based Scaling**: Scales based on memory utilization
- **Custom Metrics**: Support for custom metrics (optional)
- **Production Configuration**: PostgreSQL + Redis for scalability

## Prerequisites

1. **Metrics Server**: Required for HPA to function
2. **PostgreSQL Database**: External database for shared state
3. **Redis**: For shared sessions and caching
4. **Resource Requests**: Properly configured resource requests are essential

## Metrics Server

Install metrics-server if not already available:

```bash
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
```

Verify metrics server is working:
```bash
kubectl top nodes
kubectl top pods
```

## Autoscaling Configuration

The HPA configuration includes:

- **Min Replicas**: 2 (for high availability)
- **Max Replicas**: 10 (adjust based on your needs)
- **CPU Target**: 70% average utilization
- **Memory Target**: 80% average utilization

### Scaling Behavior

```yaml
behavior:
  scaleUp:
    stabilizationWindowSeconds: 300    # 5 minutes
    policies:
    - type: Pods
      value: 2
      periodSeconds: 60               # Max 2 pods per minute
  scaleDown:
    stabilizationWindowSeconds: 600    # 10 minutes
    policies:
    - type: Pods
      value: 1
      periodSeconds: 120              # Max 1 pod every 2 minutes
```

## Deployment

1. **Deploy database dependencies**:
   ```bash
   kubectl apply -f dependencies/
   ```

2. **Deploy Directus with autoscaling**:
   ```bash
   kubectl apply -f directus.yaml
   ```

3. **Verify HPA**:
   ```bash
   kubectl get hpa
   kubectl describe hpa directus-autoscaling-hpa
   ```

4. **Monitor scaling events**:
   ```bash
   kubectl get events --sort-by=.metadata.creationTimestamp
   ```

## Monitoring Autoscaling

### Check HPA Status
```bash
# View current HPA status
kubectl get hpa directus-autoscaling-hpa

# Detailed HPA information
kubectl describe hpa directus-autoscaling-hpa

# Watch HPA in real-time
kubectl get hpa -w
```

### Check Pod Metrics
```bash
# Current resource usage
kubectl top pods -l app.kubernetes.io/instance=directus-autoscaling

# Watch resource usage
watch kubectl top pods -l app.kubernetes.io/instance=directus-autoscaling
```

## Load Testing

To test autoscaling, you can generate load:

### Using kubectl run
```bash
# Create a load generator pod
kubectl run load-generator --image=busybox -- sh -c \
  "while true; do wget -q -O- http://directus-autoscaling/; done"

# Clean up
kubectl delete pod load-generator
```

### Using Apache Bench (ab)
```bash
# Install ab and generate load
kubectl run ab-test --image=httpd:alpine -- ab -n 10000 -c 50 http://directus-autoscaling/
```

### Using hey (HTTP load testing tool)
```bash
kubectl run hey-test --image=rcmorano/hey -- hey -z 5m -c 50 http://directus-autoscaling/
```

## Scaling Scenarios

### Expected Behavior

1. **Normal Load**: Maintains minimum 2 replicas
2. **High CPU Load**: Scales up when CPU > 70%
3. **High Memory Usage**: Scales up when memory > 80%
4. **Load Decrease**: Scales down gradually after stabilization period

### Scaling Timeline

- **Scale Up**: ~3-5 minutes after threshold breach
- **Scale Down**: ~10-15 minutes after metrics normalize
- **Stabilization**: Built-in delays prevent flapping

## Custom Metrics (Advanced)

For production workloads, consider scaling on custom metrics:

### Application Metrics
- Active database connections
- Response time percentiles
- Queue length
- Cache hit ratio

### External Metrics
- AWS CloudWatch metrics
- Prometheus metrics
- Datadog metrics

Example custom metrics configuration:
```yaml
metrics:
- type: Pods
  pods:
    metric:
      name: directus_active_connections
    target:
      type: AverageValue
      averageValue: "100"
```

## Troubleshooting

### Common Issues

1. **HPA Not Scaling**:
   - Check metrics server: `kubectl get --raw /apis/metrics.k8s.io/v1beta1/nodes`
   - Verify resource requests are set in pod spec
   - Check HPA events: `kubectl describe hpa`

2. **Metrics Not Available**:
   - Ensure metrics-server is running
   - Check resource requests in deployment
   - Verify pod is ready and serving metrics

3. **Frequent Scaling (Flapping)**:
   - Adjust stabilization windows
   - Review target utilization percentages
   - Check for resource contention

### Debug Commands

```bash
# Check metrics server
kubectl get apiservice v1beta1.metrics.k8s.io

# Raw metrics API
kubectl get --raw /apis/metrics.k8s.io/v1beta1/namespaces/default/pods

# HPA controller logs
kubectl logs -n kube-system deployment/horizontal-pod-autoscaler

# Check resource requests
kubectl describe pod <pod-name> | grep -A 5 "Requests:"
```

## Best Practices

### Resource Management
1. **Set Appropriate Requests**: HPA requires resource requests
2. **Monitor Resource Usage**: Use monitoring tools to understand patterns
3. **Right-size Limits**: Prevent resource starvation

### Scaling Configuration
1. **Conservative Targets**: Start with higher target percentages
2. **Stabilization Windows**: Prevent rapid scaling fluctuations
3. **Gradual Scaling**: Limit scaling velocity

### Application Design
1. **Stateless Design**: Ensure pods can be safely terminated
2. **Graceful Shutdown**: Handle SIGTERM properly
3. **Health Checks**: Implement proper readiness/liveness probes

## Performance Optimization

### Database Connections
- Use connection pooling
- Monitor active connections
- Consider read replicas for scaling reads

### Caching Strategy
- Redis for shared cache
- CDN for static assets
- Application-level caching

### Resource Allocation
- Profile your application
- Monitor memory leaks
- Optimize container images 