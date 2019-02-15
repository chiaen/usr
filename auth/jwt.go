package main

import (
	"errors"
	"time"

	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
	"github.com/SermoDigital/jose/jwt"
)

const (
	AccessTokenTypeBearer = "Bearer"
)

// TODO: make token claims configurable
const (
	maxExpiration         = 86400
	accessTokenExpiresIn  = 3600
	accessTokenExpiration = accessTokenExpiresIn * time.Second
	idTokenExpiration     = 120 * time.Hour
	tokenLeeway           = time.Minute

	tokenSubject  = "erp"
	tokenIssuer   = "auth.erp.lativ.com"
	tokenAudience = "api.erp.lativ.com"
)

var (
	// jwt configurations
	tokenSigningMethod             = crypto.SigningMethodHS256
	accessTokenKey     interface{} = []byte("m6VdMH+DnoSAK/0brNDG/N1JYAFJUUI4/Q8q60BU9fc=") // a random key
)

var timeNow = func() time.Time {
	return time.Now().UTC()
}

func baseJWT(expiration time.Duration) jwt.JWT {
	claims := jws.Claims{}
	claims.SetSubject(tokenSubject)
	claims.SetIssuer(tokenIssuer)
	claims.SetAudience(tokenAudience)
	now := timeNow()
	claims.SetIssuedAt(now)
	claims.SetNotBefore(now)
	claims.SetExpiration(now.Add(expiration))
	token := jws.NewJWT(claims, tokenSigningMethod)
	return token
}

func tokenValidator() *jwt.Validator {
	claims := jws.Claims{}
	claims.SetSubject(tokenSubject)
	claims.SetIssuer(tokenIssuer)
	claims.SetAudience(tokenAudience)
	return jws.NewValidator(claims, tokenLeeway, tokenLeeway, mustContainUID)
}

func mustContainUID(claims jws.Claims) error {
	uid, ok := claims.Get("uid").(string)
	if !ok {
		return errors.New("missing uid")
	} else if uid == "" {
		return errors.New("invalid uid")
	}
	return nil
}
