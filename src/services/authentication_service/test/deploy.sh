#! /usr/bin/env sh

# add jetstack repository
helm repo add jetstack https://charts.jetstack.io || true

# install cert-manager
helm upgrade --install cert-manager jetstack/cert-manager \
    --set installCRDs=true \
    --namespace default

# wait for cert manager
kubectl rollout status deployment/cert-manager --timeout=2m
kubectl rollout status deployment/cert-manager-webhook --timeout=2m
kubectl rollout status deployment/cert-manager-cainjector --timeout=2m

# install self-signed certificate
cat << 'EOF' | kubectl apply -f -
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: self-signed
spec:
  selfSigned: {}
EOF

# install service with tls enabled
helm upgrade --install authentication_service ./charts/authentication_service \
    --set image.repository=test/authentication_service \
    --set image.tag=1.10.4 \
    --set tls.enabled=true \
    --set certificate.create=true \
    --namespace=default
