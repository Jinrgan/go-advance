package data

import (
	"errors"
	userpb "go-advance/week04/app/user/api/gen/v1"
	"go-advance/week04/pkg/id"
	pkgmysql "go-advance/week04/pkg/mysql"
	"go-advance/week04/pkg/mysql/model"
	"time"

	"github.com/go-sql-driver/mysql"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type User struct {
	model.Base
	Mobile   string    `gorm:"index:idx_mobile;unique;type:varchar(11);not null"`
	Password string    `gorm:"type:varchar(100);not null"`
	NickName string    `gorm:"type:varchar(20)"`
	Birthday time.Time `gorm:"type:datetime"`
	Gender   string    `gorm:"column:gender;default:male;type:varchar(6) comment 'female 表示女，male 表示男'"`
	Role     int32     `gorm:"column:role;default:1;type:int comment '1 表示普通用户，2 表示管理员'"`
}

func ToResp(d *User) *userpb.UserEntity {
	return &userpb.UserEntity{
		Id: pkgmysql.ObjectIDToUserID(d.ID).String(),
		User: &userpb.User{
			Mobile:   d.Mobile,
			NickName: d.NickName,
			Birthday: d.Birthday.Format("2006-01-02"),
			//Birthday: u.Birthday.Unix(),
			Gender: d.Gender,
			Role:   userpb.Roles(d.Role),
		},
	}
}

const mysqlDuplicateErr = 1062

type MySQL struct {
	DB *gorm.DB
}

func (m *MySQL) CreateUser(user *User) (id.UserID, error) {
	res := m.DB.Create(user)
	if err, ok := res.Error.(*mysql.MySQLError); ok {
		if err.Number == mysqlDuplicateErr {
			return "", status.Error(codes.AlreadyExists, "用户已存在")
		}
	} else if res.Error != nil {
		zap.L().Error("cannot get create user", zap.Error(res.Error))
		return "", status.Error(codes.Internal, "")
	}

	return pkgmysql.ObjectIDToUserID(user.ID), nil
}

func (m *MySQL) GetUser(uid pkgmysql.ObjectID) (*User, error) {
	var user User
	res := m.DB.First(&user, uid)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "用户不存在")
	} else if res.Error != nil {
		zap.L().Error("cannot get user by id", zap.Error(res.Error))
		return nil, status.Error(codes.Internal, "")
	}

	return &user, nil
}

func (m *MySQL) GetUserByMobile(mobile string) (*User, error) {
	user := &User{
		Mobile: mobile,
	}
	res := m.DB.Find(&user)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "用户不存在")
	} else if res.Error != nil {
		zap.L().Error("cannot get user by id", zap.Error(res.Error))
		return nil, status.Error(codes.Internal, "")
	}

	return user, nil
}

func (m *MySQL) GetUserIDByMobile(mobile string) (pkgmysql.ObjectID, error) {
	var user User
	res := m.DB.Where(&User{Mobile: mobile}).First(&user) // TODO: select id from
	if res.RowsAffected == 0 {
		return 0, status.Error(codes.NotFound, "用户不存在")
	}
	if res.Error != nil {
		zap.L().Error("cannot get user by mobile", zap.Error(res.Error))
		return 0, status.Error(codes.Internal, "")
	}

	return user.ID, nil
}

func (m *MySQL) GetUserEncPwd(uid id.UserID) (string, error) {
	var user User
	res := m.DB.First(&user, uid)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return "", status.Error(codes.NotFound, "用户不存在")
	} else if res.Error != nil {
		zap.L().Error("cannot get user by id", zap.Error(res.Error))
		return "", status.Error(codes.Internal, "")
	}

	return user.Password, nil
}

func (m *MySQL) GetUsers(page, size int) ([]*User, error) {
	var users []*User
	// has users?
	res := m.DB.Find(&users)
	if res.Error != nil {
		zap.L().Error("cannot get users", zap.Error(res.Error))
		return nil, status.Error(codes.NotFound, "")
	}

	var resp []*User
	res = m.DB.Scopes(pkgmysql.Paginate(page, size)).Find(&users)
	if res.Error != nil {
		zap.L().Error("cannot paginate users", zap.Error(res.Error))
		return nil, status.Error(codes.Internal, "")
	}

	for _, user := range users {
		resp = append(resp, user)
	}

	return resp, nil
}

func (m *MySQL) UpdateUser(uid pkgmysql.ObjectID, user *User) (*User, error) {
	res := m.DB.Model(&User{
		Base: model.Base{ID: uid},
	}).Updates(user)
	if res.Error != nil {
		zap.L().Error("cannot update user", zap.Error(res.Error))
		return nil, status.Error(codes.Internal, "")
	}

	return user, nil
}
