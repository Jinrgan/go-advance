package biz

import (
	"fmt"
	"go-advance/week04/app/user/internal/data"
	"go-advance/week04/pkg/id"
	pkgmysql "go-advance/week04/pkg/mysql"
	"time"
)

type User struct {
	Mobile   string
	NickName string
	Password string
	Birthday time.Time
	Gender   string
	Role     int32
}

type Coder interface {
	Gen(code string) string
	Verify(code, encCode string) error
}

type UserBiz struct {
	Coder Coder // TODO: change to code verifier
	Repo  Repo
}

type Repo interface {
	CreateUser(user *data.User) (id.UserID, error)
	GetUserEncPwd(uid id.UserID) (string, error)
	GetUser(uid pkgmysql.ObjectID) (*data.User, error)
}

func (b *UserBiz) CreateUser(ob *User) (id.UserID, error) {
	encPwd := b.Coder.Gen(ob.Password)
	uid, err := b.Repo.CreateUser(&data.User{
		Mobile:   ob.Mobile,
		Password: encPwd,
		NickName: ob.NickName,
		Birthday: ob.Birthday,
		Gender:   ob.Gender,
		Role:     ob.Role,
	})
	if err != nil {
		return "", err
	}

	return uid, nil
}

func (b *UserBiz) GetUser(uid id.UserID) (*User, error) {
	objID, err := pkgmysql.ObjectIDFromID(uid)
	if err != nil {
		return nil, fmt.Errorf("cannot get object id: %v", err)
	}
	user, err := b.Repo.GetUser(objID)
	if err != nil {
		return nil, err
	}

	return &User{
		Mobile:   user.Mobile,
		NickName: user.NickName,
		Password: user.Password,
		Birthday: user.Birthday,
		Gender:   user.Gender,
		Role:     user.Role,
	}, nil
}

func (b *UserBiz) Verify(uid id.UserID, pwd string) error {
	code, err := b.Repo.GetUserEncPwd(uid)
	if err != nil {
		return fmt.Errorf("cannot get user encrypted password: %v", err)
	}

	err = b.Coder.Verify(pwd, code)
	if err != nil {
		return fmt.Errorf("cannot verify password: %v", err)
	}

	return nil
}
