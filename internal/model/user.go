package model

import "golang.org/x/crypto/bcrypt"

// User 用户模型
type User struct {
	ID           int64  `gorm:"primaryKey"`
	UserName     string `gorm:"column:user_name" json:"user_name"`
	UserAccount  string `gorm:"column:user_account" json:"user_account"`
	UserPassword string `gorm:"column:user_password" json:"user_password"`
	Status       int8   `gorm:"column:status" json:"status"`
	CreatedTime  int64  `gorm:"column:created_time" json:"created_time"`
	UpdatedTime  int64  `gorm:"column:updated_time" json:"updated_time"`
	IsDeleted    int8   `gorm:"column:is_deleted" json:"is_deleted"`
}

func (u *User) Table() string {
	return "user"
}

// HashPassword 加密密码
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.UserPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.UserPassword = string(hashedPassword)
	return nil
}

// CheckPassword 验证密码
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.UserPassword), []byte(password))
}
