package main

import (
	"crypto/sha512"
	"fmt"
	userpb "go-advance/week04/app/user/api/gen/v1"
	"go-advance/week04/app/user/config"
	confnacos "go-advance/week04/app/user/config/nacos"
	"go-advance/week04/app/user/internal/biz"
	"go-advance/week04/app/user/internal/biz/captcha"
	"go-advance/week04/app/user/internal/biz/code"
	"go-advance/week04/app/user/internal/biz/middleware"
	"go-advance/week04/app/user/internal/biz/token"
	"go-advance/week04/app/user/internal/data"
	"go-advance/week04/app/user/internal/service"
	"go-advance/week04/pkg/auth"
	"go-advance/week04/pkg/nacos"
	"go-advance/week04/pkg/web"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"

	"github.com/hashicorp/consul/api"

	ormlogger "gorm.io/gorm/logger"

	"github.com/anaskhan96/go-password-encoder"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/go-redis/redis"

	"github.com/mojocn/base64Captcha"

	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	// 1. S() 以取一个全局的 sugar，可以让我自己设置一个全局的 logger
	// 2. 日志是分级别的，debug、info、warn、error、fetal
	// 3. S 函数和 L 函数很有用，提供了一个全局的安全访问 logger 的途径
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("cannot create logger: %v", err)
	}
	zap.ReplaceGlobals(logger)

	// get service config
	confClt, err := clients.CreateConfigClient(nacos.CreateCltOpts(
		[]*constant.ServerConfig{
			{
				IpAddr: "localhost",
				Port:   8848,
			},
		},
		&constant.ClientConfig{
			TimeoutMs:           5000,
			NamespaceId:         "0fb2fa26-8c9d-43a1-8eba-7b4638ecc446",
			CacheDir:            "app/user/cmd/user-interface/tmp/nacos/cache",
			NotLoadCacheAtStart: true,
			LogDir:              "app/user/cmd/user-interface/tmp/nacos/log",
			RotateTime:          "1h",
			MaxAge:              3,
			LogLevel:            "debug",
		}))
	if err != nil {
		zap.L().Error("cannot create config client", zap.Error(err))
	}
	cfgTor := config.Configurator(&confnacos.Configurator{
		DataID: "web",
		Group:  "dev",
		Client: confClt,
	})
	var cfg config.Server
	err = cfgTor.GetConfig(GenUnmFn(&cfg))
	if err != nil {
		zap.L().Error("cannot get config", zap.Error(err))
	}
	zap.L().Info("success to get config", zap.Any("conf", cfg))
	go func() {
		err := cfgTor.Listen(GenUnmFn(&cfg))
		if err != nil {
			zap.L().Error("fail to listen", zap.Error(err))
		}

		return
	}()

	// mysql
	db, err := gorm.Open(
		mysql.Open(fmt.Sprintf("root:root@tcp(%s)/mxshop_user_srv?charset=utf8mb4&parseTime=True&loc=Local", cfg.Mysql.Addr)),
		&gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			Logger: ormlogger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags),
				ormlogger.Config{
					SlowThreshold: time.Second,
					Colorful:      true,
					LogLevel:      ormlogger.Info,
				}),
		})
	if err != nil {
		logger.Fatal("cannot connect to database", zap.Error(err))
	}

	// 代码侵入性很强 中间件
	trans, err := middleware.NewTrans("zh")
	if err != nil {
		zap.L().Error("fail to init trans", zap.Error(err))
		return
	}
	err = middleware.RegisterValidation(trans)
	if err != nil {
		zap.L().Error("cannot 01register validation", zap.Error(err))
		return
	}
	_ = &web.ErrWrapper{Trans: trans}

	// base router
	capStore := base64Captcha.DefaultMemStore
	r := gin.Default()
	r.Use(middleware.AllowAccess) // 配置跨域
	rg := r.Group("/u/v1")
	cch := captcha.NewCaptcha(capStore)
	bg := rg.Group("base")
	smsClt, err := dysmsapi.NewClientWithAccessKey("cn-beijing", "xxxx", "xxx")
	if err != nil {
		zap.L().Error("cannot create client with access key", zap.Error(err))
		return
	}
	rds := redis.NewClient(&redis.Options{Addr: cfg.Redis.Addr})
	smsH, err := biz.NewSmsHandler(smsClt, cfg.Sms.Expire, rds)
	if err != nil {
		zap.L().Error("cannot create sms handler", zap.Error(err))
		return
	}
	bg.POST("captcha", cch.Gen)
	bg.POST("send_sms", smsH.Send)

	// user handler
	pkFile, err := os.Open("app/user/private.key")
	if err != nil {
		logger.Fatal("cannot open private key", zap.Error(err))
	}
	pkBytes, err := ioutil.ReadAll(pkFile)
	if err != nil {
		logger.Fatal("cannot read private key", zap.Error(err))
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(pkBytes)
	if err != nil {
		logger.Fatal("cannot parse private key", zap.Error(err))
	}

	// get config from consul register
	conCfg := api.DefaultConfig()
	conCfg.Address = cfg.Consul.Addr
	clt, err := api.NewClient(conCfg)
	if err != nil {
		panic(err)
	}

	srv, err := clt.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%s"`, cfg.ServiceName))
	if err != nil {
		panic(err)
	}

	var srvAddr string
	for _, v := range srv {
		srvAddr = fmt.Sprintf("%s:%d", v.Address, v.Port)
		break
	}
	conn, err := grpc.Dial(srvAddr, grpc.WithInsecure())
	if err != nil {
		zap.L().Error("cannot connect to user service", zap.Error(err))
		return
	}
	srvClt := userpb.NewUserServiceClient(conn)
	uh := &service.Interface{
		CaptchaVerifier: &captcha.Verifier{Redis: rds},
		TokenGenerator:  token.NewJWTTokenGen("app/user/private.key", privateKey),
		TokenExpire:     2 * time.Hour,
		CaptchaStore:    capStore,
		Biz: &biz.UserBiz{
			Coder: &code.PasswdCoder{PwdOpts: &password.Options{
				SaltLen:      10,
				Iterations:   100,
				KeyLen:       32,
				HashFunction: sha512.New,
			}},
			Repo: &data.MySQL{DB: db},
		},
		Repo: nil,
	}
	mid, err := auth.NewMiddleware("pkg/auth/public.key", srvClt)
	if err != nil {
		zap.L().Error("cannot create middleware", zap.Error(err))
		return
	}
	urg := rg.Group("user")
	{
		urg.POST("01register", uh.Register)
		urg.GET("list", mid.HandleReq, mid.HandleAdminReq, uh.GetUsers)
		urg.POST("pwd_login", uh.Login)
	}

	zap.S().Error(r.Run(cfg.Addr))
}

func GenUnmFn(cfg *config.Server) func([]byte, confnacos.UnmarshalFn) error {
	return func(b []byte, unmarshal confnacos.UnmarshalFn) error {
		err := unmarshal(b, &cfg)
		if err != nil {
			return fmt.Errorf("cannot unmarshal config: %v", err)
		}

		return nil
	}
}
