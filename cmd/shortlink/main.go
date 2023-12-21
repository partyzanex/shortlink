package main

import (
	"database/sql"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/partyzanex/shortlink/internal/config"
	"github.com/partyzanex/shortlink/internal/link"
	transport_grpc "github.com/partyzanex/shortlink/internal/transport/grpc"
	transport_http "github.com/partyzanex/shortlink/internal/transport/http"
	"github.com/partyzanex/shortlink/pkg/logger"
	pb "github.com/partyzanex/shortlink/pkg/proto/go/shorts"
	"github.com/partyzanex/shortlink/pkg/safe"
	"github.com/partyzanex/shortlink/pkg/service"
	"github.com/partyzanex/shutdown"
)

func main() {
	app := cli.NewApp()
	app.Name = config.AppName
	app.Description = config.AppDesc
	app.Flags = []cli.Flag{
		config.EnvFlag(),
		config.DebugHostFlag(),
		config.GrpcHostFlag(),
		config.GrpcReflectionFlag(),
		config.HashLengthFlag(),
		config.HttpHostFlag(),
		config.LogLevelFlag(),
		config.PostgresReadTimeoutFlag(),
		config.PostgresUrlFlag(),
		config.PostgresWriteTimeoutFlag(),
	}

	app.Action = action

	if err := app.Run(os.Args); err != nil {
		logger.GetLogger().WithError(err).Fatal("app.Run")
	}
}

func action(ctx *cli.Context) error {
	logger.ParseLogLevel(config.LogLevel.String())

	srv := service.New()
	shutdown.Append(srv)

	err := srv.Init(service.WithDebugHost(config.DebugHost.String()))
	if err != nil {
		return errors.Wrap(err, "srv.Init")
	}

	db, err := sql.Open("postgres", config.PostgresUrl.String())
	if err != nil {
		return errors.Wrap(err, "sql.Open")
	}

	shutdown.Append(db)

	links := link.NewService(
		db,
		config.PostgresReadTimeout.Duration(),
		config.PostgresWriteTimeout.Duration(),
		config.HashLength.Int(),
	)

	grpcService := transport_grpc.NewService(links)

	err = srv.Start(
		service.WithHTTPHost(config.HttpHost.String()),
		service.WithMiddlewares(
			gin.Recovery(),
			logger.Middleware(logger.GetLogger()),
		),
		service.WithHTTPRoutes(func(r gin.IRouter) {
			r.StaticFS("/favicon.ico", http.FS(transport_http.Favicon))
			r.GET("/:id", transport_http.GetRedirectHandler(links))
		}),
		service.WithGatewayHandlers(
			pb.RegisterShortsHandlerFromEndpoint,
		),
		service.WithDialOptions(
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
		service.WithGRPCHost(config.GrpcHost.String()),
		service.GRPCReflection(config.GrpcReflection.Bool()),
		service.WithGRPCHandlers(func(s *grpc.Server) {
			pb.RegisterShortsServer(s, grpcService)
		}),
		service.WithGRPCOptions(
			grpc.ChainUnaryInterceptor(
				safe.GRPCRecover(),
			),
		),
	)
	if err != nil {
		return errors.Wrap(err, "srv.Start")
	}

	return shutdown.CloseOnSignalContext(ctx.Context, logger.GetLogger(), os.Kill, os.Interrupt)
}
