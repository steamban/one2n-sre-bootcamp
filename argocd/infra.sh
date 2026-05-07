#!/bin/bash
set -e

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🎡 Full Infrastructure Setup"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Step 1: Spin up 4-node Minikube cluster
echo ""
echo "🎡 Spinning up a 4-node Minikube cluster..."
minikube start --nodes 4

echo "⚙️ Enabling Rancher storage provisioner..."
minikube addons enable storage-provisioner-rancher

# Step 2: Wait for nodes to be ready
echo ""
echo "⏳ Waiting for nodes to be ready..."
kubectl wait --for=condition=ready nodes --all --timeout=120s

# Step 3: Get all node names
NODES=($(kubectl get nodes -o jsonpath='{.items[*].metadata.name}'))
echo "✓ Found ${#NODES[@]} node(s)"

# Step 4: Taint the control plane
echo ""
echo "🚫 Tainting the control plane..."
CONTROL_PLANE=$(kubectl get nodes -l node-role.kubernetes.io/control-plane -o jsonpath='{.items[0].metadata.name}')
kubectl taint nodes "$CONTROL_PLANE" node-role.kubernetes.io/control-plane:NoSchedule --overwrite
echo "✓ Control plane tainted with NoSchedule"

# Step 5: Ensure control-plane label exists
echo ""
echo "⏳ Ensuring control-plane label exists..."
until kubectl get nodes -l node-role.kubernetes.io/control-plane | grep -q .; do
  sleep 2
done

# Step 6: Label worker nodes
echo ""
echo "🏷️  Applying node labels to worker nodes..."
WORKER_NODES=($(kubectl get nodes -l '!node-role.kubernetes.io/control-plane' \
  -o jsonpath='{.items[*].metadata.name}' | tr ' ' '\n' | sort))

if [ ${#WORKER_NODES[@]} -ge 3 ]; then
  kubectl label nodes ${WORKER_NODES[0]} type=application --overwrite
  echo "✓ Labeled ${WORKER_NODES[0]} as type=application"

  kubectl label nodes ${WORKER_NODES[1]} type=database --overwrite
  echo "✓ Labeled ${WORKER_NODES[1]} as type=database"

  kubectl label nodes ${WORKER_NODES[2]} type=dependent_services --overwrite
  echo "✓ Labeled ${WORKER_NODES[2]} as type=dependent_services"
else
  echo "❌ Error: Expected at least 3 worker nodes, found ${#WORKER_NODES[@]}"
  exit 1
fi

# Step 7: Add ArgoCD Helm repo if not already added
echo ""
echo "📦 Adding ArgoCD Helm repository..."
if ! helm repo list 2>/dev/null | grep -q "argo"; then
  helm repo add argo https://argoproj.github.io/argo-helm
  echo "✓ Added argo Helm repo"
else
  echo "✓ ArgoCD Helm repo already present"
fi

helm repo update
echo "✓ Updated Helm repos"

# Step 8: Create argocd namespace
echo ""
echo "📍 Creating argocd namespace..."
kubectl create namespace argocd --dry-run=client -o yaml | kubectl apply -f -
echo "✓ Namespace argocd ready"

# Step 9: Install ArgoCD via Helm
echo ""
echo "🚀 Installing ArgoCD..."
helm install argocd argo/argo-cd \
  --namespace argocd \
  --set server.service.type=NodePort \
  --set server.service.nodePortHttp=30880 \
  --set server.service.nodePortHttps=30443 \
  --set installCRDs=true \
  --wait --timeout=5m

echo "✓ ArgoCD installed successfully"

# Step 10: Get ArgoCD credentials
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ Infrastructure setup complete!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "🔑 ArgoCD Initial Admin Password:"
echo "   kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath='{.data.password}' | base64 -d"
echo ""
echo "🌐 Access ArgoCD UI:"
echo "   minikube service argo-cd-argocd-server -n argocd"
echo "   Or port-forward: kubectl port-forward svc/argo-cd-argocd-server -n argocd 8080:443"
echo ""
echo "📦 Create ArgoCD Applications:"
echo "   kubectl apply -f argocd/apps/"
echo ""