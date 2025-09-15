package dao

import (
	"github.com/codedancewth/public_project/internal/model"
	"gorm.io/gorm"
	"log"
)

// GetUserList 获取用户的基础数据列表
func GetUserList(db *gorm.DB) (users []*model.User, err error) {
	err = db.Table((&model.User{}).Table()).Find(&users).Error
	if err != nil {
		log.Fatalf("find user list error [%v]", err)
	}
	return users, err
}

// GetUser 获取用户
func GetUser(db *gorm.DB, userName string) (users *model.User, err error) {
	err = db.Table((&model.User{}).Table()).Where("user_name", userName).Find(&users).Error
	if err != nil {
		log.Fatalf("get user error [%v]", err)
	}
	return users, err
}
