package service

import (
	"github.com/codedancewth/public_project/internal/proto/public_project"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const DbURL = "root:123456@tcp(127.0.0.1:3306)/project?charset=utf8mb4&parseTime=True&loc=Local"

var DImp *PublicProject

type PublicProject struct {
	public_project.UnimplementedAppServiceServer
	gDB *gorm.DB // 主库
	//mongoDB       *mongo.Client
	//rc            *redis.Client
}

type option func(imp *PublicProject)

// NewAppService 初始化逻辑
func NewAppService() public_project.AppServiceServer {
	imp := &PublicProject{}
	options := []option{
		withEasyDB(),
		//withRedis(),
		//withConnectMongo(),
	}
	for _, o := range options {
		o(imp)
	}
	DImp = imp

	return imp
}

func withEasyDB() option {
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
