package graphql

import (
	"context"
	"log"

	pb "Demo/service/grpc/pb"
	"github.com/graphql-go/graphql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Resolver GraphQL 解析器，连接到 gRPC 服务
type Resolver struct {
	grpcClient pb.UserServiceClient
	conn       *grpc.ClientConn // 添加连接字段用于关闭
}

// NewResolver 创建新的 GraphQL 解析器
func NewResolver(grpcAddress string) (*Resolver, error) {
	// 创建 gRPC 连接
	conn, err := grpc.Dial(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	// 创建 gRPC 客户端
	client := pb.NewUserServiceClient(conn)

	return &Resolver{
		grpcClient: client,
		conn:       conn,
	}, nil
}

// Close 关闭 gRPC 连接
func (r *Resolver) Close() error {
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}

// 创建 GraphQL Schema
func CreateSchema(resolver *Resolver) (graphql.Schema, error) {
	// 定义用户类型
	userType := graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"age": &graphql.Field{
				Type: graphql.Int,
			},
			"sex": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	// 定义查询类型
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"userInfo": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"userId": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// 从参数中获取用户 ID
					userID, ok := p.Args["userId"].(string)
					if !ok {
						return nil, nil
					}

					// 调用 gRPC 客户端获取用户信息
					resp, err := resolver.grpcClient.GetUserInfo(context.Background(), &pb.UserInfoRequest{
						UserId: userID,
					})

					if err != nil {
						log.Printf("gRPC 调用失败: %v", err)
						return nil, err
					}

					// 返回用户信息
					return map[string]interface{}{
						"id":   resp.Id,
						"name": resp.Name,
						"age":  resp.Age,
						"sex":  resp.Sex,
					}, nil
				},
			},
		},
	})

	// 创建 Schema
	return graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
	})
}
