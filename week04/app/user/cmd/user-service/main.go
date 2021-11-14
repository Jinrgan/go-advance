package main

import (
	"crypto/sha512"
	"fmt"
	userpb "go-advance/week04/app/user/api/gen/v1"
	"go-advance/week04/app/user/config"
	biznacos "go-advance/week04/app/user/config/nacos"
	"go-advance/week04/app/user/internal/biz"
	"go-advance/week04/app/user/internal/biz/code"
	"go-advance/week04/app/user/internal/data"
	"go-advance/week04/app/user/internal/service"
	"go-advance/week04/pkg/nacos"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"

	"github.com/hashicorp/consul/api"

	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/anaskhan96/go-password-encoder"
	uuid "github.com/satori/go.uuid"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	ormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}

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
			CacheDir:            "app/user/cmd/user-service/tmp/nacos/cache",
			NotLoadCacheAtStart: true,
			LogDir:              "app/user/cmd/user-service/tmp/nacos/log",
			RotateTime:          "1h",
			MaxAge:              3,
			LogLevel:            "debug",
		}))
	if err != nil {
		zap.L().Error("cannot create config client", zap.Error(err))
	}
	cfgTor := config.Configurator(&biznacos.Configurator{
		DataID: "service",
		Group:  "dev",
		Client: confClt,
	})
	var cfg config.Service
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
		mysql.Open(fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.Mysql.User, cfg.Mysql.Password, cfg.Mysql.Addr, cfg.Mysql.DB)),
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

	lis, err := net.Listen("tcp", cfg.Addr+":0")
	if err != nil {
		logger.Fatal("cannot listen", zap.Error(err))

	}

	s := grpc.NewServer()
	userpb.RegisterUserServiceServer(s, &service.UserService{
		Biz: &biz.UserBiz{
			Coder: &code.PasswdCoder{PwdOpts: &password.Options{
				SaltLen:      10,
				Iterations:   100,
				KeyLen:       32,
				HashFunction: sha512.New,
			}},
			Repo: &data.MySQL{DB: db},
		},
		Repo: &data.MySQL{DB: db},
	})
	grpc_health_v1.RegisterHealthServer(s, health.NewServer())

	regCfg := api.DefaultConfig()
	//cfg.Address = fmt.Sprintf("%s:%d", srvConf.Consul.Host, srvConf.Consul.Port)
	clt, err := api.NewClient(regCfg)
	if err != nil {
		logger.Error("cannot create consul client", zap.Error(err))
		return
	}

	srvPort, err := strconv.Atoi(strings.Split(lis.Addr().String(), ":")[1])
	if err != nil {
		logger.Error("cannot convey port", zap.Error(err))
		return
	}
	srvID := uuid.NewV4().String()
	err = clt.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      srvID,
		Name:    cfg.Name,
		Tags:    []string{"mxshop", "bobby"},
		Port:    srvPort,
		Address: cfg.Addr,
		Check: &api.AgentServiceCheck{
			GRPC:                           lis.Addr().String(),
			Timeout:                        "5s",
			Interval:                       "5s",
			DeregisterCriticalServiceAfter: "10s",
		},
	})
	defer func() {
		err := clt.Agent().ServiceDeregister(srvID)
		if err != nil {
			logger.Error("cannot deregister service", zap.String("service", cfg.Name+srvID), zap.Error(err))
			return
		}
	}()
	if err != nil {
		panic(err)
	}

	logger.Info("server started", zap.String("name", cfg.Name), zap.String("Addr", lis.Addr().String()))
	logger.Sugar().Fatal(s.Serve(lis))
}

func GenUnmFn(cfg *config.Service) func([]byte, biznacos.UnmarshalFn) error {
	return func(b []byte, unmarshal biznacos.UnmarshalFn) error {
		err := unmarshal(b, &cfg)
		if err != nil {
			return fmt.Errorf("cannot unmarshal config: %v", err)
		}

		return nil
	}
}
