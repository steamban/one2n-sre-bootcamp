#!/bin/bash
set -e

echo "🎡 Spinning up a 4-node Minikube cluster..."
minikube start --nodes 4

echo "⚙️ Enabling Rancher storage provisioner..."
minikube addons enable storage-provisioner-rancher

echo "⏳ Waiting for nodes to be ready..."
kubectl wait --for=condition=ready nodes --all --timeout=120s

# Get node names dynamically
NODES=($(kubectl get nodes -o jsonpath='{.items[*].metadata.name}'))

echo "🚫 Tainting the control plane..."
kubectl taint nodes $(kubectl get nodes -l node-role.kubernetes.io/control-plane -o jsonpath='{.items[0].metadata.name}') node-role.kubernetes.io/control-plane:NoSchedule --overwrite

echo "⏳ Ensuring control-plane label exists..."
until kubectl get nodes -l node-role.kubernetes.io/control-plane | grep -q .; do
  sleep 2
done

echo "🏷️  Applying node labels to worker nodes..."
WORKER_NODES=($(kubectl get nodes -l '!node-role.kubernetes.io/control-plane' \
  -o jsonpath='{.items[*].metadata.name}' | tr ' ' '\n' | sort))

if [ ${#WORKER_NODES[@]} -ge 3 ]; then
  kubectl label nodes ${WORKER_NODES[0]} type=application --overwrite
  kubectl label nodes ${WORKER_NODES[1]} type=database --overwrite
  kubectl label nodes ${WORKER_NODES[2]} type=dependent_services --overwrite
else
  echo "❌ Error: Expected at least 3 worker nodes, found ${#WORKER_NODES[@]}"
  exit 1
fi

