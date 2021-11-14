package service

import (
	"context"
	userpb "go-advance/week04/app/user/api/gen/v1"
	"go-advance/week04/app/user/internal/biz"
	"go-advance/week04/app/user/internal/data"
	"go-advance/week04/pkg/id"
	pkgmysql "go-advance/week04/pkg/mysql"
	"strconv"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go.uber.org/zap"

	"github.com/golang/protobuf/ptypes/empty"
)

type userError string

func (e userError) Error() string {
	return e.Message()
}

func (e userError) Message() string {
	return string(e)
}

type Coder interface {
	Gen(code string) string
	Verify(code, encCode string) error
}

type CaptchaVerifier interface {
	Verify(mobile, code string) error
}

type TokenGenerator interface {
	GenerateToken(account string, expire time.Duration) (string, error)
}

type Repo interface {
	GetUser(uid pkgmysql.ObjectID) (*data.User, error)
	GetUserByMobile(mobile string) (*data.User, error)
	GetUserIDByMobile(mobile string) (pkgmysql.ObjectID, error)
	GetUserEncPwd(uid id.UserID) (string, error)
	GetUsers(page, pageSize int) ([]*data.User, error)
	UpdateUser(uid pkgmysql.ObjectID, user *data.User) (*data.User, error)
}

type UserService struct {
	Biz *biz.UserBiz
	Repo
}

func (s *UserService) CreateUser(_ context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	uid, err := s.Biz.CreateUser(&biz.User{
		Mobile:   req.Mobile,
		NickName: req.NickName,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	return &userpb.CreateUserResponse{Id: uid.String()}, err
}

func (s *UserService) GetUser(_ context.Context, req *userpb.GetUserRequest) (*userpb.UserEntity, error) {
	uid, err := strconv.Atoi(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "")
	}
	user, err := s.Repo.GetUser(pkgmysql.ObjectID(uid))
	if err != nil {
		return nil, err
	}

	return dataToPb(user), nil
}

func (s *UserService) GetUserByMobile(_ context.Context, req *userpb.GetUserByMobileRequest) (*userpb.UserEntity, error) {
	user, err := s.Repo.GetUserByMobile(req.Mobile)
	if err != nil {
		return nil, err
	}

	return dataToPb(user), nil
}

func (s *UserService) UpdateUser(_ context.Context, req *userpb.UpdateUserRequest) (*empty.Empty, error) {
	uid, err := strconv.Atoi(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "")
	}
	_, err = s.Repo.UpdateUser(pkgmysql.ObjectID(uid), &data.User{
		NickName: req.NickName,
		Birthday: time.Unix(req.Birthday, 0),
		Gender:   req.Gender,
	})
	if err != nil {
		zap.L().Error("cannot update user", zap.Error(err))
		return nil, status.Error(codes.Internal, "")
	}

	return nil, nil
}

func dataToPb(user *data.User) *userpb.UserEntity {
	return &userpb.UserEntity{
		Id: user.ID.String(),
		User: &userpb.User{
			Mobile:   user.Mobile,
			NickName: user.NickName,
			Birthday: user.Birthday.Format("2006-01-01"),
			Gender:   user.Gender,
			Role:     userpb.Roles(user.Role),
		},
	}
}
