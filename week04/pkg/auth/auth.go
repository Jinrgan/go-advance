package auth

import (
	"context"
	"fmt"
	userpb "go-advance/week04/app/user/api/gen/v1"
	"go-advance/week04/pkg/auth/token"
	"go-advance/week04/pkg/id"
	"io/ioutil"
	"net/http"
	"os"

	"go.uber.org/zap"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const (
	uidField = "user-id"
)

type tokenVerifier interface {
	Verify(token string) (string, error)
}

type middleware struct {
	verifier tokenVerifier
	userClt  userpb.UserServiceClient
}

//Middleware creates a gin auth middleware
func NewMiddleware(publicKeyfile string, clt userpb.UserServiceClient) (*middleware, error) {
	f, err := os.Open(publicKeyfile)
	if err != nil {
		return nil, fmt.Errorf("cannot open public key file: %v", err)
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("cannot read public key file: %v", err)
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(b)
	if err != nil {
		return nil, fmt.Errorf("cannot parse public key: %v", err)
	}

	return &middleware{
		verifier: &token.JWTTokenVerifier{
			PublicKey: pubKey,
		},
		userClt: clt,
	}, nil
}

func (m *middleware) HandleReq(ctx *gin.Context) {
	tkn := ctx.Request.Header.Get("Token")
	if tkn == "" {
		ctx.JSON(http.StatusUnauthorized, map[string]string{"msg": "please login in"})
		ctx.Abort()
		return
	}

	uid, err := m.verifier.Verify(tkn)
	if err != nil {
		ctx.String(http.StatusUnauthorized, "token not valid: %v", err)
		ctx.Abort()
		return
	}

	ctx.Set(uidField, uid)
	ctx.Next()
}

func (m *middleware) HandleAdminReq(ctx *gin.Context) {
	v, ok := ctx.Get(uidField)
	if !ok {
		ctx.Status(http.StatusInternalServerError)
		ctx.Abort()
		return
	}

	uid, ok := v.(string)
	if !ok {
		zap.L().Error("uidField value with unknown type")
		ctx.Status(http.StatusInternalServerError)
		ctx.Abort()
		return
	}

	user, err := m.userClt.GetUser(context.Background(), &userpb.GetUserRequest{Id: uid})
	if err != nil {
		zap.L().Error("cannot get user", zap.Error(err))
		ctx.Status(http.StatusInternalServerError)
		ctx.Abort()
		return
	}

	if user.User.Role != userpb.Roles_Admin {
		ctx.JSON(http.StatusForbidden, gin.H{"msg": "no permission"})
		ctx.Abort()
		return
	}

	ctx.Next()
}

//UserIDFromContext gets account id from context
// returns abort handler error if no account id is available.
func UserIDFromContext(ctx *gin.Context) (id.UserID, error) {
	v, ok := ctx.Get(uidField)
	if !ok {
		return "", http.ErrAbortHandler
	}

	uid, ok := v.(string)
	if !ok {
		return "", http.ErrAbortHandler
	}

	return id.UserID(uid), nil
}
