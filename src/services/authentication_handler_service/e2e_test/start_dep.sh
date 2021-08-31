#! /usr/bin/env sh
# deploy service dependencies
# ref. to test connectivity - https://www.bmc.com/blogs/kubernetes-postgresql/

# install postgres as a stateful set
kubectl apply -f ./charts/service_dependencies/postgres/config.yaml
kubectl apply -f ./charts/service_dependencies/postgres/service.yaml
kubectl apply -f ./charts/service_dependencies/postgres/deployment.yaml

# install redis as a stateful set
kubectl apply -f ./charts/service_dependencies/redis/config.yaml
kubectl apply -f ./charts/service_dependencies/redis/service.yaml
kubectl apply -f ./charts/service_dependencies/redis/deployment.yaml

# install the authentication service
# ensure it sits behind a load balancer and is comprised of at least 3 replicas
kubectl apply -f ./charts/service_dependencies/authentication_service/service.yaml
kubectl apply -f ./charts/service_dependencies/authentication_service/deployment.yaml
