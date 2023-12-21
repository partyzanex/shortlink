package service

import (
	"context"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/partyzanex/shortlink/pkg/logger"
	"github.com/partyzanex/shortlink/pkg/ptr"
	"github.com/partyzanex/shutdown"
)

func (s *Service) Start(options ...Option) (err error) {
	err = ErrServiceIsStarted

	s.start.Do(func() {
		err = s.startService(options)
	})

	return err
}

func (s *Service) startService(options []Option) (err error) {
	if status := s.GetStatus(); status != StatusInitialized {
		return errors.Errorf("cannot start service, current status is %q", status)
	}

	defer func() {
		if err == nil {
			s.setStatus(StatusStarted)
		}
	}()

	s.closer.Append(shutdown.Fn(func() error {
		s.setStatus(StatusFinished)

		return nil
	}))

	cfg := newConfig(options)

	if err = s.startGRPCServer(cfg); err != nil {
		return errors.Wrap(err, "cannot start GRPC server")
	}

	gw, err := s.startGateway(cfg)
	if err != nil {
		return errors.Wrap(err, "cannot start gateway")
	}

	if err = s.startHTTPServer(cfg, gw); err != nil {
		return errors.Wrap(err, "cannot start HTTP server")
	}

	return nil
}

func (s *Service) startGateway(cfg *config) (gin.HandlerFunc, error) {
	if cfg.GRPCEnable == nil || !*cfg.GRPCEnable {
		return nil, nil
	}

	if cfg.GRPCHost == nil && *cfg.GRPCHost == "" {
		return nil, errors.New("gRPC server is not configured")
	}

	if len(cfg.gatewayHandlers) == 0 {
		return nil, errors.New("required GatewayHandler")
	}

	ctx, cancel := context.WithCancel(context.Background())

	s.closer.Append(shutdown.Fn(func() error {
		cancel()

		return nil
	}))

	gw := runtime.NewServeMux(cfg.gatewayMuxOptions...)

	for _, handler := range cfg.gatewayHandlers {
		err := handler(ctx, gw, *cfg.GRPCHost, cfg.gatewayDialOptions)
		if err != nil {
			return nil, errors.Wrap(err, "cannot set gateway handler")
		}
	}

	return func(c *gin.Context) {
		c.Status(http.StatusOK)
		gw.ServeHTTP(c.Writer, c.Request)
	}, nil
}

func (s *Service) startGRPCServer(cfg *config) error {
	if cfg.GRPCEnable == nil || !*cfg.GRPCEnable {
		return nil
	}

	if cfg.GRPCHost == nil {
		cfg.GRPCHost = ptr.Ptr(defaultGRPCHost)
	}

	server := grpc.NewServer(cfg.GRPCServerOptions...)

	if cfg.GRPCHandler != nil {
		cfg.GRPCHandler(server)
	}

	if cfg.GRPCReflection != nil && *cfg.GRPCReflection {
		reflection.Register(server)
	}

	lis, err := net.Listen("tcp", *cfg.GRPCHost)
	if err != nil {
		return errors.Wrapf(err, "cannot listen address %q", *cfg.GRPCHost)
	}

	go func() {
		if errServe := server.Serve(lis); errServe != nil {
			logger.GetLogger().WithError(errServe).Error("GRPCServer")
		}
	}()

	logger.GetLogger().Warn("gRPC server started on ", *cfg.GRPCHost)

	s.closer.Append(shutdown.Fn(func() error {
		server.GracefulStop()

		return nil
	}))

	return nil
}

func (s *Service) startHTTPServer(cfg *config, gw gin.HandlerFunc) error {
	if cfg.HTTPEnable == nil || !*cfg.HTTPEnable {
		return nil
	}

	h := gin.New()
	h.Use(cfg.HTTPMiddlewares...)

	if gw != nil {
		h.NoRoute(gw)
	}

	if cfg.HTTPCustomRoutes != nil {
		cfg.HTTPCustomRoutes(h.Group("/"))
	}

	if cfg.HTTPHost == nil {
		cfg.HTTPHost = ptr.Ptr(defaultHTTPHost)
	}

	server := &http.Server{
		Addr:    *cfg.HTTPHost,
		Handler: h,
	}

	go func() {
		if errServe := server.ListenAndServe(); errServe != nil && errServe != http.ErrServerClosed {
			logger.GetLogger().WithError(errServe).Error("HTTPServer")
		}
	}()

	logger.GetLogger().Warn("HTTP server started on ", *cfg.HTTPHost)

	s.closer.Append(s.stopServer(server))

	return nil
}
