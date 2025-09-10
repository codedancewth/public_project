package service

import (
	"context"
	"fmt"
	"github.com/codedancewth/public_project/internal/dao"
	"github.com/codedancewth/public_project/internal/proto/public_project"
	"github.com/codedancewth/public_project/pkg/utils"
	"time"
)

// Hello 测试接口的正常性
func (s *PublicProject) Hello(ctx context.Context,
	req *public_project.HelloReq) (resp *public_project.HelloRsp, err error) {

	// 数据库测试
	users, err := dao.GetUserList(DImp.gDB)
	if err != nil {
		return &public_project.HelloRsp{}, nil
	}

	fmt.Println(fmt.Sprintf("test for mysql [%v]", utils.UtilJsonMarshal(users)))

	// redis测试
	DImp.rc.Set(ctx, "hello", "hello redis", time.Second*3)
	getHello := DImp.rc.Get(ctx, "hello").String()
	fmt.Println(fmt.Sprintf("test for redis [%v]", getHello))

	return &public_project.HelloRsp{}, nil
}
