package service

import (
	"context"
	"github.com/codedancewth/public_project/config"
	"github.com/codedancewth/public_project/internal/cache"
	"github.com/codedancewth/public_project/proto/public_project"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var DImp *PublicProject

type PublicProject struct {
	public_project.UnimplementedAppServiceServer
	gDB *gorm.DB // 主库
	//mongoDB       *mongo.Client
	rc *redis.Client
}

type option func(imp *PublicProject)

// NewAppService 初始化逻辑
func NewAppService() public_project.AppServiceServer {
	imp := &PublicProject{}
	options := []option{
		withMysqlDB(),
		withRedis(),
		withCacheInit(),
		//withConnectMongo(),
	}
	for _, o := range options {
		o(imp)
	}
	DImp = imp

	return imp
}

func withMysqlDB() option {
	return func(imp *PublicProject) {
		var err error
		mysqlDB, err := gorm.Open(mysql.Open(config.DefaultConfig.MySQL.GetDSN()), &gorm.Config{})
		if err != nil {
			return
		}

		db, err := mysqlDB.DB()
		if err != nil {
			return
		}
		db.SetMaxOpenConns(30)
		db.SetMaxIdleConns(50)

		imp.gDB = mysqlDB
	}
}

func withRedis() option {
	return func(imp *PublicProject) {
		r := redis.NewClient(&redis.Options{
			Network:     "tcp",
			Addr:        config.DefaultConfig.Redis.GetAddr(),
			Password:    config.DefaultConfig.Redis.Password,
			DB:          config.DefaultConfig.Redis.DB,
			PoolSize:    config.DefaultConfig.Redis.PoolSize,
			IdleTimeout: time.Duration(config.DefaultConfig.Redis.IdleTimeout) * time.Second,
		})
		if err := r.Ping(context.Background()).Err(); err != nil {
			logrus.Panicf("withRedis ping redis failed. %s", err)
		}
		imp.rc = r
	}
}

func withCacheInit() option {
	return func(imp *PublicProject) {
		cache.InitLocalCache("", time.Minute)
	}
}
