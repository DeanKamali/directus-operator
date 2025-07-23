# Directus Operator 🚀

Hey there! 👋 So you want to run Directus on Kubernetes but don't want to deal with all the YAML complexity? You've come to the right place! This operator makes deploying Directus as easy as pie.

## What's This All About?

Look, we've all been there. You want to spin up Directus (that awesome headless CMS) on your Kubernetes cluster, but then you realize you need to:
- Set up a database (PostgreSQL? MySQL? SQLite?)
- Configure Redis for caching (do I even need it?)
- Handle ingress and TLS certificates (ugh, cert-manager again?)
- Make it scale when things get busy
- Keep it secure (because security matters!)

This operator does all that boring stuff for you. Just tell it what you want, and it'll handle the rest.

## What You Get Out of the Box

- 🗄️ **Database Magic**: PostgreSQL, MySQL, or SQLite - whatever floats your boat
- ⚡ **Redis Goodness**: Caching and session management that actually works
- 🌐 **Web-Ready**: Ingress with automatic TLS because who has time for manual certificates?
- 📈 **Auto-Scaling**: Your app gets busy? No problem, we'll spin up more pods
- 🔒 **Security First**: Production-ready security configs so you can sleep at night
- 💾 **Storage Sorted**: File uploads and persistent storage handled properly

## Getting Started

### What You'll Need

Before we dive in, make sure you've got:
- Go 1.24+ (for building the operator)
- Docker 17.03+ (for containerizing)
- kubectl 1.11.3+ (for talking to Kubernetes)
- A Kubernetes cluster v1.11.3+ (obviously!)

### The Quick Way (Recommended)

The fastest way to try this out is with our pre-made examples. We've got a Makefile that makes everything super easy:

```bash
# Clone the repo
git clone <your-repo-url>
cd directus-operator/examples

# Deploy a basic Directus instance (SQLite, perfect for testing)
make basic

# Access it locally
make port-forward-basic
# Now open http://localhost:8080

# Get the admin password
make get-password-basic
```

That's it! You've got Directus running. Want to clean up? Just run `make clean-basic`.

### Other Cool Examples

```bash
# Production setup with PostgreSQL and Redis
make production

# Web-accessible with automatic HTTPS
make ingress

# Auto-scaling production deployment
make autoscaling

# See what's running
make status

# Clean everything up
make clean-all
```

### Installing the Operator (The Hard Way)

