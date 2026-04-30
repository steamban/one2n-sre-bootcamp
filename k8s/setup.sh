#!/bin/bash
set -e

echo "🎡 Spinning up a 4-node Minikube cluster..."
minikube start --nodes 4

echo "⏳ Waiting for nodes to be ready..."
kubectl wait --for=condition=ready nodes --all --timeout=120s

# Get node names dynamically
NODES=($(kubectl get nodes -o jsonpath='{.items[*].metadata.name}'))

echo "🚫 Tainting the control plane..."
kubectl taint nodes ${NODES[0]} node-role.kubernetes.io/control-plane:NoSchedule --overwrite

echo "🏷️  Applying node labels to worker nodes..."
# Node 1 -> application | Node 2 -> database | Node 3 -> dependent_services
kubectl label node ${NODES[1]} type=application --overwrite
kubectl label node ${NODES[2]} type=database --overwrite
kubectl label node ${NODES[3]} type=dependent_services --overwrite

echo "⚙️  Installing External Secrets CRDs..."
kubectl apply -f https://raw.githubusercontent.com/external-secrets/external-secrets/main/deploy/crds/bundle.yaml --server-side

# Wait for CRDs to be established
echo "⏳ Waiting for External Secrets CRDs to be established..."
kubectl wait --for=condition=Established crd/externalsecrets.external-secrets.io --timeout=60s

echo "🚀 Starting One-Click Deployment..."

# 1. Create Namespaces
kubectl apply -f 0-namespace.yml

# 2. Deploy Infrastructure
kubectl apply -f 3-vault.yml
kubectl apply -f 4-external-secrets.yml
kubectl apply -f student-config.yml

echo "⏳ Waiting for Vault to be ready..."
kubectl wait --for=condition=ready pod -l app=vault -n vault --timeout=90s

# 3. Auto-Inject Secrets into Vault
echo "🔑 Seeding Vault with Database Credentials..."
kubectl exec -n vault deployments/vault -- sh -c \
  "VAULT_ADDR='http://127.0.0.1:8200' VAULT_TOKEN='root' vault kv put secret/student-api/db username='postgres' password='supersecretpassword' dbname='student_db'"

echo "⏳ Waiting for External Secret to sync..."
kubectl wait --for=condition=ready externalsecret/student-api-db-creds -n student-api --timeout=60s

# 4. Deploy Database and Application
echo "📦 Deploying Postgres and Student API..."
kubectl apply -f 1-database.yml
kubectl apply -f 2-application.yml

echo "⏳ Waiting for the Student API to be up and healthy..."
kubectl wait --for=condition=ready pod -l app=student-api -n student-api --timeout=120s

echo "✅ Setup Complete!"
echo "🌐 API available on port 8080"

kubectl port-forward -n student-api svc/student-api 8080:8080