package service

import (
	"context"
	"github.com/codedancewth/public_project/internal/dao"
	"github.com/codedancewth/public_project/proto/public_project"
)

func (s *PublicProject) GetUserInfo(ctx context.Context,
	req *public_project.GetUserInfoReq) (resp *public_project.GetUserInfoRep, err error) {
	userName := ctx.Value("userName").(string)
	user, err := dao.GetUser(DImp.gDB, userName)
	if err != nil {
		return nil, err
	}
	return &public_project.GetUserInfoRep{
		User: &public_project.User{
			Id:          user.ID,
			UserName:    user.UserName,
			UserAccount: user.UserAccount,
			Status:      int32(user.Status),
			CreatedTime: user.CreatedTime,
			UpdatedTime: user.UpdatedTime,
			IsDeleted:   int32(user.IsDeleted),
		},
	}, nil
}
