FROM golang:1.17-alpine as builder

RUN apk update && apk upgrade -U -a

ARG REVISION

RUN mkdir -p /service/

WORKDIR /service

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -ldflags "-s -w \
	-X github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/version.REVISION=${REVISION}" \
	-a -o bin/service cmd/service/*

RUN CGO_ENABLED=0 go build -ldflags "-s -w \
	-X github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/version.REVISION=${REVISION}" \
	-a -o bin/servicecli cmd/servicecli/*

FROM alpine:3.14

ARG BUILD_DATE
ARG VERSION
ARG REVISION

LABEL maintainer="yoanyomba"

RUN addgroup -S app \
	&& adduser -S -G app app \
	&& apk --no-cache add \
	ca-certificates curl netcat-openbsd

WORKDIR /home/app

COPY --from=builder /service/bin/service .
COPY --from=builder /service/bin/servicecli /usr/local/bin/servicecli
COPY ./ui ./ui
RUN chown -R app:app ./

USER app

CMD ["./service"]
