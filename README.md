# spec-cc-0104

基于 Go-Zero 微服务架构的数据源管理系统，采用 AI 辅助的规范驱动开发 (SDD) 模式。

## 📋 项目简介

本项目是一个数据源配置管理系统，提供数据源的创建、查询、管理等核心功能。支持多种数据源类型（MySQL、PostgreSQL、Redis、MongoDB、SQL Server），采用微服务架构设计，支持容器化部署。

## 🚀 技术栈

- **语言**: Go 1.24+
- **框架**: Go-Zero v1.9+
- **数据库**: MySQL 8.0
- **ORM**: GORM v1.31+
- **架构**: 微服务 (API/RPC/Job/Consumer)
- **通用库**: idrm-go-base v0.1.0+
- **部署**: Docker + Kubernetes

## 📁 项目结构

```
spec-cc-0104/
├── api/                      # API 服务
│   ├── doc/                  # API 定义 (.api 文件)
│   │   ├── api.api          # 主 API 定义
│   │   ├── base.api         # 基础类型定义
│   │   └── datasource/      # 数据源模块 API
│   ├── etc/                  # 配置文件
│   │   ├── api.yaml         # API 服务配置
│   │   └── spec_cc_0104_api.yaml
│   └── internal/             # 内部实现
│       ├── handler/          # 请求处理层
│       │   ├── datasource/  # 数据源处理器
│       │   └── routes.go    # 路由注册
│       ├── logic/            # 业务逻辑层
│       │   └── datasource/  # 数据源业务逻辑
│       ├── svc/              # 服务上下文
│       ├── types/            # 类型定义
│       └── config/           # 配置结构
├── model/                    # 数据模型层
│   └── datasource/          # 数据源模型
├── migrations/               # 数据库迁移文件
│   └── datasource/          # 数据源表结构
├── specs/                    # SDD 规格文档
│   └── datasource/          # 数据源功能规格
├── deploy/                   # 部署配置
│   ├── docker/              # Docker 配置
│   │   ├── Dockerfile.api   # API 服务镜像
│   │   ├── docker-compose.yaml
│   │   └── build.sh         # 构建脚本
│   └── k8s/                 # Kubernetes 配置
│       ├── base/            # 基础配置
│       └── overlays/        # 环境覆盖
│           ├── dev/         # 开发环境
│           └── prod/        # 生产环境
├── rpc/                      # RPC 服务（预留）
├── consumer/                 # 消息消费者（预留）
├── job/                      # 定时任务（预留）
├── Makefile                  # 构建脚本
├── go.mod                    # Go 模块定义
└── README.md                 # 项目说明
```

## 🛠️ 快速开始

### 环境要求

- Go 1.24+
- MySQL 8.0+
- Docker & Docker Compose (可选)
- Make (可选)

### 本地开发

1. **克隆项目**
```bash
git clone <repository-url>
cd spec-cc-0104
```

2. **安装依赖**
```bash
go mod download
```

3. **配置数据库**
```bash
# 创建数据库
mysql -u root -p
CREATE DATABASE idrm CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

# 执行迁移
mysql -u root -p idrm < migrations/datasource/datasource.sql
```

4. **配置环境变量**
```bash
# 复制配置文件
cp api/etc/api.yaml api/etc/api.yaml.local

# 编辑配置文件，设置数据库连接等信息
```

5. **运行服务**
```bash
# 方式1: 使用 Makefile
make run

# 方式2: 直接运行
go run api/spec_cc_0104_api.go -f api/etc/api.yaml
```

服务默认运行在 `http://localhost:8888`

### Docker 部署

1. **使用 Docker Compose**
```bash
cd deploy/docker

# 设置环境变量
export DB_PASSWORD=your_password
export ACCESS_SECRET=your_secret

# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f api
```

2. **构建镜像**
```bash
make docker-build
```

## 📚 API 文档

### 数据源管理

#### 创建数据源
```http
POST /api/v1/datasource
Content-Type: application/json

{
  "name": "生产数据库",
  "type": "mysql",
  "host": "192.168.1.100",
  "port": 3306,
  "database": "production",
  "username": "admin",
  "password": "password123",
  "description": "生产环境主数据库",
  "status": "enabled",
  "sort_order": 1
}
```

#### 查询数据源列表
```http
GET /api/v1/datasource?offset=1&limit=10&sort=created_at&direction=desc&keyword=mysql
```

### 生成 Swagger 文档

```bash
# 生成 JSON 格式
make swagger

# 生成 YAML 格式
make swagger-yaml
```

文档将生成在 `api/doc/swagger/` 目录下。

