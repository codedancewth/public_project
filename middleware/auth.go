package middleware

import (
	"context"
	"github.com/codedancewth/public_project/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

// AuthInterceptor gRPC 一元拦截器
func AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 排除登录和刷新令牌等不需要认证的方法
		if WhiteRouter(info.FullMethod) {
			return handler(ctx, req)
		}

		// 从 metadata 中获取 token
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "缺少认证信息")
		}

		// 获取 authorization 头
		authHeaders := md["authorization"]
		if len(authHeaders) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "未提供认证令牌")
		}

		// 解析 Bearer token
		authHeader := authHeaders[0]
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return nil, status.Errorf(codes.Unauthenticated, "令牌格式错误")
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 验证 token
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "无效的令牌: %v", err)
		}

		// 将用户信息添加到上下文
		newCtx := context.WithValue(ctx, "userId", claims.UserID)
		newCtx = context.WithValue(newCtx, "userName", claims.Username)

		return handler(newCtx, req)
	}
}

func WhiteRouter(path string) bool {
	if strings.Contains(path, "Login") || strings.Contains(path, "Refresh") {
		return true
	}

	return false
}
