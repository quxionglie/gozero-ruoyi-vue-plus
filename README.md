# Go-Zero Ruoyi Vue Plus
 
使用go-zero实现[RuoYi-Vue-Plus](https://github.com/dromara/RuoYi-Vue-Plus), 没有使用gorm框架，直接使用go-zero自带的框架。
暂不实现多租户功能。

所有代码的迁移基本使用Cursor AI进行代码转换的，[cursor_chat.txt](cursor_chat.txt)是Cursor对话记录。

对应java版本：[RuoYi-Vue-Plus](https://github.com/dromara/RuoYi-Vue-Plus)  tag: v5.5.2

对应ui版本：[plus-ui](https://github.com/JavaLionLi/plus-ui) tag: v5.5.2-v2.5.2 , 需更换RSA 公钥/私钥 VITE_APP_RSA_PUBLIC_KEY、VITE_APP_RSA_PRIVATE_KEY

RuoYi-Vue-Plus与plus-ui可查看当前项目对应git分支。

基于 [go-zero](https://github.com/zeromicro/go-zero) 框架的高性能后端服务项目，集成 MySQL 和 Redis，提供 RESTful API 服务。
 
## 🛠️ 技术栈

- **框架**: [go-zero](https://github.com/zeromicro/go-zero) v1.6.1
- **数据库**: MySQL 5.7+
- **缓存**: Redis 6.0+
- **语言**: Go 1.21+

## 📁 项目结构

```
gozero-ruoyi-vue-plus/
├── api/                          # API 定义文件（goctl 使用）
│   ├── auth.api                 # 认证相关接口定义
│   ├── main.api                 # 主 API 文件（导入其他 API）
│   ├── monitor.api              # 监控相关接口定义
│   ├── resource.api             # 资源管理接口定义（OSS）
│   ├── sse.api                  # Server-Sent Events 接口定义
│   └── system.api               # 系统管理接口定义
├── etc/                          # 配置文件目录
│   └── admin-api.yaml           # 主配置文件
├── internal/                     # 内部代码（不对外暴露）
│   ├── config/                   # 配置结构定义
│   │   └── config.go            # 配置结构体
│   ├── handler/                  # HTTP 请求处理器
│   │   ├── auth/                # 认证相关处理器
│   │   ├── monitor/             # 监控相关处理器
│   │   ├── resource/            # 资源管理处理器（OSS）
│   │   ├── sse/                 # SSE 处理器
│   │   ├── sys/                 # 系统管理处理器
│   │   └── routes.go            # 路由注册
│   ├── logic/                    # 业务逻辑层
│   │   ├── auth/                # 认证相关业务逻辑
│   │   ├── monitor/             # 监控相关业务逻辑
│   │   ├── resource/            # 资源管理业务逻辑（OSS）
│   │   ├── sse/                 # SSE 业务逻辑
│   │   └── sys/                 # 系统管理业务逻辑
│   ├── middleware/               # 中间件
│   │   └── corsmiddleware.go    # CORS 中间件
│   ├── model/                    # 数据模型层
│   │   └── sys/                 # 系统相关数据模型
│   ├── oss/                      # 对象存储服务
│   │   ├── client.go            # OSS 客户端接口定义
│   │   ├── manager.go           # OSS 管理器
│   │   ├── minio_client.go      # MinIO 客户端实现
│   │   └── minio_client_ext.go  # MinIO 客户端扩展
│   ├── svc/                      # 服务上下文
│   │   └── service_context.go   # 服务上下文，包含 DB、Redis、模型等
│   ├── types/                    # 类型定义
│   │   └── types.go             # API 请求/响应类型定义
│   └── util/                     # 工具函数
│       ├── captcha.go           # 验证码工具
│       ├── decrypt.go           # 解密工具
│       ├── jwt.go               # JWT 工具
│       ├── response.go          # 响应工具
│       ├── snowflake.go         # 雪花算法 ID 生成
│       ├── sse_manager.go       # SSE 管理器
│       └── validator.go         # 验证器
├── admin.go                      # 应用程序入口
├── cursor_chat.txt               # Cursor AI 对话记录
├── go.mod                        # Go 模块定义
├── Makefile                      # Make 命令文件
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
git clone https://github.com/quxionglie/gozero-ruoyi-vue-plus.git
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

```bash
go install github.com/zeromicro/go-zero/tools/goctl@latest
```

### 4. 配置数据库

#### 创建 MySQL 数据库

```sql
CREATE DATABASE ry-vue DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
```

#### 修改配置文件

编辑 `etc/admin-api.yaml` 文件，修改数据库连接信息：

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

服务启动后，默认监听 `0.0.0.0:58888`，可以测试健康检查接口：

```bash
curl 'http://localhost:58888/auth/code' \
  -H 'Accept: application/json, text/plain, */*' \
  -H 'clientid: e5cd7e4891bf95d1d19206ce24a7b32e' \
  --insecure
```

预期响应：

```json
{
  "message": "pong",
  "status": "ok"
}
```
## 📖 开发指南
 
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

⭐ 如果这个项目对你有帮助，请给个 Star 支持一下！