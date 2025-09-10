package service

import (
	"context"
	"github.com/codedancewth/public_project/proto/public_project"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

const DbURL = "root:123456@tcp(127.0.0.1:3306)/project?charset=utf8mb4&parseTime=True&loc=Local"

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
		mysqlDB, err := gorm.Open(mysql.Open(DbURL), &gorm.Config{})
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
			Addr:        "127.0.0.1:6379",
			Password:    "",
			DB:          0,
			PoolSize:    10,
			IdleTimeout: time.Minute * 30,
		})
		if err := r.Ping(context.Background()).Err(); err != nil {
			logrus.Panicf("withRedis ping redis failed. %s", err)
		}
		imp.rc = r
	}
}
