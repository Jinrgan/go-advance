package api

import "time"

type User struct {
	ID        string        `json:"id"`
	NickName  string        `json:"nick_name"`
	Token     string        `json:"token"`
	ExpiredIn time.Duration `json:"expired_in"`
}

type UserLoginResponse struct {
	NickName  string        `json:"nick_name"`
	Token     string        `json:"token"`
	ExpiredIn time.Duration `json:"expired_in"`
}

type Captcha struct {
	ID      string `json:"captcha_id"`
	PicPath string `json:"pic_path"`
}

type RegisterForm struct {
	Mobile   string `form:"mobile" json:"mobile" binding:"required,mobile"` // 手机号码格式有规范可寻，自定义 validator
	Password string `form:"password" json:"password" binding:"required,min=3,max=20"`
	Code     string `form:"code" json:"code" binding:"required,min=6,max=6"`
}

type LoginForm struct {
	ID       string `form:"id" json:"id" binding:"required,id"` // 手机号码格式有规范可寻，自定义 validator
	Password string `form:"password" json:"password" binding:"required,min=3,max=20"`
}

type LoginByMobileForm struct {
	Mobile    string `form:"mobile" json:"mobile" binding:"required,mobile"` // 手机号码格式有规范可寻，自定义 validator
	Password  string `form:"password" json:"password" binding:"required,min=3,max=20"`
	Captcha   string `form:"captcha" json:"captcha" binding:"required,min=5,max=5"`
	CaptchaID string `form:"captcha_id" json:"captcha_id" binding:"required"`
}

type SendSmsForm struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required,mobile"` // 手机号码格式有规范可寻，自定义 validator
	Type   string `form:"type" json:"type" binding:"required,oneof=1 2"`
}
