# Go-Zero Ruoyi Vue Plus

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Go-Zero](https://img.shields.io/badge/go--zero-1.6.1-00ADD8?style=flat)](https://github.com/zeromicro/go-zero)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

基于 [go-zero](https://github.com/zeromicro/go-zero) 框架的高性能后端服务项目，集成 MySQL 和 Redis，提供 RESTful API 服务。

## ✨ 特性

- 🚀 **高性能**: 基于 go-zero 框架，提供出色的性能表现
- 🏗️ **架构清晰**: 遵循 go-zero 最佳实践，代码结构清晰易维护
- 💾 **数据持久化**: 集成 MySQL 数据库，支持事务和连接池
- ⚡ **缓存支持**: 集成 Redis 缓存，提升数据访问性能
- 📝 **完整日志**: 内置日志系统，支持文件输出和日志轮转
- 🔧 **易于配置**: 基于 YAML 配置文件，支持灵活配置
- 🛠️ **开发友好**: 提供 Makefile 命令，简化开发和部署流程

## 🛠️ 技术栈

- **框架**: [go-zero](https://github.com/zeromicro/go-zero) v1.6.1
- **数据库**: MySQL 5.7+
- **缓存**: Redis 6.0+
- **语言**: Go 1.21+

## 📁 项目结构

```
gozero-ruoyi-vue-plus/
├── etc/                          # 配置文件目录
│   └── admin-api.yaml           # 主配置文件
├── internal/                     # 内部代码（不对外暴露）
│   ├── config/                   # 配置结构定义
│   │   └── config.go            # 配置结构体
│   ├── handler/                  # HTTP 请求处理器
│   │   └── routes.go            # 路由注册和处理器
│   └── svc/                      # 服务上下文
│       └── servicecontext.go    # 服务上下文，包含 DB 和 Redis 连接
├── docs/                         # 文档目录
├── logs/                         # 日志文件目录（自动创建）
├── bin/                          # 编译输出目录
├── admin.go                      # 应用程序入口
├── go.mod                        # Go 模块定义
├── go.sum                        # Go 模块校验和
├── Makefile                      # Make 命令文件
├── .gitignore                    # Git 忽略文件
└── README.md                     # 项目说明文档
```

## 🚀 快速开始

### 前置要求

- Go 1.21 或更高版本
- MySQL 5.7+ 或更高版本
- Redis 6.0+ 或更高版本
- Make（可选，用于便捷命令）

### 1. 克隆项目

```bash
git clone <repository-url>
cd gozero-ruoyi-vue-plus
```

### 2. 安装依赖

```bash
# 使用 Makefile
make deps

# 或直接使用 go 命令
go mod tidy
go mod download
```

### 3. 安装 goctl 工具（可选）

go-zero 提供了强大的代码生成工具 goctl，建议安装：

```bash
go install github.com/zeromicro/go-zero/tools/goctl@latest
```

### 4. 配置数据库

#### 创建 MySQL 数据库

```sql
CREATE DATABASE IF NOT EXISTS ruoyi DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
```

#### 修改配置文件

编辑 `etc/admin-api.yaml` 文件，修改数据库连接信息：

```yaml
Name: gozero-ruoyi-vue-plus
Host: 0.0.0.0
Port: 58888
Timeout: 60000

# MySQL配置
Mysql:
  DataSource: root:password@tcp(127.0.0.1:3306)/ruoyi?charset=utf8mb4&parseTime=True&loc=Local

# Redis配置
Redis:
  Host: 127.0.0.1:6379
  Type: node
  Pass: ""

# 日志配置
Log:
  ServiceName: gozero-ruoyi-vue-plus
  Mode: console
  Path: logs
  Level: info
  Compress: true
  KeepDays: 7
  StackCooldownMillis: 100
```

**注意**: 请根据实际情况修改数据库用户名、密码、主机地址和端口。

### 5. 启动服务

#### 方式一：使用 Makefile（推荐）

```bash
# 运行项目
make run

# 或先编译再运行
make build
./bin/gozero-ruoyi-vue-plus -f etc/admin-api.yaml
```

#### 方式二：直接使用 Go 命令

```bash
# 开发模式运行
go run admin.go -f etc/admin-api.yaml

# 编译后运行
go build -o bin/gozero-ruoyi-vue-plus admin.go
./bin/gozero-ruoyi-vue-plus -f etc/admin-api.yaml
```

### 6. 验证服务

服务启动后，默认监听 `0.0.0.0:8888`，可以测试健康检查接口：

```bash
curl http://localhost:8888/api/ping
```

预期响应：

```json
{
  "message": "pong",
  "status": "ok"
}
```

## ⚙️ 配置说明

### 服务配置

- `Name`: 服务名称
- `Host`: 服务监听地址，`0.0.0.0` 表示监听所有网络接口
- `Port`: 服务监听端口

### MySQL 配置

- `DataSource`: MySQL 连接字符串
  - 格式: `用户名:密码@tcp(主机:端口)/数据库名?参数`
  - 示例: `root:password@tcp(127.0.0.1:3306)/ruoyi?charset=utf8mb4&parseTime=True&loc=Local`
  - 参数说明:
    - `charset=utf8mb4`: 字符集
    - `parseTime=True`: 解析时间类型
    - `loc=Local`: 时区设置

### Redis 配置

- `Host`: Redis 服务器地址和端口（格式: `host:port`）
- `Type`: 连接类型
  - `node`: 单节点模式
  - `sentinel`: 哨兵模式
  - `cluster`: 集群模式
- `Pass`: Redis 密码（如果设置了密码）
- `DB`: 数据库编号（0-15）

### 日志配置

- `ServiceName`: 服务名称（用于日志标识）
- `Mode`: 日志模式
  - `file`: 文件模式
  - `console`: 控制台模式
  - `volume`: 容器卷模式
- `Path`: 日志文件路径
- `Level`: 日志级别（debug/info/error/severe）
- `Compress`: 是否压缩日志文件
- `KeepDays`: 日志保留天数
- `StackCooldownMillis`: 堆栈冷却时间（毫秒）

## 📖 开发指南

### 添加新的 API 路由

1. **在 `internal/handler/routes.go` 中注册路由**:

```go
func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
    server.AddRoutes(
        []rest.Route{
            {
                Method:  http.MethodGet,
                Path:    "/api/users",
                Handler: GetUsersHandler(serverCtx),
            },
            {
                Method:  http.MethodPost,
                Path:    "/api/users",
                Handler: CreateUserHandler(serverCtx),
            },
        },
    )
}
```

2. **创建对应的 Handler 函数**:

```go
func GetUsersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 实现业务逻辑
        httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{
            "users": []string{},
        })
    }
}
```

### 数据库操作

通过 `ServiceContext` 中的 `DB` 进行数据库操作：

```go
// 在 handler 中使用
func SomeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 使用 svcCtx.DB 执行数据库操作
        rows, err := svcCtx.DB.Query("SELECT * FROM users")
        // ...
    }
}
```

### Redis 操作

通过 `ServiceContext` 中的 `RedisConn` 进行 Redis 操作：

```go
// 在 handler 中使用
func SomeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 使用 svcCtx.RedisConn 执行 Redis 操作
        err := svcCtx.RedisConn.Set("key", "value")
        // ...
    }
}
```

### 使用 goctl 生成代码

go-zero 提供了强大的代码生成工具，可以快速生成 API 定义、模型代码等：

```bash
# 生成 API 定义
goctl api go -api api/user.api -dir . -style gozero

# 从数据库生成模型代码
goctl model mysql datasource -url="user:password@tcp(127.0.0.1:3306)/ruoyi" -table="user" -dir="./model" -c
```

更多 goctl 使用说明请参考 [go-zero 文档](https://go-zero.dev/cn/docs/design/overview)。

## 🛠️ Makefile 命令

项目提供了便捷的 Makefile 命令：

```bash
# 查看所有可用命令
make help

# 安装依赖
make deps

# 构建项目
make build

# 运行项目
make run

# 格式化代码
make fmt

# 代码检查
make vet

# 清理构建文件和日志
make clean
```

## 📝 API 文档

### 健康检查

**接口**: `GET /api/ping`

**描述**: 服务健康检查接口

**请求示例**:

```bash
curl http://localhost:8888/api/ping
```

**响应示例**:

```json
{
  "message": "pong",
  "status": "ok"
}
```

## 🚢 部署

### 编译二进制文件

```bash
# 开发环境
go build -o bin/gozero-ruoyi-vue-plus admin.go

# 生产环境（交叉编译）
# Linux
GOOS=linux GOARCH=amd64 go build -o bin/gozero-ruoyi-vue-plus admin.go

# Windows
GOOS=windows GOARCH=amd64 go build -o bin/gozero-ruoyi-vue-plus.exe admin.go

# macOS
GOOS=darwin GOARCH=amd64 go build -o bin/gozero-ruoyi-vue-plus admin.go
```

### Docker 部署

创建 `Dockerfile`:

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o gozero-ruoyi-vue-plus admin.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

COPY --from=builder /app/gozero-ruoyi-vue-plus .
COPY --from=builder /app/etc ./etc

EXPOSE 58888
CMD ["./gozero-ruoyi-vue-plus", "-f", "etc/admin-api.yaml"]
```

构建和运行:

```bash
docker build -t gozero-ruoyi-vue-plus .
docker run -p 8888:8888 gozero-ruoyi-vue-plus
```

## ⚠️ 注意事项

1. **生产环境配置**: 部署到生产环境时，请务必修改默认密码和敏感配置
2. **环境变量**: 建议使用环境变量或配置中心管理敏感信息，避免硬编码
3. **日志管理**: 日志文件会自动保存到 `logs` 目录，注意定期清理
4. **数据库连接**: 确保数据库连接池配置合理，避免连接泄漏
5. **Redis 连接**: 生产环境建议使用 Redis 密码认证
6. **安全建议**: 
   - 启用 HTTPS
   - 配置 CORS
   - 实现请求限流
   - 添加认证和授权

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 📄 许可证

本项目采用 MIT 许可证。详情请参阅 [LICENSE](LICENSE) 文件。

## 🔗 相关链接

- [go-zero 官方文档](https://go-zero.dev/cn/docs/)
- [go-zero GitHub](https://github.com/zeromicro/go-zero)
- [Go 官方文档](https://golang.org/doc/)

## 📞 联系方式

如有问题或建议，请通过以下方式联系：

- 提交 Issue
- 发送邮件

---

⭐ 如果这个项目对你有帮助，请给个 Star 支持一下！