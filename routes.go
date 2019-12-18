package main

import (
	"context"
	"net/http"
	"net/http/pprof"

	"github.com/mateuszdyminski/go-template/api"
	"github.com/mateuszdyminski/go-template/app"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func newRouter(ctx context.Context, l *zap.Logger, repo app.Repository) *mux.Router {
	r := mux.NewRouter()

	// register Prometheus/Metrics middleware
	prom := api.NewMetricsMiddleware()
	r.Use(prom.Handler)

	// register request ID middleware
	r.Use(api.RequestIDMiddleware)

	// register logging middleware
	httpLogger := api.NewLoggingMiddleware(l)
	r.Use(httpLogger.Handler)

	// register version middleware
	r.Use(api.VersionMiddleware)

	apiHandler := api.NewAPIHandler(ctx, l, repo)

	r.HandleFunc("/api/version", apiHandler.Versionz).Methods(http.MethodGet)
	r.HandleFunc("/api/health", apiHandler.Healthz).Methods(http.MethodGet)
	r.HandleFunc("/api/ready", apiHandler.Readyz).Methods(http.MethodGet)

	r.Handle("/metrics", promhttp.Handler())

	return r
}

func newPprofRouter() *mux.Router {
	r := mux.NewRouter()

	// pprof endpoints configuration
	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)

	return r
}
