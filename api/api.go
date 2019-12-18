package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/mateuszdyminski/go-template/app"
	"go.uber.org/zap"
)

// AppName - variable injected by -X flag during build
var AppName = "unknown"

// AppVersion - variable injected by -X flag during build
var AppVersion = "unknown"

// APIVersion - variable injected by -X flag during build
var APIVersion = "unknown"

// LastCommitTime - variable injected by -X flag during build
var LastCommitTime = "unknown"

// LastCommitUser - variable injected by -X flag during build
var LastCommitUser = "unknown"

// LastCommitHash - variable injected by -X flag during build
var LastCommitHash = "unknown"

// BuildTime - variable injected by -X flag during build
var BuildTime = "unknown"

// StartTime - time(int UTC) when application starts
var StartTime = time.Now().UTC()

// ApiHandler general handler responsible for providing information about app itself.
// Available handlers: Versionz, Readyz, Healthz
type ApiHandler interface {
	Versionz(http.ResponseWriter, *http.Request)
	Readyz(http.ResponseWriter, *http.Request)
	Healthz(http.ResponseWriter, *http.Request)
}

type apiHandler struct {
	l       *zap.SugaredLogger
	repo    app.Repository
	healthy int32
}

func NewAPIHandler(ctx context.Context, l *zap.Logger, repo app.Repository) ApiHandler {
	a := &apiHandler{l: l.Sugar(), repo: repo, healthy: 1}
	go a.watchSignals(ctx)
	return a
}

func (a *apiHandler) watchSignals(ctx context.Context) {
	<-ctx.Done()
	atomic.StoreInt32(&a.healthy, 0)
}

func (a *apiHandler) Versionz(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{
		"appName":       AppName,
		"version":       AppVersion,
		"buildTime":     BuildTime,
		"gitCommitHash": LastCommitHash,
		"gitCommitUser": LastCommitUser,
		"gitCommitTime": LastCommitTime,
	}

	MustWriteJSON(a.l, w, r, resp, http.StatusOK)
}

func (a *apiHandler) Healthz(w http.ResponseWriter, r *http.Request) {
	if atomic.LoadInt32(&a.healthy) != 1 {
		WriteErrJSON(a.l, w, r, errors.New("graceful shutdown started"), http.StatusServiceUnavailable)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	dbResp := "up"
	ok, err := a.repo.OK(ctx)
	if err != nil {
		dbResp = fmt.Sprintf("db down - %s", err)
	}

	resp := map[string]string{
		"uptime":   time.Since(StartTime).String(),
		"dbStatus": dbResp,
	}

	if ok {
		MustWriteJSON(a.l, w, r, resp, http.StatusOK)
	} else {
		MustWriteJSON(a.l, w, r, resp, http.StatusInternalServerError)
	}
}

func (a *apiHandler) Readyz(w http.ResponseWriter, r *http.Request) {
	if atomic.LoadInt32(&a.healthy) != 1 {
		WriteErrJSON(a.l, w, r, errors.New("graceful shutdown started"), http.StatusServiceUnavailable)
		return
	}

	MustWriteJSON(a.l, w, r, struct{ Msg string }{"OK"}, http.StatusOK)
}
