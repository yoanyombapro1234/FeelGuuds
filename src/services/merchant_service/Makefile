# Makefile for releasing service
#
# The release version is controlled from pkg/version

TAG?=latest
NAME:=service
DOCKER_REPOSITORY:=stefanprodan
DOCKER_IMAGE_NAME:=$(DOCKER_REPOSITORY)/$(NAME)
GIT_COMMIT:=$(shell git describe --dirty --always)
VERSION:=$(shell grep 'VERSION' pkg/version/version.go | awk '{ print $$4 }' | tr -d '"')
EXTRA_RUN_ARGS?=

run:
	go run -ldflags "-s -w -X github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/version.REVISION=$(GIT_COMMIT)" cmd/service/* \
	--level=debug --grpc-port=9999 --backend-url=https://httpbin.org/status/401 --backend-url=https://httpbin.org/status/500 \
	--ui-logo=https://raw.githubusercontent.com/stefanprodan/service/gh-pages/cuddle_clap.gif $(EXTRA_RUN_ARGS)

test:
	go test ./... -coverprofile cover.out

build:
	GIT_COMMIT=$$(git rev-list -1 HEAD) && CGO_ENABLED=0 go build  -ldflags "-s -w -X github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/version.REVISION=$(GIT_COMMIT)" -a -o ./bin/service ./cmd/service/*
	GIT_COMMIT=$$(git rev-list -1 HEAD) && CGO_ENABLED=0 go build  -ldflags "-s -w -X github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/version.REVISION=$(GIT_COMMIT)" -a -o ./bin/servicecli ./cmd/servicecli/*

fmt:
	gofmt -l -s -w ./
	goimports -l -w ./

build-charts:
	helm lint charts/*
	helm package charts/*

build-container:
	docker build -t $(DOCKER_IMAGE_NAME):$(VERSION) .

build-base:
	docker build -f Dockerfile.base -t $(DOCKER_REPOSITORY)/service-base:latest .

push-base: build-base
	docker push $(DOCKER_REPOSITORY)/service-base:latest

test-container:
	@docker rm -f service || true
	@docker run -dp 9898:9898 --name=service $(DOCKER_IMAGE_NAME):$(VERSION)
	@docker ps
	@TOKEN=$$(curl -sd 'test' localhost:9898/token | jq -r .token) && \
	curl -sH "Authorization: Bearer $${TOKEN}" localhost:9898/token/validate | grep test

push-container:
	docker tag $(DOCKER_IMAGE_NAME):$(VERSION) $(DOCKER_IMAGE_NAME):latest
	docker push $(DOCKER_IMAGE_NAME):$(VERSION)
	docker push $(DOCKER_IMAGE_NAME):latest
	docker tag $(DOCKER_IMAGE_NAME):$(VERSION) quay.io/$(DOCKER_IMAGE_NAME):$(VERSION)
	docker tag $(DOCKER_IMAGE_NAME):$(VERSION) quay.io/$(DOCKER_IMAGE_NAME):latest
	docker push quay.io/$(DOCKER_IMAGE_NAME):$(VERSION)
	docker push quay.io/$(DOCKER_IMAGE_NAME):latest

version-set:
	@next="$(TAG)" && \
	current="$(VERSION)" && \
	sed -i '' "s/$$current/$$next/g" pkg/version/version.go && \
	sed -i '' "s/tag: $$current/tag: $$next/g" charts/service/values.yaml && \
	sed -i '' "s/tag: $$current/tag: $$next/g" charts/service/values-prod.yaml && \
	sed -i '' "s/appVersion: $$current/appVersion: $$next/g" charts/service/Chart.yaml && \
	sed -i '' "s/version: $$current/version: $$next/g" charts/service/Chart.yaml && \
	sed -i '' "s/service:$$current/service:$$next/g" kustomize/deployment.yaml && \
	sed -i '' "s/service:$$current/service:$$next/g" deploy/webapp/frontend/deployment.yaml && \
	sed -i '' "s/service:$$current/service:$$next/g" deploy/webapp/backend/deployment.yaml && \
	sed -i '' "s/service:$$current/service:$$next/g" deploy/bases/frontend/deployment.yaml && \
	sed -i '' "s/service:$$current/service:$$next/g" deploy/bases/backend/deployment.yaml && \
	echo "Version $$next set in code, deployment, chart and kustomize"

release:
	git tag $(VERSION)
	git push origin $(VERSION)

swagger:
	go get github.com/swaggo/swag/cmd/swag
	cd pkg/api && $$(go env GOPATH)/bin/swag init -g server.go
