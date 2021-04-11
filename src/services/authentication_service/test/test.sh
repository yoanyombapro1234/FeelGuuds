#1 /usr/bin/env sh

set -e

# wait for service
kubectl rollout status deployment/authentication_service --timeout=3m

# test service
helm test service
