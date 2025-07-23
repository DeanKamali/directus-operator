# Directus Operator Examples

This directory contains example configurations for deploying Directus using the Directus Operator. Each example demonstrates different deployment scenarios and configurations.

**TL;DR**: We've made this super easy with a Makefile. Just run `make basic` to get started! 🚀

## Prerequisites

Before deploying any examples, ensure you have:

1. A Kubernetes cluster with the Directus Operator installed
2. `kubectl` configured to access your cluster
3. Necessary permissions to create resources in your namespace

## Examples Overview

| Example | Description | Use Case |
|---------|-------------|----------|
| [basic](./basic/) | Simple SQLite-based deployment | Quick testing and development |
| [production](./production/) | PostgreSQL + Redis deployment | Production workloads |
| [with-ingress](./with-ingress/) | Includes ingress configuration | Web-accessible deployment |
| [autoscaling](./autoscaling/) | Horizontal Pod Autoscaler setup | High-availability production |

## Quick Start

We've made this ridiculously easy with a Makefile. Just run:

```bash
# Deploy basic example (SQLite, perfect for testing)
make basic

# Or go straight to production-ready
make production

# Want it web-accessible with HTTPS?
make ingress

# Need auto-scaling?
make autoscaling

# See what's running
make status

# Clean up everything
make clean-all
```

That's it! No more copy-pasting YAML files and wondering if you got the indentation right.

## Example Structure

Each example directory contains:
- `directus.yaml` - Main Directus resource configuration
- `dependencies/` - Additional required resources (databases, secrets, etc.)
- `README.md` - Specific instructions for that example

## Customization

Before deploying, review and customize these common settings:

- **Namespace**: Update `metadata.namespace` in all resources
- **Admin Email**: Set `spec.adminEmail` to your email address
- **Image**: Choose appropriate Directus container image tag
- **Resources**: Adjust CPU/memory limits based on your requirements
- **Ingress**: Configure hostname and TLS settings for your domain

## Makefile Commands

We've included a simple Makefile with just the essentials:

```bash
make help        # Show all commands
make basic       # Deploy SQLite example  
make production  # Deploy PostgreSQL + Redis
make ingress     # Deploy with HTTPS
make autoscaling # Deploy with auto-scaling
make status      # Check what's running
make access      # Access basic deployment (port-forward + password)
make clean       # Remove everything
```

## Support

For more information about configuration options, see the [main README](../README.md) or consult the [API documentation](../api/v1/directus_types.go). 