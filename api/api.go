package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"
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

// Versionz godoc
// @Summary Application version information
// @Description returns information about application version
// @Tags API
// @Produce json
// @Router /api/version [get]
// @Failure 500 {object} api.HTTPError
// @Success 200 {object} api.VersionResp
func (a *apiHandler) Versionz(w http.ResponseWriter, r *http.Request) {
	resp := VersionResp{
		AppName:        AppName,
		APIVersion:     APIVersion,
		AppVersion:     AppVersion,
		BuildTime:      BuildTime,
		LastCommitHash: LastCommitHash,
		LastCommitUser: LastCommitUser,
		LastCommitTime: LastCommitTime,
	}

	MustWriteJSON(a.l, w, r, resp, http.StatusOK)
}

// Healthz godoc
// @Summary Application health information
// @Description returns information whether application is up and running as well as information whether connection to DB is up. Endpoint returns http status 500 when there no connection to DB or 503 when service starts shutdown process
// @Tags API
// @Produce json
// @Router /api/health [get]
// @Failure 500 {object} api.HealthResp
// @Failure 503 {object} api.HTTPError
// @Success 200 {object} api.HealthResp
func (a *apiHandler) Healthz(w http.ResponseWriter, r *http.Request) {
	if atomic.LoadInt32(&a.healthy) != 1 {
		WriteErrJSON(a.l, w, r, errors.New("graceful shutdown started"), http.StatusServiceUnavailable)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp := HealthResp{Uptime: time.Since(StartTime).String()}
	ok, err := a.repo.OK(ctx)
	if err != nil {
		resp.DBError = err.Error()
	}
	resp.DBStatus = strconv.FormatBool(ok)

	if ok {
		MustWriteJSON(a.l, w, r, resp, http.StatusOK)
	} else {
		MustWriteJSON(a.l, w, r, resp, http.StatusInternalServerError)
	}
}

// Healthz godoc
// @Summary Application ready information
// @Description returns information whether application is ready for handling traffic. Endpoint returns http status 503 when service starts shutdown process
// @Tags API
// @Produce json
// @Router /api/ready [get]
// @Failure 500 {object} api.HTTPError
// @Failure 503 {object} api.HTTPError
// @Success 200 {object} api.HealthResp
func (a *apiHandler) Readyz(w http.ResponseWriter, r *http.Request) {
	if atomic.LoadInt32(&a.healthy) != 1 {
		WriteErrJSON(a.l, w, r, errors.New("graceful shutdown started"), http.StatusServiceUnavailable)
		return
	}

	MustWriteJSON(a.l, w, r, HealthResp{Msg: "OK"}, http.StatusOK)
}

// VersionResp - struct represents response for /version endpoint.
type VersionResp struct {
	AppName        string `json:"appName"`
	APIVersion     string `json:"apiVersion"`
	AppVersion     string `json:"version"`
	BuildTime      string `json:"buildTime"`
	LastCommitHash string `json:"gitCommitHash"`
	LastCommitUser string `json:"gitCommitUser"`
	LastCommitTime string `json:"gitCommitTime"`
}

// HealthResp - struct represents response for /health endpoint.
type HealthResp struct {
	Msg      string `json:"msg,omitempty"`
	Uptime   string `json:"uptime,omitempty"`
	DBStatus string `json:"dbStatus,omitempty"`
	DBError  string `json:"dbError,omitempty"`
}
