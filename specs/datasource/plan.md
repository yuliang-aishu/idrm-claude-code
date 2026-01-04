# DataSource Technical Plan

> **Branch**: `feature/datasource`
> **Spec Path**: `specs/datasource/`
> **Created**: 2026-01-04
> **Status**: Draft

---

## Summary

基于 Go-Zero 微服务架构实现数据源管理模块，支持多种数据源类型的统一管理。采用 UUID v7 作为主键，AES-256-GCM 加密敏感信息，Single-tenant 模式。核心功能包括数据源列表查询（分页、搜索、排序）和新增创建（连接测试验证）。遵循 IDRM 分层架构：Handler → Logic → Model → Database。

---

## Technical Context

| Item | Value |
|------|-------|
| **Language** | Go 1.24+ |
| **Framework** | Go-Zero v1.9+ |
| **Storage** | MySQL 8.0 |
| **Cache** | Redis 7.0 |
| **ORM** | GORM / SQLx |
| **Testing** | go test |
| **Common Lib** | idrm-go-base v0.1.0+ |
| **Encryption** | AES-256-GCM |
| **Primary Key** | UUID v7 |
| **Tenant Mode** | Single-tenant |
| **Auth Model** | Single admin role |

---

## 通用库 (idrm-go-base)

**安装**:
```bash
go get github.com/jinguoxing/idrm-go-base@latest
```

### 模块初始化

| 模块 | 初始化方式 |
|------|-----------|
| validator | `validator.Init()` 在 main.go |
| telemetry | `telemetry.Init(cfg)` 在 main.go |
| response | `httpx.SetErrorHandler(response.ErrorHandler)` |
| middleware | `rest.WithMiddlewares(...)` |
| errorx | 在 Logic 层使用错误码 |

### 自定义错误码

| 功能 | 范围 | 位置 |
|------|------|------|
| 数据源管理 | 30400-30499 | `api/internal/errorx/codes.go` |

**错误码定义**:
- 30400: 参数错误（必填参数为空、无效类型等）
- 30409: 资源冲突（名称重复）
- 30411: 资源不存在（查询不存在的 ID）
- 30413: 连接测试失败

### 第三方库确认

> 如需使用通用库以外的第三方库，请在此列出并说明原因:

| 库 | 原因 | 确认状态 |
|----|------|----------|
| github.com/google/uuid | 生成 UUID v7 主键 | ✅ 已确认（项目规范要求） |
| golang.org/x/crypto/... | AES-256-GCM 加密实现 | ✅ 已确认（加密需求） |

---

## Go-Zero 开发流程

按以下顺序完成技术设计和代码生成：

| Step | 任务 | 方式 | 产出 |
|------|------|------|------|
| 1 | 定义 API 文件 | AI 实现 | `api/doc/datasource/datasource.api` |
| 2 | 生成 Handler/Types | goctl | `api/internal/handler/`, `types/` |
| 3 | 定义 DDL 文件 | AI 手写 | `migrations/datasource/datasource.sql` |
| 4 | 实现 Model 接口 | AI 手写 | `model/datasource/datasource/` |
| 5 | 实现 Logic 层 | AI 实现 | `api/internal/logic/datasource/` |

> ⚠️ **重要**：goctl 必须在 `api/doc/api.api` 入口文件上执行，不能针对单个功能文件！

**goctl 命令**:
```bash
# 步骤1：在 api/doc/api.api 中 import datasource 模块
# 步骤2：执行 goctl 生成代码（针对整个项目）
goctl api go -api api/doc/api.api -dir api/ --style=go_zero --type-group
```

---

## File Structure

### 文件产出清单

| 序号 | 文件 | 生成方式 | 位置 |
|------|------|----------|------|
| 1 | API 文件 | AI 实现 | `api/doc/datasource/datasource.api` |
| 2 | DDL 文件 | AI 实现 | `migrations/datasource/datasource.sql` |
| 3 | Handler | goctl 生成 | `api/internal/handler/datasource/` |
| 4 | Types | goctl 生成 | `api/internal/types/` |
| 5 | Logic | AI 实现 | `api/internal/logic/datasource/` |
| 6 | Model | AI 实现 | `model/datasource/datasource/` |

