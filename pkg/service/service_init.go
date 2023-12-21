package service

import (
	"context"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/partyzanex/shortlink/pkg/logger"
	"github.com/partyzanex/shortlink/pkg/ptr"
	closer "github.com/partyzanex/shutdown"
)

func (s *Service) Init(options ...Option) (err error) {
	err = ErrInitialized

	s.init.Do(func() {
		err = s.initialize(options)
	})

	return err
}

func (s *Service) initialize(options []Option) (err error) {
	if s.GetStatus() > 0 {
		return ErrInitialize
	}

	s.startDebugServer(newConfig(options))
	s.setStatus(StatusInitialized)

	return nil
}

func (s *Service) startDebugServer(cfg *config) {
	if cfg.debugEnable == nil || !*cfg.debugEnable {
		return
	}

	if cfg.debugHost == nil {
		cfg.debugHost = ptr.Ptr(defaultDebugHost)
	}

	s.setDebugHost(*cfg.debugHost)

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Any("/version", s.versionHandler())
	router.Any("/live", s.liveness)
	router.Any("/ready", s.readiness)
	router.Any("/metrics", s.metricsHandler())
	pprof.Register(router, "/dev/pprof")

	if cfg.customRoutes != nil {
		cfg.customRoutes(router)
	}

	server := &http.Server{
		Addr:    *cfg.debugHost,
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			logger.GetLogger().WithError(err).Fatal("DebugServer")
		}
	}()

	s.closer.Append(s.stopServer(server))
}

func (*Service) stopServer(server *http.Server) closer.Fn {
	return func() error {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			return errors.Wrap(err, "cannot shutdown debug server")
		}

		return nil
	}
}

func (*Service) versionHandler() gin.HandlerFunc {
	var (
		ts = time.Now()
		v  = runtime.Version()
	)

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	return func(c *gin.Context) {
		c.JSON(http.StatusOK, &Version{
			Name:      "",
			Version:   "",
			Hostname:  hostname,
			StartedAt: ts,
			Uptime:    time.Since(ts),
			GoVersion: v,
		})
	}
}

type Version struct {
	Name      string        `json:"name"`
	Version   string        `json:"version"`
	Hostname  string        `json:"hostname"`
	StartedAt time.Time     `json:"started_at"`
	Uptime    time.Duration `json:"uptime"`
	GoVersion string        `json:"go_version"`
}

func (s *Service) liveness(c *gin.Context) {
	if status := s.GetStatus(); status == StatusFinished {
		c.Status(http.StatusServiceUnavailable)
	} else {
		c.Status(http.StatusOK)
	}
}

func (s *Service) readiness(c *gin.Context) {
	if s.GetStatus() != StatusStarted {
		c.Status(http.StatusServiceUnavailable)
	} else {
		c.Status(http.StatusOK)
	}
}

func (s *Service) metricsHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
