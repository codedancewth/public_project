package service

import (
	"context"
	"fmt"
	"github.com/codedancewth/public_project/internal/dao"
	"github.com/codedancewth/public_project/pkg/utils"
	"github.com/codedancewth/public_project/proto/public_project"
	"net/http"
	"time"
)

// Login 接口登录
func (s *PublicProject) Login(ctx context.Context,
	req *public_project.LoginReq) (resp *public_project.LoginRep, err error) {

	users, err := dao.GetUser(DImp.gDB, req.Username)
	if err != nil {
		return nil, err
	}

	// 验证密码
	if err := user.CheckPassword(loginReq.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 生成JWT令牌
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法生成令牌"})
		return
	}

	// 生成刷新令牌
	refreshToken := utils.GenerateRefreshToken()

	// 存储刷新令牌到Redis，有效期7天
	err = DImp.rc.Set(ctx, fmt.Sprintf("refresh_token:%s", refreshToken),
		user.Username, 7*24*time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法存储会话"})
		return
	}

	return &public_project.LoginRep{
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresIn:    15 * 60,
	}, nil
}