### 代码结构

```
api/internal/
├── handler/datasource/
│   ├── create_datasource_handler.go    # goctl 生成
│   ├── list_datasource_handler.go
│   └── routes.go
├── logic/datasource/
│   ├── create_datasource_logic.go      # AI 实现
│   └── list_datasource_logic.go
├── types/
│   └── types.go                       # goctl 生成
└── svc/
    └── servicecontext.go              # 手动维护

model/datasource/datasource/
├── interface.go                       # 接口定义
├── types.go                           # 数据结构
├── vars.go                            # 常量/错误
├── factory.go                         # ORM 工厂
├── gorm_dao.go                        # GORM 实现
└── sqlx_model.go                      # SQLx 实现
```

---

## Architecture Overview

遵循 IDRM 分层架构：

```
HTTP Request → Handler → Logic → Model → Database
```

| 层级 | 职责 | 最大行数 |
|------|------|----------|
| Handler | 解析参数、格式化响应、调用 Logic | 30 |
| Logic | 业务逻辑实现、事务管理、连接测试 | 50 |
| Model | 数据访问、CRUD 操作、加密/解密 | 50 |

**职责说明**:
- **Handler**: 接收 HTTP 请求，参数绑定和校验，调用 Logic 层，返回统一响应格式
- **Logic**: 实现数据源业务规则（名称唯一性校验、连接测试、状态管理），调用 Model 层
- **Model**: 数据持久化，密码 AES-256-GCM 加密，UUID v7 生成，软删除实现

---

## Interface Definitions

根据 Spec 中的 Data Considerations 定义 Model 接口：

```go
type DataSourceModel interface {
    // 插入数据源
    Insert(ctx context.Context, data *DataSource) (*DataSource, error)

    // 根据 ID 查询数据源
    FindOne(ctx context.Context, id string) (*DataSource, error)

    // 更新数据源
    Update(ctx context.Context, data *DataSource) error

    // 软删除数据源
    Delete(ctx context.Context, id string) error

    // 列表查询（支持分页、搜索、筛选）
    FindList(ctx context.Context, query *DataSourceQuery) ([]*DataSource, int64, error)

    // 检查名称是否重复
    CheckNameExists(ctx context.Context, name string, excludeId ...string) (bool, error)

    // 连接测试
    TestConnection(ctx context.Context, config *DataSource) error

    // 事务支持
    WithTx(tx interface{}) DataSourceModel
    Trans(ctx context.Context, fn func(ctx context.Context, model DataSourceModel) error) error
}

// 查询参数
type DataSourceQuery struct {
    Offset    int    `form:"offset,default=1"`
    Limit     int    `form:"limit,default=10"`
    Keyword   string `form:"keyword,optional"`
    Status    string `form:"status,optional"`
    Sort      string `form:"sort,default=created_at"`
    Direction string `form:"direction,default=desc"`
}
```

---

## Data Model

### DDL

**位置**: `migrations/datasource/datasource.sql`

