package graphql

import (
	"context"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/graphql-go/handler"
)

// Server GraphQL HTTP 服务器
type Server struct {
	schema     *handler.Handler
	address    string
	httpServer *http.Server
	resolver   *Resolver
}

// NewServer 创建 GraphQL 服务器
func NewServer(grpcAddress string, httpAddress string) (*Server, error) {
	// 创建解析器
	resolver, err := NewResolver(grpcAddress)
	if err != nil {
		return nil, err
	}

	// 创建 Schema
	schema, err := CreateSchema(resolver)
	if err != nil {
		return nil, err
	}

	// 创建 HTTP 处理器
	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: false, // 禁用默认的GraphiQL
	})

	return &Server{
		schema:   h,
		address:  httpAddress,
		resolver: resolver,
	}, nil
}

// Start 启动 GraphQL 服务器
func (s *Server) Start() error {
	mux := http.NewServeMux()
	
	// GraphQL API 端点
	mux.Handle("/query", s.schema)
	
	// 自定义 Playground 界面
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 如果是根路径，重定向到playground
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/playground", http.StatusMovedPermanently)
			return
		}
		http.NotFound(w, r)
	})
	
	// 提供静态Playground页面
	mux.HandleFunc("/playground", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("service", "graphql", "static", "playground.html"))
	})
	
	s.httpServer = &http.Server{
		Addr:    s.address,
		Handler: mux,
	}
	
	log.Printf("GraphQL 服务器启动在 %s/query", s.address)
	log.Printf("GraphQL Playground 可访问：%s/playground", s.address)
	return s.httpServer.ListenAndServe()
}

// Stop 优雅地停止服务器
func (s *Server) Stop() error {
	log.Println("正在关闭 GraphQL 服务器...")
	
	// 创建一个上下文，设置 5 秒超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// 优雅关闭 HTTP 服务器
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Printf("关闭 HTTP 服务器出错: %v", err)
		return err
	}
	
	// 关闭 gRPC 连接
	if err := s.resolver.Close(); err != nil {
		log.Printf("关闭 gRPC 连接出错: %v", err)
		return err
	}
	
	log.Println("GraphQL 服务器已关闭")
	return nil
}
