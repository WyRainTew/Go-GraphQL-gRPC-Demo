package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"Demo/service/database"
	"Demo/service/graphql"
	grpcServer "Demo/service/grpc"
)

func main() {
	// 命令行参数
	dbPath := flag.String("db", "user.db", "SQLite 数据库文件路径")
	grpcAddr := flag.String("grpc", ":50051", "gRPC 服务器地址")
	httpAddr := flag.String("http", ":8080", "GraphQL HTTP 服务器地址")
	flag.Parse()

	log.Println("启动用户信息服务...")

	// 初始化数据库
	db, err := database.NewDatabase(*dbPath)
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("关闭数据库失败: %v", err)
		}
	}()
	log.Println("数据库初始化成功")

	// 创建 gRPC 服务器
	server, err := grpcServer.NewServer(*grpcAddr, db)
	if err != nil {
		log.Fatalf("创建 gRPC 服务器失败: %v", err)
	}
	log.Println("gRPC 服务器创建成功")

	// 启动 gRPC 服务器
	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("启动 gRPC 服务器失败: %v", err)
		}
	}()

	// 创建并启动 GraphQL 服务器
	graphqlServer, err := graphql.NewServer(*grpcAddr, *httpAddr)
	if err != nil {
		log.Fatalf("创建 GraphQL 服务器失败: %v", err)
	}
	log.Println("GraphQL 服务器创建成功")

	// 启动 GraphQL 服务器
	go func() {
		if err := graphqlServer.Start(); err != nil {
			log.Fatalf("启动 GraphQL 服务器失败: %v", err)
		}
	}()

	log.Printf("服务已启动:\n- gRPC: %s\n- GraphQL: %s/graphql\n", *grpcAddr, *httpAddr)
	log.Println("使用 Ctrl+C 停止服务")

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("正在关闭服务...")
	
	// 关闭 GraphQL 服务器
	if err := graphqlServer.Stop(); err != nil {
		log.Printf("关闭 GraphQL 服务器出错: %v", err)
	}
	
	// 关闭 gRPC 服务器
	server.Stop()
	
	log.Println("所有服务已关闭")
}
