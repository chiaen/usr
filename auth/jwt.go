package auth

import (
	"errors"
	"fmt"
	"log"
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

	tokenSubject  = "usr"
	tokenIssuer   = "auth.usr"
	tokenAudience = "api.usr"
)

var (
	tokenSigningMethod             = crypto.SigningMethodHS256
	accessTokenKey     interface{} = []byte("m6VdMH+DnoSAK/0brNDG/N1JYAFJUUI4/Q8q60BU9fc=") // a random key

	ErrTokenInvalid = errors.New("token invalid")
	ErrTokenExpired = errors.New("token expired")
)

type AccessToken interface {
	fmt.Stringer

	ExpiresIn() int32
}

type accessTokenImpl struct {
	jwt.JWT
	signKey interface{}
	cache   string
}

func (t *accessTokenImpl) parse(tokenStr string) error {
	parsed, err := jws.ParseJWT([]byte(tokenStr))
	if err != nil {
		return err
	}
	if err = parsed.Validate(t.signKey, tokenSigningMethod, tokenValidator()); err == jwt.ErrTokenIsExpired {
		return ErrTokenExpired
	} else if err != nil {
		return ErrTokenInvalid
	}
	t.JWT = parsed
	t.cache = tokenStr
	return nil
}

func (t *accessTokenImpl) String() string {
	if t.cache != "" {
		return t.cache
	} else if jwtBytes, err := t.Serialize(t.signKey); err != nil {
		log.Printf("access token serialization failed", "error", err, "claims", t.Claims())
		return ""
	} else {
		t.cache = string(jwtBytes)
	}
	return t.cache
}

func (t *accessTokenImpl) ExpiresIn() int32 {
	exp, _ := t.Claims().Expiration()
	nbf, _ := t.Claims().NotBefore()
	return int32(exp.Unix() - nbf.Unix())
}

func IssueToken(uid string) AccessToken {
	j := baseJWT(accessTokenExpiration)
	j.Claims().Set("uid", uid)
	return &accessTokenImpl{
		JWT:     j,
		signKey: accessTokenKey,
	}
}

// ParseAccessToken parses a token string and transforms into an access token.
func ParseAccessToken(tokenStr string) (AccessToken, error) {
	token := &accessTokenImpl{
		JWT:     baseJWT(accessTokenExpiration),
		signKey: accessTokenKey,
	}
	if err := token.parse(tokenStr); err != nil {
		return nil, err
	}
	return token, nil
}

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
