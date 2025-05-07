package service

import (
	"encoding/base64"
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtUtil struct {
	secretKey []byte
	TokenTTL  int64
}

var ErrJwtBroken = errors.New("JWT cannot be handled")
var ErrJwtSigning = errors.New("GenerateJWT: Couldn't sign JWT")
var ErrJwtMalformed = errors.New("ParseJWT: JWT is malformed")
var ErrJwtInvalidSignature = errors.New("ParseJWT: JWT signature is invalid")
var ErrJwtNotAavailableNow = errors.New("ParseJWT: JWT is not active at this time")
var ErrJwtClaimsErr = errors.New("ParseJWT: Couldn't extract JWT claims")

func NewJwtUtil() *JwtUtil {
	base64Key := os.Getenv("JWT_SECRET_KEY")
	byteKey, err := base64.StdEncoding.DecodeString(base64Key)

	if err != nil {
		log.Fatal("Missing env variable JWT_SECRET_KEY, or the env variable isn't base64 encoded.")
	}

	var tokenTTL int64 = 24 * 60 * 60
	return &JwtUtil{
		secretKey: byteKey,
		TokenTTL:  tokenTTL,
	}
}

func (ju *JwtUtil) GenerateJWT(subject string) (string, error) {
	currentTime := time.Now().Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": subject,
			"iat": currentTime,
			"exp": currentTime + ju.TokenTTL,
		})

	signedToken, err := token.SignedString(ju.secretKey)

	if err != nil {
		return "", ErrJwtSigning
	}

	return signedToken, nil
}

func (ju *JwtUtil) ParseJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return ju.secretKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenMalformed):
			return "", ErrJwtMalformed
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			return "", ErrJwtInvalidSignature
		case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
			return "", ErrJwtNotAavailableNow
		default:
			return "", ErrJwtBroken
		}
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userEmail, ok := claims["sub"].(string)
		if !ok {
			return "", ErrJwtClaimsErr
		}

		return userEmail, nil
	} else {
		return "", ErrJwtClaimsErr
	}
}
