package service

import (
	"context"
	"fmt"
	"github.com/codedancewth/public_project/internal/proto/public_project"
)

// Hello 测试接口的正常性
func (s *PublicProject) Hello(ctx context.Context,
	req *public_project.HelloReq) (resp *public_project.HelloRsp, err error) {
	fmt.Println("that's is ok")
	return &public_project.HelloRsp{}, nil
}
