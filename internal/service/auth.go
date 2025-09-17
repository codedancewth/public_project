package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/codedancewth/public_project/internal/dao"
	"github.com/codedancewth/public_project/pkg/utils"
	"github.com/codedancewth/public_project/proto/public_project"
	"github.com/sirupsen/logrus"
	"time"
)

// Login 接口登录
func (s *PublicProject) Login(ctx context.Context,
	req *public_project.LoginReq) (resp *public_project.LoginRep, err error) {

	user, err := dao.GetUser(DImp.gDB, req.Username)
	if err != nil {
		return nil, err
	}

	// 验证密码
	if err = user.CheckPassword(req.Password); err != nil {
		logrus.Infof("check password error %+v", err)
		return
	}

	// 生成JWT令牌
	token, err := utils.GenerateToken(uint(user.ID), user.UserName)
	if err != nil {
		logrus.Infof("generate token error %+v", err)
		return
	}

	// 生成刷新令牌
	refreshToken, _ := utils.GenerateRefreshToken()

	if refreshToken == "" {
		err = errors.New("refresh token fail")
		return
	}

	// 存储刷新令牌到Redis，有效期7天
	err = DImp.rc.Set(ctx, fmt.Sprintf("refresh_token:%s", refreshToken),
		user.UserName, 7*24*time.Hour).Err()
	if err != nil {
		logrus.Infof("refresh token token error %+v", err)
		return
	}

	return &public_project.LoginRep{
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresIn:    utils.TokenExpiredIn,
	}, nil
}

func (s *PublicProject) RefreshToken(ctx context.Context,
	req *public_project.RefreshTokenReq) (resp *public_project.RefreshTokenRep, err error) {

	// 验证刷新令牌
	username, err := DImp.rc.Get(ctx, fmt.Sprintf("refresh_token:%s", req.RefreshToken)).Result()

	if err != nil {
		return
	}

	// 获取用户信息
	user, err := dao.GetUser(DImp.gDB, username)
	if err != nil {
		return nil, err
	}

	if user == nil || user.ID == 0 {
		return &public_project.RefreshTokenRep{}, nil
	}

	// 生成新的JWT令牌
	token, err := utils.GenerateToken(uint(user.ID), user.UserName)
	if err != nil {
		return
	}
	return &public_project.RefreshTokenRep{
		Token:     token,
		ExpiresIn: utils.TokenExpiredIn,
	}, nil
}

func (s *PublicProject) Logout(ctx context.Context,
	req *public_project.LogoutReq) (resp *public_project.LogoutRep, err error) {
	if req.RefreshToken != "" {
		ctx := context.Background()
		// 删除刷新令牌
		_ = DImp.rc.Del(ctx, fmt.Sprintf("refresh_token:%s", req.RefreshToken)).Err()
	}

	return &public_project.LogoutRep{Message: "success logout"}, nil
}
