package service

import (
	"go-advance/week04/app/user/api"
	userpb "go-advance/week04/app/user/api/gen/v1"
	"go-advance/week04/app/user/internal/biz"
	"go-advance/week04/pkg/auth"
	"go-advance/week04/pkg/id"
	pkgmysql "go-advance/week04/pkg/mysql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
)

type Interface struct {
	CaptchaVerifier CaptchaVerifier
	TokenGenerator  TokenGenerator
	TokenExpire     time.Duration
	CaptchaStore    base64Captcha.Store
	Biz             *biz.UserBiz
	Repo
}

func (s *Interface) Register(ctx *gin.Context) {
	var form api.RegisterForm
	err := ctx.ShouldBind(&form)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	err = s.CaptchaVerifier.Verify(form.Mobile, form.Code) // TODO: move to middleware
	if err != nil {
		ctx.String(http.StatusBadRequest, "wrong captcha, please try again")
		return
	}

	uid, err := s.Biz.CreateUser(&biz.User{
		Mobile:   form.Mobile,
		Password: form.Password,
	})
	if err != nil {
		zap.L().Error("cannot create user", zap.Error(err))
		ctx.String(http.StatusInternalServerError, "cannot create user")
		return
	}

	ctx.String(http.StatusOK, "%s", uid)
}

func (s *Interface) Login(ctx *gin.Context) {
	var form api.LoginForm
	err := ctx.ShouldBind(&form)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	err = s.Biz.Verify(id.UserID(form.ID), form.Password)
	if err != nil {
		ctx.String(http.StatusUnauthorized, "wrong password")
		return
	}

	return
}

func (s *Interface) LoginByMobile(ctx *gin.Context) {
	var form api.LoginByMobileForm
	err := ctx.ShouldBind(&form)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	if !s.CaptchaStore.Verify(form.CaptchaID, form.Captcha, true) {
		ctx.String(http.StatusBadRequest, "wrong captcha")
		return
	}

	uid, err := s.GetUserIDByMobile(form.Mobile)
	if err != nil {
		zap.L().Error("cannot get user id by mobile", zap.String("mobile", form.Mobile), zap.Error(err))
		ctx.Status(http.StatusInternalServerError)
		return
	}
	err = s.Biz.Verify(pkgmysql.ObjectIDToUserID(uid), form.Password)
	if err != nil {
		ctx.String(http.StatusUnauthorized, "wrong password")
		return
	}

	tkn, err := s.TokenGenerator.GenerateToken(uid.String(), s.TokenExpire)
	if err != nil {
		zap.L().Error("cannot generate token", zap.Error(err))
		ctx.Status(http.StatusInternalServerError)
		return
	}

	user, err := s.Biz.GetUser(pkgmysql.ObjectIDToUserID(uid))
	if err != nil {
		zap.L().Error("cannot get user", zap.Error(err))
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, &api.UserLoginResponse{
		NickName:  user.NickName,
		Token:     tkn,
		ExpiredIn: s.TokenExpire,
	})

	return
}

func (s *Interface) GetUsers(ctx *gin.Context) {
	uid, err := auth.UserIDFromContext(ctx)
	if err != nil {
		ctx.Status(http.StatusUnauthorized)
		return
	}
	zap.S().Infof("user login in: %v", uid)

	pn, err := strconv.Atoi(ctx.DefaultQuery("page_number", "0"))
	if err != nil {
		ctx.String(http.StatusBadRequest, "wrong page number")
		return
	}

	ps, err := strconv.Atoi(ctx.DefaultQuery("page_size", "0"))
	if err != nil {
		ctx.String(http.StatusBadRequest, "wrong page size")
		return
	}

	users, err := s.Repo.GetUsers(pn, ps)
	if err != nil {
		zap.L().Error("cannot get users", zap.Error(err))
		ctx.Status(http.StatusInternalServerError)
		return
	}

	var res []*userpb.UserEntity
	for _, user := range users {
		res = append(res, dataToPb(user))
	}

	ctx.JSON(http.StatusOK, res)

	return
}
