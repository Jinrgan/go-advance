package token

import (
	"crypto/rsa"
	"time"

	"github.com/dgrijalva/jwt-go"
)

//JWTTokenGen generates a JWT token
type JWTTokenGen struct {
	privateKey *rsa.PrivateKey
	issuer     string
	nowFunc    func() time.Time
}

//GenerateToken generates a token
func (g *JWTTokenGen) GenerateToken(accountID string, expire time.Duration) (string, error) {
	nowSec := g.nowFunc().Unix()
	tkn := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.StandardClaims{
		ExpiresAt: nowSec + int64(expire.Seconds()),
		IssuedAt:  nowSec,
		Issuer:    g.issuer,
		Subject:   accountID,
	})

	return tkn.SignedString(g.privateKey)
}

//NewJWTTokenGen creates a JWTTokenGen
func NewJWTTokenGen(issuer string, privateKey *rsa.PrivateKey) *JWTTokenGen {
	return &JWTTokenGen{
		privateKey: privateKey,
		issuer:     issuer,
		nowFunc:    time.Now,
	}
}
