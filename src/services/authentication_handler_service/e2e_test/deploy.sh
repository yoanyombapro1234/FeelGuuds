#! /usr/bin/env sh

set -e

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd -P)

# add jetstack repository
helm repo add jetstack https://charts.jetstack.io || true

# install cert-manager
helm upgrade --install cert-manager jetstack/cert-manager \
    --set installCRDs=true \
    --namespace default

# wait for cert manager
kubectl rollout status deployment/cert-manager --timeout=4m
kubectl rollout status deployment/cert-manager-webhook --timeout=4m
kubectl rollout status deployment/cert-manager-cainjector --timeout=4m

# install self-signed certificate
cat << 'EOF' | kubectl apply -f -
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: self-signed
spec:
  selfSigned: {}
EOF

# deploy service dependencies
$SCRIPT_DIR/start_dep.sh

# install service with tls enabled
helm upgrade --install service ./charts/authentication_handler_service \
    --set image.repository=feelguuds/authentication_handler_service \
    --set image.tag=latest \
    --set tls.enabled=true \
    --set certificate.create=true \
    --namespace=default
