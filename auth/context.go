package auth

import (
	"context"
	"github.com/SermoDigital/jose/jws"
)

func WithAccessTokenContext(ctx context.Context, token AccessToken) context.Context {
	return context.WithValue(ctx, ctxAccessTokenKey{}, token)
}

func UserIDfromValidAccessToken(ctx context.Context) (uid string, ok bool) {
	token, ok := ctx.Value(ctxAccessTokenKey{}).(AccessToken)
	parsed, err := jws.ParseJWT([]byte(token.String()))
	if err != nil {
		return "", false
	}
	if err := parsed.Validate(accessTokenKey, tokenSigningMethod, tokenValidator()); err != nil {
		return "", false
	}
	return parsed.Claims().Get("uid").(string), true
}

type ctxAccessTokenKey struct{}