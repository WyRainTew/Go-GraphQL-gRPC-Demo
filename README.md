# GraphQL + gRPC + SQLite 用户信息服务

这是一个基于 Go 的微服务演示项目，展示了如何使用 GraphQL、gRPC 和 SQLite 构建一个简单的用户信息查询服务。项目采用了三层架构，实现了从前端到数据库的完整流程。

![架构图](https://raw.githubusercontent.com/WyRainBow/assets/main/graphql-grpc-sqlite-arch.png)

## 项目架构

该项目采用了三层架构：

1. **GraphQL 层**：提供 HTTP API，接收客户端查询
2. **gRPC 层**：实现核心业务逻辑
3. **数据库层**：使用 SQLite 存储用户数据

架构设计：
```
┌───────────┐     ┌───────────┐     ┌───────────┐
│  GraphQL  │ --> │   gRPC    │ --> │  SQLite   │
│  (8080)   │     │  (50051)  │     │           │
└───────────┘     └───────────┘     └───────────┘
```

## 技术栈

- **后端**: Go 1.21+
- **API**: GraphQL (github.com/graphql-go/graphql)
- **微服务**: gRPC
- **ORM**: GORM
- **数据库**: SQLite

## 功能

- 查询用户信息: `userInfo(userId: String!): User`
- 用户数据包含: id, name, age, sex

## 快速开始

### 前置条件

- Go 1.21+
- Protocol Buffers (protoc)
- 对应的 Go 插件:
  ```bash
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
  ```

### 安装

1. 克隆仓库
```bash
git clone https://github.com/WyRainBow/go-graphql-grpc-demo.git
cd go-graphql-grpc-demo
```

2. 安装依赖
```bash
go mod tidy
```

3. 重新生成 Protocol Buffers 代码 (可选)
```bash
cd service/grpc
protoc --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    user.proto
cd ../..
```

### 运行

启动服务:
```bash
go run main.go
```

默认情况下:
- GraphQL API 在 `http://localhost:8080/query`
- GraphQL Playground 在 `http://localhost:8080/playground`
- gRPC 服务在 `:50051`

### 自定义配置

你可以通过命令行参数自定义配置:
```bash
go run main.go -db=./mydb.db -grpc=:50052 -http=:8081
```

## 使用示例

访问 GraphQL Playground: http://localhost:8080/playground

执行以下查询:

```graphql
{
  userInfo(userId: "aaa") {
    id
    name
    age
    sex
  }
}
```

这将返回预设的测试用户数据:

```json
{
  "data": {
    "userInfo": {
      "id": "aaa",
      "name": "测试用户",
      "age": 30,
      "sex": "男"
    }
  }
}
```

## 项目结构

```
Demo/
├── main.go                 # 主程序
├── service/
│   ├── database/           # 数据库层
│   │   └── database.go     # SQLite数据库操作
│   ├── grpc/               # gRPC 服务
│   │   ├── server.go       # gRPC服务器实现
│   │   ├── user.proto      # Protocol Buffers定义
│   │   └── pb/             # 生成的 protobuf 代码
│   └── graphql/            # GraphQL API
│       ├── resolver.go     # GraphQL解析器
│       ├── server.go       # GraphQL HTTP服务器
│       └── static/         # 静态资源
│           └── playground.html  # GraphQL Playground界面
├── go.mod                  # Go模块定义
└── go.sum                  # 依赖校验和
```

## 流程说明

1. 客户端向 GraphQL 服务器发送查询请求
2. GraphQL 服务器解析请求并调用 gRPC 客户端
3. gRPC 客户端向 gRPC 服务器发送请求
4. gRPC 服务器查询 SQLite 数据库获取数据
5. 数据通过 gRPC、GraphQL 层层返回给客户端

## 许可证

MIT
