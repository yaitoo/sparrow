package auth

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Authorize is the pluggable function that performs authentication.
type Authorize func(ctx context.Context) (context.Context, error)

// Authorizer allows a given gRPC service implementation to override the global `Authorize`.
//
// If a service implements the AuthorizeOverride method, it takes precedence over the `Authorize` method,
// and will be called instead of Authorize for all method invocations within that service.
type Authorizer interface {
	Authorize(ctx context.Context, fullMethodName string) (context.Context, error)
}

// UnaryServerInterceptor returns a new unary server interceptors that performs per-request auth.
func UnaryServerInterceptor(authorize Authorize) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var newCtx context.Context
		var err error
		if overrideSrv, ok := info.Server.(Authorizer); ok {
			newCtx, err = overrideSrv.Authorize(ctx, info.FullMethod)
		} else {
			newCtx, err = authorize(ctx)
		}
		if err != nil {
			return nil, err
		}
		return handler(newCtx, req)
	}
}

// StreamServerInterceptor returns a new unary server interceptors that performs per-request auth.
func StreamServerInterceptor(authorize Authorize) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		var newCtx context.Context
		var err error
		if overrideSrv, ok := srv.(Authorizer); ok {
			newCtx, err = overrideSrv.Authorize(stream.Context(), info.FullMethod)
		} else {
			newCtx, err = authorize(stream.Context())
		}
		if err != nil {
			return err
		}
		wrapped := WrapServerStream(stream)
		wrapped.WrappedContext = newCtx
		return handler(srv, wrapped)
	}
}
