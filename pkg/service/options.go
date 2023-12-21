package service

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"github.com/partyzanex/shortlink/pkg/ptr"
)

const (
	defaultDebugHost = ":8084"
	defaultGRPCHost  = ":8082"
	defaultHTTPHost  = ":8080"
	timeout          = time.Second * 5
)

type SetCustomRoute func(r gin.IRouter)
type SetGRPCHandler func(s *grpc.Server)
type SetGatewayHandler func(
	ctx context.Context,
	mux *runtime.ServeMux,
	endpoint string,
	opts []grpc.DialOption,
) error

type config struct {
	debugEnable  *bool
	debugHost    *string
	customRoutes SetCustomRoute

	GRPCEnable        *bool
	GRPCHost          *string
	GRPCReflection    *bool
	GRPCHandler       SetGRPCHandler
	GRPCServerOptions []grpc.ServerOption

	HTTPEnable       *bool
	HTTPHost         *string
	HTTPMiddlewares  []gin.HandlerFunc
	HTTPCustomRoutes SetCustomRoute

	gatewayEnable      *bool
	gatewayHandlers    []SetGatewayHandler
	gatewayMuxOptions  []runtime.ServeMuxOption
	gatewayDialOptions []grpc.DialOption
}

func newConfig(options []Option) *config {
	c := new(config)

	for _, option := range options {
		option(c)
	}

	return c
}

type Option func(c *config)

func WithDebug(enable bool) Option {
	return func(c *config) {
		c.debugEnable = ptr.Ptr(enable)
	}
}

func WithDebugHost(host string) Option {
	return func(c *config) {
		c.debugHost = ptr.Ptr(host)

		if c.debugEnable == nil {
			c.debugEnable = ptr.Ptr(true)
		}
	}
}

func WithCustomRoutes(fn SetCustomRoute) Option {
	return func(c *config) {
		c.customRoutes = fn

		if c.debugEnable == nil {
			c.debugEnable = ptr.Ptr(true)
		}
	}
}

func GRPCEnable(enable bool) Option {
	return func(c *config) {
		c.GRPCEnable = ptr.Ptr(enable)
	}
}

func WithGRPCHost(host string) Option {
	return func(c *config) {
		c.GRPCHost = ptr.Ptr(host)

		if c.GRPCEnable == nil {
			c.GRPCEnable = ptr.Ptr(true)
		}
	}
}

func GRPCReflection(reflection bool) Option {
	return func(c *config) {
		c.GRPCReflection = ptr.Ptr(reflection)
	}
}

func WithGRPCHandlers(fn SetGRPCHandler) Option {
	return func(c *config) {
		c.GRPCHandler = fn

		if c.GRPCEnable == nil {
			c.GRPCEnable = ptr.Ptr(true)
		}
	}
}

func WithGRPCOptions(options ...grpc.ServerOption) Option {
	return func(c *config) {
		c.GRPCServerOptions = options

		if c.GRPCEnable == nil {
			c.GRPCEnable = ptr.Ptr(true)
		}
	}
}

func HTTPEnable(enable bool) Option {
	return func(c *config) {
		c.HTTPEnable = ptr.Ptr(enable)
	}
}

func WithHTTPHost(host string) Option {
	return func(c *config) {
		c.HTTPHost = ptr.Ptr(host)

		if c.HTTPEnable == nil {
			c.HTTPEnable = ptr.Ptr(true)
		}
	}
}

func WithMiddlewares(middlewares ...gin.HandlerFunc) Option {
	return func(c *config) {
		c.HTTPMiddlewares = middlewares

		if c.HTTPEnable == nil {
			c.HTTPEnable = ptr.Ptr(true)
		}
	}
}

func WithHTTPRoutes(fn SetCustomRoute) Option {
	return func(c *config) {
		c.HTTPCustomRoutes = fn

		if c.HTTPEnable == nil {
			c.HTTPEnable = ptr.Ptr(true)
		}
	}
}

func GatewayEnable(enable bool) Option {
	return func(c *config) {
		c.gatewayEnable = ptr.Ptr(enable)
	}
}

func WithGatewayHandlers(fn ...SetGatewayHandler) Option {
	return func(c *config) {
		c.gatewayHandlers = fn

		if c.gatewayEnable == nil {
			c.gatewayEnable = ptr.Ptr(true)
		}
	}
}

func WithServerMuxOptions(options ...runtime.ServeMuxOption) Option {
	return func(c *config) {
		c.gatewayMuxOptions = options

		if c.gatewayEnable == nil {
			c.gatewayEnable = ptr.Ptr(true)
		}
	}
}

func WithDialOptions(options ...grpc.DialOption) Option {
	return func(c *config) {
		c.gatewayDialOptions = options

		if c.gatewayEnable == nil {
			c.gatewayEnable = ptr.Ptr(true)
		}
	}
}
