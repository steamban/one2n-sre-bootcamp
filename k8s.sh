#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

MINIKUBE_NODES="${MINIKUBE_NODES:-4}"
MINIKUBE_DRIVER="${MINIKUBE_DRIVER:-docker}"
VAULT_NAMESPACE="${VAULT_NAMESPACE:-vault}"
EXTERNAL_SECRETS_NAMESPACE="${EXTERNAL_SECRETS_NAMESPACE:-external-secrets}"
APP_NAMESPACE="${APP_NAMESPACE:-student-api}"
ARGOCD_NAMESPACE="${ARGOCD_NAMESPACE:-argocd}"
DB_USER="${DB_USER:-studentapp}"
DB_PASSWORD="${DB_PASSWORD:-complicated}"
VAULT_TOKEN="${VAULT_TOKEN:-root}"

log() {
  echo "==> $1"
}

ensure_minikube_running() {
  local host_status
  host_status="$(minikube status --format='{{.Host}}' 2>/dev/null || true)"
  if [[ "$host_status" == "Running" ]]; then
    log "Minikube cluster already running, skipping start"
    return
  fi
  log "Starting Minikube cluster with $MINIKUBE_NODES nodes"
  minikube start --nodes "$MINIKUBE_NODES" --driver="$MINIKUBE_DRIVER"
  minikube addons enable storage-provisioner-rancher
}

label_node() {
  local node="$1"
  local role="$2"
  local current
  current="$(kubectl get node "$node" -o jsonpath='{.metadata.labels.type}' 2>/dev/null || true)"
  if [[ "$current" == "$role" ]]; then
    return
  fi
  kubectl label nodes "$node" type="$role" --overwrite
}

ensure_node_labels() {
  log "Ensuring node labels"
  local control_plane
  control_plane="$(kubectl get nodes -l node-role.kubernetes.io/control-plane -o jsonpath='{.items[0].metadata.name}' 2>/dev/null || true)"

  if [[ -z "$control_plane" ]]; then
    until kubectl get nodes -l node-role.kubernetes.io/control-plane | grep -q .; do
      sleep 2
    done
    control_plane="$(kubectl get nodes -l node-role.kubernetes.io/control-plane -o jsonpath='{.items[0].metadata.name}')"
  fi

  kubectl taint nodes "$control_plane" node-role.kubernetes.io/control-plane:NoSchedule --overwrite

  local workers=()
  while IFS= read -r node; do
    [[ -n "$node" ]] && workers+=("$node")
  done < <(kubectl get nodes -l '!node-role.kubernetes.io/control-plane' -o jsonpath='{range .items[*]}{.metadata.name}{"\n"}{end}')

  if [[ ${#workers[@]} -ge 3 ]]; then
    label_node "${workers[0]}" application
    label_node "${workers[1]}" database
    label_node "${workers[2]}" dependent_services
  else
    log "Warning: Expected at least 3 worker nodes, found ${#workers[@]}"
  fi
}

ensure_namespaces() {
  for ns in "$VAULT_NAMESPACE" "$EXTERNAL_SECRETS_NAMESPACE" "$ARGOCD_NAMESPACE" "$APP_NAMESPACE"; do
    kubectl get ns "$ns" &>/dev/null || kubectl create ns "$ns"
  done
}

ensure_helm_repo() {
  local name="$1"
  local url="$2"
  if helm repo list | awk '{print $1}' | grep -qx "$name"; then
    return
  fi
  log "Adding Helm repo $name"
  helm repo add "$name" "$url"
}

ensure_helm_repos() {
  ensure_helm_repo hashicorp https://helm.releases.hashicorp.com
  ensure_helm_repo external-secrets https://charts.external-secrets.io
  ensure_helm_repo argo https://argoproj.github.io/argo-helm
  log "Updating Helm repos"
  helm repo update
}

wait_for_eso_crds() {
  log "Waiting for External Secrets CRDs"
  kubectl wait --for=condition=established crd/externalsecrets.external-secrets.io --timeout=300s
  kubectl wait --for=condition=established crd/secretstores.external-secrets.io --timeout=300s
}

seed_vault() {
  log "Checking Vault seed data"
  if kubectl exec vault-0 -n "$VAULT_NAMESPACE" -- vault kv get secret/student-api/db >/dev/null 2>&1; then
    log "Vault secret student-api/db already exists, skipping seed"
    return
  fi
  log "Seeding Vault secret student-api/db"
  kubectl exec vault-0 -n "$VAULT_NAMESPACE" -- vault login "$VAULT_TOKEN" >/dev/null
  kubectl exec vault-0 -n "$VAULT_NAMESPACE" -- vault kv put secret/student-api/db username="$DB_USER" password="$DB_PASSWORD" dbname="student_db"
}

deploy_core() {
  ensure_namespaces
  ensure_helm_repos

  log "Deploying Vault"
  helm upgrade --install vault hashicorp/vault \
    --namespace "$VAULT_NAMESPACE" \
    -f "$ROOT_DIR/infra/helm/vault-values.yaml" \
    --wait --timeout=5m
  kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=vault -n "$VAULT_NAMESPACE" --timeout=300s
  seed_vault

  log "Deploying External Secrets Operator"
  helm upgrade --install external-secrets external-secrets/external-secrets \
    --namespace "$EXTERNAL_SECRETS_NAMESPACE" \
    -f "$ROOT_DIR/infra/helm/external-secrets-values.yaml" \
    --wait --timeout=5m
  kubectl wait --for=condition=ready pod -l app.kubernetes.io/instance=external-secrets -n "$EXTERNAL_SECRETS_NAMESPACE" --timeout=300s
  wait_for_eso_crds
}

deploy_gitops() {
  ensure_namespaces
  ensure_helm_repos

  log "Deploying ArgoCD"
  helm upgrade --install argocd argo/argo-cd \
    --namespace "$ARGOCD_NAMESPACE" \
    -f "$ROOT_DIR/infra/helm/argocd-values.yaml" \
    --wait --timeout=5m
  kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=argocd-server -n "$ARGOCD_NAMESPACE" --timeout=300s

  log "Applying ArgoCD Application manifests"
  kubectl apply -f "$ROOT_DIR/apps/postgres.yaml"
  kubectl apply -f "$ROOT_DIR/apps/student-api.yaml"
}

deploy_cluster() {
  ensure_minikube_running
  ensure_node_labels
}

deploy_all() {
  deploy_cluster
  deploy_core
  deploy_gitops
  log "Deployment complete. Run the following to access ArgoCD:"
  echo "  kubectl port-forward svc/argocd-server -n argocd 8080:443"
  echo "  ArgoCD password: kubectl get secret argocd-initial-admin-secret -n argocd -o jsonpath='{.data.password}' | base64 -d"
}

main() {
  local command="${1:-up}"
  case "$command" in
    up)         deploy_all ;;
    cluster)    deploy_cluster ;;
    core)       deploy_core ;;
    gitops)     deploy_gitops ;;
    *)
      echo "Usage: $0 {up|cluster|core|gitops}" >&2
      exit 1
      ;;
  esac
}

main "$@"
