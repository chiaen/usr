package auth

import (
	"context"
	"strings"

	middle "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	HeaderAuthorization = "authorization"
)

// authFromRequestMetadata tries to obtain an access token string from an incoming request metadata.
func authFromRequestMetadata(ctx context.Context) string {
	val := metautils.ExtractIncoming(ctx).Get(HeaderAuthorization)
	if val == "" {
		return ""
	}
	splits := strings.SplitN(val, " ", 2)
	if len(splits) < 2 {
		return ""
	}
	if !strings.EqualFold(splits[0], "BEARER") {
		return ""
	}
	return splits[1]
}

// ensureAccessToken ensures an access token is bring in incoming request metadata,
// otherwise returns an unauthenticated status error.
func ensureAccessToken(ctx context.Context) (AccessToken, error) {
	tokenStr := authFromRequestMetadata(ctx)
	if tokenStr == "" {
		return nil, status.Error(codes.Unauthenticated, "missing token")
	}
	token, err := ParseAccessToken(tokenStr)
	if err == ErrTokenInvalid {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	} else if err == ErrTokenExpired {
		return nil, status.Error(codes.Unauthenticated, "token expired")
	} else if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}
	return token, nil
}

var (
	tokenAuthFunc middle.AuthFunc = func(ctx context.Context) (context.Context, error) {
		t, err := ensureAccessToken(ctx)
		if err != nil {
			return nil, err
		}
		return WithAccessTokenContext(ctx, t), nil
	}
	UnaryTokenVerifier = middle.UnaryServerInterceptor(tokenAuthFunc)
)
