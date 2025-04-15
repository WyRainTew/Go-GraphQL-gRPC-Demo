package grpc

import (
	"context"
	"log"
	"net"

	"Demo/service/database"
	pb "Demo/service/grpc/pb" // 导入生成的 protobuf 代码
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserServer 实现 gRPC UserService 服务
type UserServer struct {
	pb.UnimplementedUserServiceServer // 必须嵌入此类型
	db                                *database.Database
}

// NewUserServer 创建 UserServer 实例
func NewUserServer(db *database.Database) *UserServer {
	return &UserServer{db: db}
}

// GetUserInfo 实现 GetUserInfo RPC 方法
func (s *UserServer) GetUserInfo(ctx context.Context, req *pb.UserInfoRequest) (*pb.UserInfoResponse, error) {
	// 获取请求参数
	userID := req.UserId

	// 从数据库查询用户
	user, err := s.db.GetUserByID(userID)
	if err != nil {
		log.Printf("查询用户失败: %v", err)
		return nil, status.Errorf(codes.NotFound, "未找到用户: %v", err)
	}

	// 构建响应
	return &pb.UserInfoResponse{
		Id:   user.ID,
		Name: user.Name,
		Age:  int32(user.Age),
		Sex:  user.Sex,
	}, nil
}

// Server 封装 gRPC 服务器
type Server struct {
	server   *grpc.Server
	listener net.Listener
}

// NewServer 创建新的 gRPC 服务器
func NewServer(address string, db *database.Database) (*Server, error) {
	// 创建 TCP 监听器
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	// 创建 gRPC 服务器
	server := grpc.NewServer()

	// 注册服务实现
	pb.RegisterUserServiceServer(server, NewUserServer(db))

	return &Server{
		server:   server,
		listener: listener,
	}, nil
}

// Start 启动 gRPC 服务器
func (s *Server) Start() error {
	log.Printf("gRPC 服务器启动在 %s", s.listener.Addr())
	return s.server.Serve(s.listener)
}

// Stop 优雅地停止服务器
func (s *Server) Stop() {
	s.server.GracefulStop()
}
