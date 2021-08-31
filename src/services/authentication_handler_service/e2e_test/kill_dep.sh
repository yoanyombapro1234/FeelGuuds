#! /usr/bin/env sh

kubectl delete -f ./charts/service_dependencies/postgres/config.yaml
kubectl delete -f ./charts/service_dependencies/postgres/service.yaml
kubectl delete -f ./charts/service_dependencies/postgres/deployment.yaml

kubectl delete -f ./charts/service_dependencies/redis/config.yaml
kubectl delete -f ./charts/service_dependencies/redis/service.yaml
kubectl delete -f ./charts/service_dependencies/redis/deployment.yaml

kubectl delete -f ./charts/service_dependencies/authentication_service/service.yaml
kubectl delete -f ./charts/service_dependencies/authentication_service/deployment.yaml
