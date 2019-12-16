# Arguments
ARG GOLANG_VERSION
ARG ALPINE_VERSION

# Build part
FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} AS builder

RUN apk --no-cache add make git; \
    adduser -D -h /tmp/build build
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

COPY --chown=build version version
COPY --chown=build main.go main.go
RUN make build

# Exec part
FROM gcr.io/distroless/base

ARG VERSION
ARG NAME

# Appication Configuration
ENV PORT="8080"

# Copy from builder
COPY --from=builder /tmp/build/${NAME}-${VERSION} /usr/bin/app

# Exec
CMD ["app"]