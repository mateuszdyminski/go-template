# Go-template

TBD

## Getting Started

TBD

### Features

* 12-factor app compliant
* Inteligent health checks (readiness and liveness) - they are checking connection to DB as well
* Graceful shutdown on interrupt signals
* Instrumented with Prometheus
* Structured logging with zap
* Layered docker builds
* Multi-stage docker builds
* Repository for connecting PostgresDB
* Swagger docs available under `/swagger` endpoint

### Web API

* `GET` /version returns information about app version, last commiter, etc
* `GET` /metrics returns metrics for prometheus purpose
* `GET` /health returns liveness probe
* `GET` /ready returns readiness probe
* `GET` /swagger.json returns the API Swagger docs, used for Linkerd service profiling and Gloo routes discovery

### Prerequisites

You need to have working `go` environment:

* Install `go` - [https://golang.org/dl/](https://golang.org/dl/)
* Have working `make` - [https://www.gnu.org/software/make/](https://www.gnu.org/software/make/)
* Install `docker` - [https://docs.docker.com/install/](https://docs.docker.com/install/)
* Install `golangci-lint` - [https://github.com/golangci/golangci-lint](https://github.com/golangci/golangci-lint)
* Install `swag` - [https://github.com/swaggo/swag](https://github.com/swaggo/swag)
* Install `misspell` - [https://github.com/client9/misspell](https://github.com/client9/misspell)
* [For local development] Install `docker-compose` - [https://docs.docker.com/compose/install/](https://docs.docker.com/compose/install/)

### Installing

A step by step series of examples how to get a development env running:

Run `docker-compose` with all required components:

```bash
docker-compose up -d postgres
```

And run application locally:

```bash
make run
```

Or with docker-compose:

```bash
docker-compose up -d application
```

Now you can go to [http://localhost:8080/swagger/](http://localhost:8080/swagger) and check whether it's working.

## Running the tests

```bash
make test
```

### And coding style tests

```bash
make lint
```

## Deployment

To build, pack binary into Docker image and push it into [dockerhub.com](dockerhub.com):

```bash
make release
```

```bash
kubectl apply -k github.com/mateuszdyminski/go-template/kustomize
```

## Build final binary

```bash
make build
```

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/your/project/tags).

## Authors

TBD

## License

MIT

## Acknowledgments

TBD
