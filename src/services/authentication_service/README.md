# FeelGuuds Service Template

[![e2e](https://github.com/yoanyombapro1234/Microservice-Template-Golang/workflows/e2e/badge.svg)](https://github.com/yoanyombapro1234/Microservice-Template-Golang/blob/master/.github/workflows/e2e.yml)
[![test](https://github.com/yoanyombapro1234/Microservice-Template-Golang/workflows/test/badge.svg)](https://github.com/yoanyombapro1234/Microservice-Template-Golang/blob/master/.github/workflows/test.yml)
[![cve-scan](https://github.com/yoanyombapro1234/Microservice-Template-Golang/workflows/cve-scan/badge.svg)](https://github.com/yoanyombapro1234/Microservice-Template-Golang/blob/master/.github/workflows/cve-scan.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/yoanyombapro1234/Microservice-Template-Golang)](https://goreportcard.com/report/github.com/yoanyombapro1234/Microservice-Template-Golang)
[![Docker Pulls](https://img.shields.io/docker/pulls/yoanyombapro1234/Microservice-Template-Golang)](https://hub.docker.com/r/yoanyombapro1234/Microservice-Template-Golang)

Please reference the keratin/authn [documentation](https://keratin.github.io/authn-server/#/guide-deploying_with_docker) for further details
specific to the endpoints this service exposes.

To run the service locally and its dependecies run the following

```
# to run the service with postgres db
cd src/services/authentication_service
make up-postgres

# to run the service with mysql db
cd src/services/authentication_service
make up-mysql

# to shut down the service and it dependencies (postgres)
cd src/services/authentication_service
make down-postgres

# to shut down the service and it dependencies (mysql)
cd src/services/authentication_service
make up-mysql
```