```sql
CREATE TABLE `datasource` (
    `id` CHAR(36) NOT NULL COMMENT 'ID (UUID v7)',
    `name` varchar(100) NOT NULL COMMENT '数据源名称',
    `type` varchar(50) NOT NULL COMMENT '数据源类型：mysql/postgresql/redis/mongodb/sqlserver',
    `host` varchar(200) NOT NULL COMMENT '连接地址',
    `port` int NOT NULL COMMENT '连接端口',
    `database` varchar(100) DEFAULT NULL COMMENT '数据库名',
    `username` varchar(100) NOT NULL COMMENT '连接用户名',
    `password` varchar(500) NOT NULL COMMENT '连接密码（AES-256-GCM 加密）',
    `description` varchar(500) DEFAULT NULL COMMENT '描述信息',
    `status` varchar(20) NOT NULL DEFAULT 'enabled' COMMENT '状态：enabled/disabled',
    `sort_order` int NOT NULL DEFAULT 0 COMMENT '排序权重',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_name` (`name`),
    KEY `idx_status` (`status`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='数据源配置表';
```

### Go Struct

```go
type DataSource struct {
    Id          string    `gorm:"primaryKey;size:36"`      // UUID v7
    Name        string    `gorm:"size:100;not null;uniqueIndex"`
    Type        string    `gorm:"size:50;not null"`
    Host        string    `gorm:"size:200;not null"`
    Port        int       `gorm:"not null"`
    Database    string    `gorm:"size:100"`
    Username    string    `gorm:"size:100;not null"`
    Password    string    `gorm:"size:500;not null"`       // 加密存储
    Description string    `gorm:"size:500"`
    Status      string    `gorm:"size:20;not null;default:'enabled'"`
    SortOrder   int       `gorm:"not null;default:0"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
    DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// 响应时隐藏敏感信息
type DataSourceResp struct {
    Id          string    `json:"id"`
    Name        string    `json:"name"`
    Type        string    `json:"type"`
    Host        string    `json:"host"`
    Port        int       `json:"port"`
    Database    string    `json:"database"`
    Description string    `json:"description"`
    Status      string    `json:"status"`
    SortOrder   int       `json:"sort_order"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

---

## API Contract

根据 Spec 中的 Acceptance Criteria 生成 API 定义：

**位置**: `api/doc/datasource/datasource.api`

```api
syntax = "v1"

import "../../base.api"

type (
    // 创建数据源请求
    CreateDataSourceReq {
        Name        string `json:"name" validate:"required,max=100"`                           // 数据源名称
        Type        string `json:"type" validate:"required,oneof=mysql postgresql redis mongodb sqlserver"` // 数据源类型
        Host        string `json:"host" validate:"required,max=200"`                           // 连接地址
        Port        int    `json:"port" validate:"required,min=1,max=65535"`                   // 连接端口
        Database    string `json:"database" validate:"max=100"`                                // 数据库名（部分类型必填）
        Username    string `json:"username" validate:"required,max=100"`                       // 连接用户名
        Password    string `json:"password" validate:"required,max=200"`                       // 连接密码
        Description string `json:"description" validate:"max=500"`                            // 描述信息
        Status      string `json:"status" validate:"oneof=enabled disabled"`                   // 状态
        SortOrder   int    `json:"sort_order" validate:"omitempty"`                           // 排序权重
    }

    // 创建数据源响应
    CreateDataSourceResp {
        Id        string    `json:"id"`
        CreatedAt time.Time `json:"created_at"`
    }

    // 列表查询请求（继承分页和搜索）
    ListDataSourceReq {
        PageInfoWithKeyword
        Status   string `form:"status,optional"`    // 状态筛选：enabled/disabled
    }

    // 列表查询响应
    ListDataSourceResp {
        Entries    []DataSourceResp `json:"entries"`
        TotalCount int64           `json:"total_count"`
    }
)

@server(
    prefix: /api/v1/datasource
    group: datasource
)
service spec-cc-0104-api {
    @handler CreateDataSource
    post / (CreateDataSourceReq) returns (CreateDataSourceResp)

    @handler ListDataSource
    get / (ListDataSourceReq) returns (ListDataSourceResp)
}
```

---

## Testing Strategy

| 类型 | 方法 | 覆盖率 |
|------|------|--------|
| 单元测试 | 表驱动测试，Mock Model | > 80% |
| 集成测试 | 测试数据库（MySQL） | 核心流程 |
| 连接测试 | 模拟连接（不连接真实数据库） | 所有支持类型 |

**测试用例**:
- 创建数据源成功/失败场景
- 列表查询（分页、搜索、排序、筛选）
- 参数校验（必填、范围、枚举）
- 名称唯一性校验
- 连接测试（模拟成功/失败）
- 软删除验证
- 并发创建测试

---

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2026-01-04 | - | 初始版本，定义数据源管理技术方案 |
