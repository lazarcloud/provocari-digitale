package jwt

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/lazarcloud/provocari-digitale/api/globals"
)

type JWTClaims struct {
	Type       string `json:"type"`
	Issuer     string `json:"issuer"`
	IssuedAt   int64  `json:"issuedAt"`
	ExpiresAt  int64  `json:"expiresAt"`
	UserId     string `json:"userId"`
	AccessRole string `json:"accessRole"`
}

func (c *JWTClaims) Valid() error {
	vErr := new(jwt.ValidationError)

	now := time.Now().Unix()

	if c.Issuer != globals.AuthIssuer {
		vErr.Inner = fmt.Errorf("invalid issuer")
		vErr.Errors |= jwt.ValidationErrorIssuer
	}
	if c.IssuedAt == 0 {
		vErr.Inner = fmt.Errorf("missing issuedAt")
		vErr.Errors |= jwt.ValidationErrorIssuedAt
	}
	if (c.IssuedAt - now) > 0 {
		vErr.Inner = fmt.Errorf("token used before issued")
		vErr.Errors |= jwt.ValidationErrorIssuedAt
	}
	if c.ExpiresAt == 0 {
		vErr.Inner = fmt.Errorf("missing expiresAt")
		vErr.Errors |= jwt.ValidationErrorExpired
	}
	if now > c.ExpiresAt {
		vErr.Inner = fmt.Errorf("token is expired")
		vErr.Errors |= jwt.ValidationErrorExpired
	}
	if !globals.AuthIsValidRole(c.AccessRole) {
		vErr.Inner = fmt.Errorf("invalid accessRole")
		vErr.Errors |= jwt.ValidationErrorClaimsInvalid
	}
	if !globals.AuthIsValidType(c.Type) {
		vErr.Inner = fmt.Errorf("invalid type")
		vErr.Errors |= jwt.ValidationErrorClaimsInvalid
	}

	if vErr.Errors > 0 {
		return vErr
	}

	return nil
}

func CreateJWTWithClaims(tokenType string, duration time.Duration, userId string, accessRole string) (string, error) {
	claims := &JWTClaims{
		Type:       tokenType,
		Issuer:     globals.AuthIssuer,
		IssuedAt:   time.Now().Unix(),
		ExpiresAt:  time.Now().Add(duration).Unix(),
		UserId:     userId,
		AccessRole: accessRole,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(globals.AuthSecretJWTKey)

	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func ParseJWT(tokenStr string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return globals.AuthSecretJWTKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
