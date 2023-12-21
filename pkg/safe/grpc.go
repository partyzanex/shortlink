package safe

import (
	"context"

	"google.golang.org/grpc"
)

// GRPCRecover writes recovering result to GRPC error.
func GRPCRecover() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		defer RecoverToError(&err)

		resp, err = handler(ctx, req)

		return resp, err
	}
}
