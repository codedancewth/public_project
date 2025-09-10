package service

import (
	"context"
	"fmt"
	"github.com/codedancewth/public_project/internal/dao"
	"github.com/codedancewth/public_project/internal/proto/public_project"
	"github.com/codedancewth/public_project/pkg/utils"
)

// Hello 测试接口的正常性
func (s *PublicProject) Hello(ctx context.Context,
	req *public_project.HelloReq) (resp *public_project.HelloRsp, err error) {
	fmt.Println("that's is ok")

	users, err := dao.GetUserList(DImp.gDB)
	if err != nil {
		return &public_project.HelloRsp{}, nil
	}

	fmt.Println(utils.UtilJsonMarshal(users))

	return &public_project.HelloRsp{}, nil
}
