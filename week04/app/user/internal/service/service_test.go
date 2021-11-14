package service

//import (
//	"context"
//	"crypto/sha512"
//	"fmt"
//	userpb "go-advance/week04/app/user/api/gen/v1"
//	"go-advance/week04/app/user/internal/data"
//	"go-advance/week04/pkg/mysql"
//	mysqltesting "go-advance/week04/pkg/mysql/testing"
//	"math/rand"
//	"os"
//	"testing"
//
//	"github.com/anaskhan96/go-password-encoder"
//	"github.com/google/go-cmp/cmp"
//	"go.uber.org/zap"
//	"google.golang.org/grpc/codes"
//	"google.golang.org/grpc/status"
//	"google.golang.org/protobuf/testing/protocmp"
//)
//
//func TestService_CreateUser(t *testing.T) {
//	s := newService(t)
//
//	cases := []struct {
//		name        string
//		nickName    string
//		passwd      string
//		mobile      string
//		wantErrCode codes.Code
//	}{
//		{
//			name:     "user",
//			nickName: "user",
//			passwd:   "admin",
//			mobile:   "18866668888",
//		},
//		{
//			name:        "existing_user",
//			nickName:    "existing_user",
//			passwd:      "admin",
//			mobile:      "18866668888",
//			wantErrCode: codes.AlreadyExists,
//		},
//		{
//			name:     "another_user",
//			nickName: "another_user",
//			passwd:   "admin",
//			mobile:   "18866668889",
//		},
//	}
//
//	ctx := context.Background()
//	for _, cc := range cases {
//		t.Run(cc.name, func(t *testing.T) {
//			var code codes.Code
//			_, err := s.CreateUser(ctx, &userpb.CreateUserRequest{
//				NickName: cc.nickName,
//				Password: cc.passwd,
//				Mobile:   cc.mobile,
//			})
//			if err != nil {
//				if s, ok := status.FromError(err); ok {
//					code = s.Code()
//				} else {
//					t.Errorf("error creating product: %v", err)
//				}
//			}
//			if code != cc.wantErrCode {
//				t.Errorf("wrong err code: want %d, got %d", cc.wantErrCode, code)
//			}
//		})
//	}
//}
//
//func TestService_GetUser(t *testing.T) {
//	s := newService(t)
//
//	users := []*data.User{
//		{
//			NickName: "user_1",
//			Mobile:   "mobile_1",
//		},
//		{
//			NickName: "user_2",
//			Mobile:   "mobile_2",
//		},
//	}
//
//	res := s.DB.Create(&users)
//	if res.Error != nil {
//		t.Fatalf("cannot create users: %v", res.Error)
//	}
//
//	cases := []struct {
//		name        string
//		id          string
//		want        *userpb.UserEntity
//		wantErrCode codes.Code
//	}{
//		{
//			name: "exist_user",
//			id:   mysql.ObjectIDToUserID(users[0].ID).String(),
//			want: data.ToResp(users[0]),
//		},
//		{
//			name: "another_exist_user",
//			id:   mysql.ObjectIDToUserID(users[1].ID).String(),
//			want: data.ToResp(users[1]),
//		},
//		{
//			name:        "not_exist_user",
//			id:          "999",
//			wantErrCode: codes.NotFound,
//		},
//	}
//
//	for _, cc := range cases {
//		t.Run(cc.name, func(t *testing.T) {
//			code := codes.OK
//			got, err := s.GetUser(context.Background(), &userpb.GetUserRequest{Id: cc.id})
//			if err != nil {
//				if s, ok := status.FromError(err); ok {
//					code = s.Code()
//				} else {
//					t.Errorf("operation failed: %v", err)
//				}
//			}
//			if code != cc.wantErrCode {
//				t.Errorf("wrong err code: want %d, got %d", cc.wantErrCode, code)
//			}
//
//			if diff := cmp.Diff(cc.want, got, protocmp.Transform()); diff != "" {
//				t.Errorf("result differs; -want +got: %s", diff)
//			}
//		})
//	}
//}
//
//func TestService_GetUserByMobile(t *testing.T) {
//	s := newService(t)
//
//	users := []*data.User{
//		{
//			NickName: "user_1",
//			Mobile:   "mobile_1",
//		},
//		{
//			NickName: "user_2",
//			Mobile:   "mobile_2",
//		},
//	}
//
//	res := s.DB.Create(&users)
//	if res.Error != nil {
//		t.Fatalf("cannot create users: %v", res.Error)
//	}
//
//	cases := []struct {
//		name        string
//		mobile      string
//		want        *userpb.UserEntity
//		wantErrCode codes.Code
//	}{
//		{
//			name:   "exist_user",
//			mobile: "mobile_1",
//			want:   data.ToResp(users[0]),
//		},
//		{
//			name:   "another_exist_user",
//			mobile: "mobile_2",
//			want:   data.ToResp(users[1]),
//		},
//		{
//			name:        "not_exist_user",
//			mobile:      "999",
//			wantErrCode: codes.NotFound,
//		},
//	}
//
//	for _, cc := range cases {
//		t.Run(cc.name, func(t *testing.T) {
//			code := codes.OK
//			got, err := s.GetUserByMobile(context.Background(), &userpb.GetUserByMobileRequest{Mobile: cc.mobile})
//			if err != nil {
//				if s, ok := status.FromError(err); ok {
//					code = s.Code()
//				} else {
//					t.Errorf("operation failed: %v", err)
//				}
//			}
//			if code != cc.wantErrCode {
//				t.Errorf("wrong err code: want %d, got %d", cc.wantErrCode, code)
//			}
//
//			if diff := cmp.Diff(cc.want, got, protocmp.Transform()); diff != "" {
//				t.Errorf("result differs; -want +got: %s", diff)
//			}
//		})
//	}
//}
//
//func TestService_GetUsers(t *testing.T) {
//	s := newService(t)
//
//	users := []*data.User{
//		{
//			NickName: "user_1",
//			Mobile:   "mobile_1",
//		},
//		{
//			NickName: "user_2",
//			Mobile:   "mobile_2",
//		},
//		{
//			NickName: "user_3",
//			Mobile:   "mobile_3",
//		},
//	}
//
//	res := s.DB.Create(&users)
//	if res.Error != nil {
//		t.Fatalf("cannot create users: %v", res.Error)
//	}
//
//	cases := []struct {
//		name        string
//		number      uint32
//		size        uint32
//		want        *userpb.GetUsersResponse
//		wantErrCode codes.Code
//	}{
//		{
//			name: "get_all_users",
//			want: &userpb.GetUsersResponse{
//				Total: int64(len(users)),
//				Users: []*userpb.UserEntity{
//					data.ToResp(users[0]),
//					data.ToResp(users[1]),
//					data.ToResp(users[2]),
//				},
//			},
//		},
//		{
//			name:   "get_all_users_in_two_pages",
//			number: 1,
//			size:   2,
//			want: &userpb.GetUsersResponse{
//				Total: int64(len(users)),
//				Users: []*userpb.UserEntity{
//					data.ToResp(users[0]),
//					data.ToResp(users[1]),
//				},
//			},
//		},
//		{
//			name:   "get_users_out_of_range",
//			number: 3,
//			size:   2,
//			want: &userpb.GetUsersResponse{
//				Total: int64(len(users)),
//			},
//		},
//	}
//
//	for _, cc := range cases {
//		t.Run(cc.name, func(t *testing.T) {
//			var code codes.Code
//			got, err := s.GetUsers(context.Background(), &userpb.GetUsersRequest{
//				Number: cc.number,
//				Size:   cc.size,
//			})
//			if err != nil {
//				if s, ok := status.FromError(err); ok {
//					code = s.Code()
//				} else {
//					t.Errorf("operation failed: %v", err)
//				}
//			}
//			if cc.wantErrCode != code {
//				t.Errorf("wrong err code: want %d, got %d", cc.wantErrCode, code)
//			}
//			if diff := cmp.Diff(cc.want, got, protocmp.Transform()); diff != "" {
//				t.Errorf("result differs; -want +got: %s", diff)
//			}
//		})
//	}
//}
//
//func TestService_UpdateUser(t *testing.T) {
//	s := newService(t)
//
//	users := []*data.User{
//		{
//			NickName: "user_1",
//			Mobile:   "mobile_1",
//		},
//		{
//			NickName: "user_2",
//			Mobile:   "mobile_2",
//		},
//	}
//
//	res := s.DB.Create(&users)
//	if res.Error != nil {
//		t.Fatalf("cannot create users: %v", res.Error)
//	}
//
//	cases := []struct {
//		name        string
//		nickName    string
//		birth       int64
//		wantErrCode codes.Code
//	}{
//		{
//			name: "no_change",
//		},
//		{
//			name:     "change_nickname",
//			nickName: "changed_name",
//		},
//	}
//
//	for _, cc := range cases {
//		t.Run(cc.name, func(t *testing.T) {
//			code := codes.OK
//			_, err := s.UpdateUser(context.Background(), &userpb.UpdateUserRequest{
//				Id:       "1",
//				NickName: cc.nickName,
//				Gender:   "",
//				Birthday: cc.birth,
//			})
//			if err != nil {
//				if s, ok := status.FromError(err); ok {
//					code = s.Code()
//				} else {
//					t.Errorf("operation failed: %v", err)
//				}
//			}
//			if code != cc.wantErrCode {
//				t.Errorf("wrong err code: want %d, got %d", cc.wantErrCode, code)
//			}
//		})
//	}
//}
//
//func TestService_CheckPassword(t *testing.T) {
//	s := newService(t)
//
//	cases := []struct {
//		name          string
//		passwd        string
//		genEcpPwdFunc func(passwd string) string
//		wantErrCode   codes.Code
//	}{
//		{
//			name:   "passwd",
//			passwd: "passwd_1",
//			genEcpPwdFunc: func(passwd string) string {
//				salt, ecdPwd := password.Encode(passwd, s.PwdOpts)
//				ecpPwd := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, ecdPwd)
//				return ecpPwd
//			},
//		},
//		{
//			name:   "another_passwd",
//			passwd: "passwd_2",
//			genEcpPwdFunc: func(passwd string) string {
//				salt, ecdPwd := password.Encode(passwd, s.PwdOpts)
//				ecpPwd := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, ecdPwd)
//				return ecpPwd
//			},
//		},
//		{
//			name:   "bad_passwd",
//			passwd: "passwd_2",
//			genEcpPwdFunc: func(passwd string) string {
//				salt, ecdPwd := password.Encode(passwd, s.PwdOpts)
//				b := []byte(ecdPwd)
//				b[rand.Intn(len(b))]++
//				ecpPwd := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, string(b))
//				return ecpPwd
//			},
//			wantErrCode: codes.Unauthenticated,
//		},
//	}
//
//	for _, cc := range cases {
//		t.Run(cc.name, func(t *testing.T) {
//			code := codes.OK
//			_, err := s.CheckPassword(context.Background(), &userpb.CheckPasswordRequest{
//				Password:          cc.passwd,
//				EncryptedPassword: cc.genEcpPwdFunc(cc.passwd),
//			})
//			if err != nil {
//				if s, ok := status.FromError(err); ok {
//					code = s.Code()
//				} else {
//					t.Errorf("operation failed: %v", err)
//				}
//			}
//			if code != cc.wantErrCode {
//				t.Errorf("wrong err code: want %d, got %d", cc.wantErrCode, code)
//			}
//		})
//	}
//}
//
//func newService(t *testing.T) *UserService {
//	db, err := mysqltesting.NewDB()
//	if err != nil {
//		t.Fatalf("cannot get database: %v", err)
//	}
//
//	err = db.AutoMigrate(&data.User{})
//	if err != nil {
//		t.Fatalf("cannot create tables: %v", err)
//	}
//
//	logger, err := zap.NewDevelopment()
//	if err != nil {
//		t.Fatalf("cannot create logger: %v", err)
//	}
//
//	s := &UserService{
//		Repo: db,
//		PwdOpts: &password.Options{
//			SaltLen:      10,
//			Iterations:   100,
//			KeyLen:       32,
//			HashFunction: sha512.New,
//		},
//		Logger: logger,
//	}
//
//	return s
//}
//
//func TestMain(m *testing.M) {
//	os.Exit(mysqltesting.RunWithMysqlInDocker(m))
//}