## 🔧 常用命令

### 开发命令

```bash
make api           # 从 .api 文件生成 API 代码
make swagger       # 生成 Swagger 文档
make gen           # 生成 API 代码 + Swagger 文档
make fmt           # 格式化代码
make lint          # 代码检查
make test          # 运行测试
make build         # 编译二进制文件
make run           # 运行服务
make clean         # 清理构建产物
make deps          # 安装依赖
```

### Docker 命令

```bash
make docker-build  # 构建 Docker 镜像
make docker-run    # 运行 Docker 容器
make docker-stop   # 停止 Docker 容器
make docker-push   # 推送镜像到仓库
```

### Kubernetes 命令

```bash
make k8s-deploy        # 部署到 K8s (默认: dev)
make k8s-deploy-dev    # 部署到开发环境
make k8s-deploy-prod   # 部署到生产环境
make k8s-manifest      # 查看生成的 Manifest
make k8s-delete        # 删除 K8s 部署
make k8s-status        # 查看部署状态
```

## 🏗️ 架构设计

### 分层架构

```
┌─────────────────┐
│   Handler 层    │  参数绑定、校验、调用 Logic
├─────────────────┤
│   Logic 层      │  业务逻辑、事务管理
├─────────────────┤
│   Model 层      │  数据访问 (GORM)
├─────────────────┤
│   Database      │  MySQL
└─────────────────┘
```

### 职责划分

| 层 | 职责 | 禁止行为 |
|---|------|----------|
| Handler | 参数绑定、校验、调用 Logic | 包含业务逻辑 |
| Logic | 业务逻辑、事务管理 | 直接操作 HTTP |
| Model | 数据访问 (GORM) | 包含业务逻辑 |

## 📝 开发规范

### 命名规范

- 文件名: `snake_case.go`
- 包名: `lowercase`
- 结构体: `PascalCase`
- 方法: `PascalCase`
- 变量: `camelCase`
- 常量: `UPPER_SNAKE_CASE`

### 错误处理

使用 `idrm-go-base/errorx` 进行错误处理：

```go
import "github.com/jinguoxing/idrm-go-base/errorx"

// 使用预定义错误码
if user == nil {
    return nil, errorx.NewWithCode(errorx.ErrCodeNotFound)
}

// 自定义业务错误码
if user.Status == 0 {
    return nil, errorx.New(30102, "用户已禁用")
}
```

### 日志规范

使用 `logx` 进行日志记录，自动包含 traceId：

```go
import "github.com/zeromicro/go-zero/core/logx"

logx.WithContext(ctx).Infof("user login: %s", phone)
```

### 数据库规范

- 所有表使用 UUID v7 作为主键
- 使用软删除 (`deleted_at`)
- 自动时间戳 (`created_at`, `updated_at`)

## 🔐 配置说明

### 环境变量

主要配置项通过环境变量设置：

```yaml
# 服务配置
API_PORT=8888
ENVIRONMENT=dev
LOG_LEVEL=info

# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_NAME=idrm
DB_USER=root
DB_PASSWORD=your_password

# Redis 配置（可选）
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# 认证配置
ACCESS_SECRET=your_secret_key
ACCESS_EXPIRE=7200
```

### 配置文件

配置文件位于 `api/etc/api.yaml`，支持环境变量替换。

## 🧪 测试

```bash
# 运行所有测试
make test

# 运行特定包的测试
go test -v ./api/internal/logic/datasource/...

# 运行测试并查看覆盖率
go test -v -cover ./...
```

## 📦 部署

### Docker 部署

```bash
# 构建镜像
make docker-build

# 使用 docker-compose
cd deploy/docker
docker-compose up -d
```

### Kubernetes 部署

```bash
# 部署到开发环境
make k8s-deploy-dev

# 部署到生产环境
make k8s-deploy-prod

# 查看状态
make k8s-status
```

## 🔄 SDD 工作流程

本项目遵循 Spec-Driven Development (规范驱动开发) 流程：

1. **Context** - 理解项目规范和上下文
2. **Specify** - 使用 EARS 格式定义需求
3. **Design** - 创建技术方案 (API/DDL/Model)
4. **Tasks** - 拆分任务清单
5. **Implement** - 按任务顺序编码实现

详细流程请参考 `CLAUDE.md` 文件。

## 📄 许可证

[待添加]

## 👥 贡献

[待添加]

## 📞 联系方式

[待添加]

---

**注意**: 本项目遵循严格的开发规范，请参考 `CLAUDE.md` 了解详细的开发流程和约束条件。

