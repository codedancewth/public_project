package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/codedancewth/public_project/internal/service"
	"github.com/codedancewth/public_project/proto/public_project"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"net/http"
)

var (
	listenInternal = flag.String("listen-internal", ":8080", "listen of internal http service")
)

const GrpcPort = 22066
const ProjectName = "public_project" // 服务名

func getGrpcPort() int64 {
	return GrpcPort
}

func gracefulStop(server *grpc.Server) {
	return
}

func main() {
	// 获取启动参数
	var configFile string
	flag.StringVar(&configFile, "f", "", "Configuration file.")
	flag.StringVar(&configFile, "c", "", "Configuration file.")
	flag.Parse()

	// 初始化对应的中间件和配置
	imp := service.NewAppService()

	// start gRPC server
	grpcServer := grpc.NewServer()
	public_project.RegisterAppServiceServer(grpcServer, imp)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", getGrpcPort()))
	if err != nil {
		logrus.Panicf("failed to listen: %+v", err)
	}
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logrus.Panicf("grpc server start failed, err: %+v", err)
		}
	}()

	// start grpc-gateway http server
	go func() {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		gwMux := runtime.NewServeMux()
		dialOpts := []grpc.DialOption{grpc.WithInsecure()}
		endpoint := "127.0.0.1:" + fmt.Sprintf("%d", getGrpcPort())
		if err := public_project.RegisterAppServiceHandlerFromEndpoint(ctx, gwMux, endpoint, dialOpts); err != nil {
			logrus.Panicf("register http gateway failed, err: %+v", err)
		}

		if err := http.ListenAndServe(*listenInternal, gwMux); err != nil {
			logrus.Panicf("http server start failed, err: %+v", err)
		}
	}()

	select {}
}