If you want to install the operator itself (maybe you're developing or want to use your own registry):

**Build and push your image:**
```bash
make docker-build docker-push IMG=<your-registry>/directus-operator:tag
```

**Install the CRDs:**
```bash
make install
```

**Deploy the operator:**
```bash
make deploy IMG=<your-registry>/directus-operator:tag
```

> **Heads up**: If you get RBAC errors, you might need cluster-admin privileges. It happens! 🤷‍♂️

**Deploy Directus instances:**
```bash
# Use the built-in samples
kubectl apply -k config/samples/

# Or use our fancy examples (recommended!)
cd examples && make basic
```

>**NOTE**: Ensure that the samples have default values to test it out.

### Cleaning Up

**Remove your Directus instances:**
```bash
# If you used our examples
cd examples && make clean-all

# If you used the samples
kubectl delete -k config/samples/
```

**Remove the operator itself:**
```bash
make uninstall  # Removes the CRDs
make undeploy   # Removes the operator
```

## Project Distribution

Following the options to release and provide this solution to the users.

### By providing a bundle with all YAML files

1. Build the installer for the image built and published in the registry:

```sh
make build-installer IMG=<some-registry>/directus-operator:tag
```

**NOTE:** The makefile target mentioned above generates an 'install.yaml'
file in the dist directory. This file contains all the resources built
with Kustomize, which are necessary to install this project without its
dependencies.

2. Using the installer

Users can just run 'kubectl apply -f <URL for YAML BUNDLE>' to install
the project, i.e.:

```sh
kubectl apply -f https://raw.githubusercontent.com/<org>/directus-operator/<tag or branch>/dist/install.yaml
```

### By providing a Helm Chart

1. Build the chart using the optional helm plugin

```sh
operator-sdk edit --plugins=helm/v1-alpha
```

2. See that a chart was generated under 'dist/chart', and users
can obtain this solution from there.

**NOTE:** If you change the project, you need to update the Helm Chart
using the same command above to sync the latest changes. Furthermore,
if you create webhooks, you need to use the above command with
the '--force' flag and manually ensure that any custom configuration
previously added to 'dist/chart/values.yaml' or 'dist/chart/manager/manager.yaml'
is manually re-applied afterwards.

## Examples That Actually Work

We've put together some real-world examples that you can actually use. No more "left as an exercise for the reader" nonsense!

| What You Get | Why You'd Want It | How Complex |
|-------------|------------------|-------------|
| [Basic](./examples/basic/) | SQLite, single pod, just works | 🟢 Super Easy |
| [Production](./examples/production/) | PostgreSQL + Redis, proper setup | 🟡 Medium |
| [Web-Ready](./examples/with-ingress/) | Automatic HTTPS, custom domain | 🟡 Medium |
| [Auto-Scaling](./examples/autoscaling/) | Handles traffic spikes like a boss | 🔴 Advanced |

Each example has its own README with actual instructions that work. Revolutionary, I know!

## The "I Just Want It Working" Guide

Seriously, this is all you need:

```bash
# Get the code
git clone <this-repo>
cd directus-operator/examples

# Deploy (takes ~30 seconds)
make demo

# Access at http://localhost:8080
# Username: admin@example.com
# Password: run `make get-password-basic`
```

That's it. You now have Directus running on Kubernetes. 🎉

## Customizing Your Setup

### Database Choices (Pick Your Poison)
- **SQLite**: Great for dev/testing, terrible for production. You've been warned! 😄
- **PostgreSQL**: The gold standard. Reliable, fast, battle-tested.
- **MySQL**: If you're into that sort of thing. Works fine too.

### Caching (Make It Fast)
- **Redis**: Seriously, just use Redis. Your users will thank you.
- **File-based**: For when you want to keep things simple (or cheap).

### Getting It Online
- **NGINX Ingress**: We set up all the security headers and stuff
- **cert-manager**: Automatic HTTPS certificates because it's 2024
- **Custom domains**: `your-awesome-cms.com` instead of `localhost:8080`

### When Things Get Busy
- **Auto-scaling**: More traffic = more pods. Automatically.
- **Resource limits**: So one bad query doesn't take down your cluster
- **Shared storage**: Files work the same across all your pods

## Need the Technical Details?

If you're the type who reads the manual (respect! 📚), check out:
- [API Types](./api/v1/directus_types.go) - All the knobs and dials you can tweak
- [CRD Docs](./config/crd/) - The full Kubernetes resource specification

## Want to Help Out?

Found a bug? Got an idea? Want to make this thing even better? Awesome! Here's how:

### Easy Ways to Contribute
- 🐛 **Found a Bug?** Open an issue and tell us what broke
- 💡 **Got an Idea?** Feature requests are always welcome
- 📖 **Improve Docs** Make our examples even clearer
- 🧪 **Test Stuff** Try it on different clusters and let us know how it goes

### For the Code Warriors
- Fork it, branch it, code it, test it, PR it
- We're pretty chill about code style - just make it readable
- Tests are nice but don't stress if you're fixing a typo

### Getting Your Dev Environment Going
```bash
git clone <this-repo>
cd directus-operator
go mod download  # Get the dependencies
make test        # Make sure everything works
make run         # Run locally (connects to your current kubectl context)
```

Run `make help` to see all the available commands. There are... quite a few. 😅

### Shoutout to Kubebuilder
This whole thing is built with [Kubebuilder](https://book.kubebuilder.io/introduction.html), which is basically the best way to build Kubernetes operators. If you're into that sort of thing, check it out!

## License

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

