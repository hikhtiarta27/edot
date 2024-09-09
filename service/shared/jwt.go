package shared

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	tokenLookup = "header:" + echo.HeaderAuthorization
	claims      = jwt.StandardClaims{}
)

type Jwt struct {
	secret string
}

func New(secret string) *Jwt {
	return &Jwt{
		secret: secret,
	}
}

func (j *Jwt) GenerateToken(duration time.Duration, userId string) (string, error) {

	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   userId,
		ExpiresAt: now.Add(duration).Unix(),
		IssuedAt:  now.Unix(),
		NotBefore: now.Unix(),
		Issuer:    "edot",
	})

	generatedToken, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", err
	}

	return generatedToken, nil
}

func (j Jwt) tokenLookupFunc(auth string, c echo.Context) (interface{}, error) {
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != "HS256" {
			return nil, errors.New("invalid signing method")
		}
		return []byte(j.secret), nil
	}

	token, err := jwt.Parse(auth, keyFunc)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}

func (j Jwt) Validate() echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:         claims,
		TokenLookup:    tokenLookup,
		ParseTokenFunc: j.tokenLookupFunc,
		ErrorHandlerWithContext: func(err error, c echo.Context) error {
			if errors.Is(err, middleware.ErrJWTMissing) {
				return FailResponse(c, "missing or malformed jwt", nil)
			}

			return SuccessResponse(c, "invalid or expired jwt", nil)
		},
	})
}
