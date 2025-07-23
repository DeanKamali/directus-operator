# Directus Operator Examples - Simple Makefile
# Just the commands you actually need

.PHONY: help basic production ingress autoscaling status clean

help: ## Show this help
	@echo "🚀 Directus Operator Examples"
	@echo ""
	@echo "Deploy:"
	@echo "  make basic       - SQLite deployment (good for testing)"
	@echo "  make production  - PostgreSQL + Redis (production ready)"
	@echo "  make ingress     - With HTTPS (needs ingress controller)"
	@echo "  make autoscaling - Auto-scaling setup (needs metrics-server)"
	@echo ""
	@echo "Manage:"
	@echo "  make status      - Check what's running"
	@echo "  make access      - Access basic deployment"
	@echo "  make clean       - Remove everything"
	@echo ""
	@echo "Need help? Each example has its own README with details."

basic: ## Deploy basic SQLite example
	@echo "🔨 Deploying basic Directus..."
	kubectl apply -f basic/directus.yaml
	@echo "✅ Done! Access with: make access"

production: ## Deploy production example with PostgreSQL and Redis
	@echo "🏭 Deploying production setup..."
	kubectl apply -f production/dependencies/
	@echo "⏳ Waiting for database..."
	kubectl wait --for=condition=ready pod -l app=postgresql --timeout=300s
	kubectl wait --for=condition=ready pod -l app=redis --timeout=300s
	kubectl apply -f production/directus.yaml
	@echo "✅ Production deployment started!"

ingress: ## Deploy with HTTPS (update domain in with-ingress/directus.yaml first!)
	@echo "🌍 Deploying with HTTPS..."
	@echo "⚠️  Update the domain in with-ingress/directus.yaml first!"
	kubectl apply -f with-ingress/cert-issuer.yaml
	kubectl apply -f with-ingress/directus.yaml
	@echo "✅ Ingress deployment started!"

autoscaling: ## Deploy with auto-scaling
	@echo "📈 Deploying auto-scaling setup..."
	@kubectl get apiservice v1beta1.metrics.k8s.io >/dev/null 2>&1 || { echo "❌ Need metrics-server: kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml"; exit 1; }
	kubectl apply -f autoscaling/dependencies/
	@echo "⏳ Waiting for dependencies..."
	kubectl wait --for=condition=ready pod -l app=postgresql --timeout=300s
	kubectl wait --for=condition=ready pod -l app=redis --timeout=300s
	kubectl apply -f autoscaling/directus.yaml
	@echo "✅ Autoscaling deployment started!"

status: ## Check what's running
	@echo "📊 Current deployments:"
	@kubectl get directus 2>/dev/null || echo "   No Directus instances found"
	@echo ""
	@kubectl get hpa 2>/dev/null && echo "" || true

access: ## Access basic deployment (port-forward + show password)
	@echo "🔑 Admin password:"
	@kubectl get secret basic-directus-secrets -o jsonpath='{.data.ADMIN_PASSWORD}' 2>/dev/null | base64 -d || echo "Basic deployment not found"
	@echo ""
	@echo "🌐 Starting port forward to http://localhost:8080"
	@echo "   (Press Ctrl+C to stop)"
	kubectl port-forward svc/basic-directus 8080:80

clean: ## Remove all deployments
	@echo "🗑️  Removing all deployments..."
	kubectl delete -f basic/directus.yaml --ignore-not-found=true
	kubectl delete -f production/directus.yaml --ignore-not-found=true
	kubectl delete -f production/dependencies/ --ignore-not-found=true
	kubectl delete -f with-ingress/directus.yaml --ignore-not-found=true
	kubectl delete -f with-ingress/cert-issuer.yaml --ignore-not-found=true
	kubectl delete -f autoscaling/directus.yaml --ignore-not-found=true
	kubectl delete -f autoscaling/dependencies/ --ignore-not-found=true
	@echo "✅ Clean up complete!" 