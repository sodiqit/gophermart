package auth

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryAuthInterceptor(tokenService TokenService, protectedMethods []string) grpc.UnaryServerInterceptor {
	protectedMethodsMap := make(map[string]bool)
	for _, method := range protectedMethods {
		protectedMethodsMap[method] = true
	}

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {

		_, ok := protectedMethodsMap[info.FullMethod]
		if !ok {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "Metadata not provided")
		}

		values := md.Get("token")
		if len(values) == 0 {
			return nil, status.Error(codes.Unauthenticated, "No token provided")
		}

		token := values[0]

		claims, err := tokenService.Validate(token)

		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "Invalid token")
		}

		ctx = context.WithValue(ctx, ClaimsContextKey, claims)

		return handler(ctx, req)
	}
}
