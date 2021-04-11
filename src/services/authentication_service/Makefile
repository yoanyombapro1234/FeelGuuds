# Makefile for releasing service
#
# The release version is controlled from pkg/version

TAG?=latest
NAME:=authentication_service

build-charts:
	helm lint charts/*
	helm package charts/*

down_mysql:
	echo "shutting down mysql backend database"
	docker-compose -f docker/docker-compose-mysql.yml down

down_postgres:
	echo "shutting down postgresql backend database"
	docker-compose -f docker/docker-compose-postgres.yml down

up_mysql:
	echo "starting authentication service with mysql backend database"
	docker-compose -f docker/docker-compose-mysql.yml up --remove-orphans

up_postgres:
	echo "starting authentication service with postgres backend database"
	docker-compose -f docker/docker-compose-postgres.yml up --remove-orphans
