package biz

import (
	"fmt"
	"go-advance/week04/app/user/api"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/go-redis/redis"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

type SmsHandler struct {
	client      *dysmsapi.Client
	request     *requests.CommonRequest
	ExpireInSec int
	redis       *redis.Client
}

func NewSmsHandler(clt *dysmsapi.Client, exp int, rdsClt *redis.Client) (*SmsHandler, error) {
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = "cn-beijing"
	request.QueryParams["SignName"] = "幕学生鲜"    // 阿里云验证过的项目名 自己设置
	request.QueryParams["TemplateCode"] = "xxx" // 阿里云的短信模板号 自己设置

	return &SmsHandler{
		client:      clt,
		request:     request,
		ExpireInSec: exp,
		redis:       rdsClt,
	}, nil
}

func mustGenerateSmsCode(length int) string {
	code, err := generateSmsCode(length)
	if err != nil {
		panic(err)
	}

	return code
}

func generateSmsCode(length int) (string, error) {
	rand.Seed(time.Now().UnixNano())

	var bud strings.Builder
	for i := 0; i < length; i++ {
		_, err := fmt.Fprintf(&bud, "%d", rand.Intn(10))
		if err != nil {
			return "", fmt.Errorf("cannot print to builder: %v", err)
		}
	}

	return bud.String(), nil
}

func (s *SmsHandler) Send(ctx *gin.Context) {
	var form api.SendSmsForm
	err := ctx.ShouldBind(&form)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	code := mustGenerateSmsCode(6)
	s.request.QueryParams["PhoneNumbers"] = form.Mobile                // 手机号
	s.request.QueryParams["TemplateParam"] = "{\"code\":" + code + "}" // 短信模板中的验证码内容 自己生成   之前试过直接返回，但是失败，加上code成功。
	resp, err := s.client.ProcessCommonRequest(s.request)
	if err != nil {
		zap.L().Error("cannot process common request", zap.Error(err))
		return
	}

	err = s.client.DoAction(s.request, resp)
	if err != nil {
		zap.L().Error("cannot do action", zap.Error(err))
		return
	}

	s.redis.Set(form.Mobile, code, time.Duration(s.ExpireInSec)*time.Second)

	ctx.String(http.StatusOK, "send sms successful")

	return
}
