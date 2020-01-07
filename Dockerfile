# Arguments
ARG GOLANG_VERSION
ARG ALPINE_VERSION

# Build part
FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} AS builder

ARG SWAG_VERSION=1.6.3

RUN apk --no-cache add make git tar; \
    adduser -D -h /tmp/build build; \
    wget -qO- https://github.com/swaggo/swag/releases/download/v${SWAG_VERSION}/swag_${SWAG_VERSION}_Linux_x86_64.tar.gz | tar -xzf - -C /tmp; \
    mv /tmp/swag /usr/bin/swag
USER build
WORKDIR /tmp/build

COPY --chown=build Makefile Makefile
COPY --chown=build go.mod go.mod
COPY --chown=build go.sum go.sum

RUN go mod download

ARG VERSION
ARG NAME
ARG BUILD_TIME
ARG LAST_COMMIT_USER
ARG LAST_COMMIT_HASH
ARG LAST_COMMIT_TIME

COPY --chown=build api api
COPY --chown=build app app
COPY --chown=build repository repository
COPY --chown=build config.go config.go
COPY --chown=build routes.go routes.go
COPY --chown=build main.go main.go
RUN make swag
RUN make build

# Exec part
FROM gcr.io/distroless/base

ARG VERSION
ARG NAME

# Appication Configuration
ENV DEBUG=""
ENV APP_HTTP_PORT="8080"
ENV APP_HTTP_PPROF_PORT="8090"
ENV APP_HTTP_GRACEFUL_TIMEOUT="5"
ENV APP_HTTP_GRACEFUL_SLEEP="1"
ENV APP_POSTGRES_HOST="postgres"
ENV APP_POSTGRES_PORT="5432"
ENV APP_POSTGRES_USER="postgres"
ENV APP_POSTGRES_PASSWORD="password"
ENV APP_POSTGRES_DBNAME="app_db"

# Copy from builder
COPY --from=builder /tmp/build/${NAME}-${VERSION} /usr/bin/app

# Exec
CMD ["app"]
