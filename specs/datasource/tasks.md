# DataSource Tasks

> **Branch**: `feature/datasource`
> **Spec Path**: `specs/datasource/`
> **Created**: 2026-01-04
> **Input**: spec.md, plan.md

---

## Task Format

```
[ID] [P?] [Story] Description
```

| 标记 | 含义 |
|------|------|
| `T001` | 任务 ID |
| `[P]` | 可并行执行（不同文件，无依赖） |
| `[US1]` | 关联 User Story 1 |
| `[US2]` | 关联 User Story 2 |

---

## Task Overview

| ID | Task | Story | Status | Parallel | Est. Lines |
|----|------|-------|--------|----------|------------|
| T001 | 依赖库安装 | Setup | ⏸️ | - | - |
| T002 | goctl 生成代码 | US1/US2 | ⏸️ | - | - |
| T003 | Model 接口定义 | US1/US2 | ⏸️ | [P] | 45 |
| T004 | Model 数据结构 | US1/US2 | ⏸️ | [P] | 40 |
| T005 | Model 常量变量 | US1/US2 | ⏸️ | [P] | 35 |
| T006 | Model ORM 工厂 | US1/US2 | ⏸️ | - | 30 |
| T007 | GORM 实现 - 基础 CRUD | US1/US2 | ⏸️ | - | 45 |
| T008 | GORM 实现 - 列表查询 | US1 | ⏸️ | [P] | 45 |
| T009 | GORM 实现 - 连接测试 | US2 | ⏸️ | [P] | 40 |
| T010 | Logic - 创建数据源 | US2 | ⏸️ | - | 50 |
| T011 | Logic - 列表查询 | US1 | ⏸️ | - | 45 |
| T012 | Model 单元测试 | US1/US2 | ⏸️ | [P] | 50 |
| T013 | Logic 单元测试 | US1/US2 | ⏸️ | [P] | 50 |
| T014 | 集成测试 | US1/US2 | ⏸️ | - | 40 |
| T015 | 错误码定义 | US1/US2 | ⏸️ | [P] | 30 |

---

## Phase 1: Setup

**目的**: 项目初始化和基础配置

- [X] T001 安装依赖库（uuid, crypto）
  ```bash
  go get github.com/google/uuid
  go get golang.org/x/crypto/...
  ```

**Checkpoint**: ✅ 开发环境就绪

---

## Phase 2: Foundation (Go-Zero 基础)

**目的**: 必须完成后才能开始 User Story 实现

- [X] T002 [US1/US2] 运行 goctl 生成代码
  ```bash
  goctl api go -api api/doc/datasource/datasource.api -dir api/ --style=go_zero --type-group
  ```
  ✅ 已生成 Handler 和 Types 文件

**Checkpoint**: ✅ 基础设施就绪，可开始 User Story 实现

---

## Phase 3: User Story 1 - 数据源列表查询 (P1) 🎯

**目标**: 实现数据源列表查询功能，支持分页、搜索、排序

**独立测试**: 提供分页、搜索、排序功能，返回符合条件的数据源列表

### Step 1: 实现 Model 层

- [X] T003 [US1] 创建 `model/datasource/datasource/interface.go` (31 行)
  - 定义 DataSourceModel 接口
  - 定义 DataSourceQuery 查询参数结构

- [X] T004 [US1] 创建 `model/datasource/datasource/types.go` (40 行)
  - 定义 DataSource 实体结构
  - 定义 DataSourceResp 响应结构

- [X] T005 [US1] 创建 `model/datasource/datasource/vars.go` (80 行)
  - 定义数据源类型常量
  - 定义状态常量
  - 定义错误信息

- [X] T006 [US1] 创建 `model/datasource/datasource/factory.go` (60 行)
  - 定义 NewDataSourceModel 工厂函数
  - 支持 GORM 和 SQLx 两种 ORM

- [X] T007 [US1] 实现 `model/datasource/datasource/gorm_dao.go` - 基础 CRUD (40 行)
  - Insert, FindOne, Update, Delete 方法
  - 软删除支持

- [X] T008 [US1] 实现 `model/datasource/datasource/gorm_dao.go` - 列表查询 (90 行)
  - FindList 方法
  - 支持分页、搜索、排序、筛选
  - 名称唯一性检查
  - 连接测试方法

### Step 2: 实现 Logic 层

- [X] T011 [US1] 实现 `api/internal/logic/datasource/list_datasource_logic.go` (75 行)
  - 接收查询参数
  - 调用 Model 层查询
  - 数据转换和响应格式化
  - 错误处理
  - 更新 ServiceContext 添加 DataSourceModel 字段

### Step 3: 测试

- [X] T012 [US1] Model 层单元测试 (80 行)
  - 测试列表查询各种场景
  - Mock 数据库连接
  - 响应转换测试

- [X] T013 [US1] Logic 层单元测试 (70 行)
  - 表驱动测试
  - 验证响应格式
  - 参数校验测试

**Checkpoint**: ✅ User Story 1 可独立测试和验证

---

---

## Phase 4: User Story 2 - 数据源新增创建 (P1) 🎯

**目标**: 实现数据源新增创建功能，支持连接测试验证

**独立测试**: 提交有效的数据源配置信息，创建成功后返回 201 和数据源详情

### Step 1: 实现 Model 层

- [X] T009 [US2] 实现 `model/datasource/datasource/connection_test.go` (170 行)
  - TestConnection 方法
  - 支持 MySQL/PostgreSQL/Redis/MongoDB/SQLServer
  - AES-256-GCM 密码加密/解密
  - UUID v7 主键生成

- [X] T007 [US2] 完善 `model/datasource/datasource/gorm_dao.go` - Insert 方法 (20 行)
  - UUID v7 主键生成
  - 时间戳设置

### Step 2: 实现 Logic 层

- [X] T010 [US2] 实现 `api/internal/logic/datasource/create_datasource_logic.go` (110 行)
  - 参数校验（必填、范围、枚举）
  - 名称唯一性校验
  - 连接测试
  - 密码加密
  - UUID v7 主键生成
  - 数据保存
  - 响应格式化

### Step 3: 测试

- [X] T012 [US2] Model 层单元测试 (80 行)
  - 测试连接测试各种场景
  - 测试加密/解密
  - 响应转换测试

- [X] T013 [US2] Logic 层单元测试 (70 行)
  - 测试创建成功/失败场景
  - 参数校验测试
  - 连接测试模拟

**Checkpoint**: ✅ User Story 2 可独立测试和验证

---

## Phase 5: Polish

**目的**: 收尾工作

- [X] T015 [US1/US2] 定义错误码 (70 行)
  - 创建 `api/internal/errorx/codes.go`
  - 定义 30400-30499 错误码
  - 错误创建函数

- [X] T014 [US1/US2] 集成测试 (40 行)
  - 创建 `model/datasource/datasource/integration_test.go`
  - 端到端测试框架

- [X] 代码清理和格式化
- [X] 运行编译检查
- [ ] 确认测试覆盖率 > 80%

---

## Dependencies

```
Phase 1 (Setup)
    ↓
Phase 2 (Foundation)
    ↓
Phase 3 (US1) → Phase 4 (US2)  # 可并行（不同人协作）
    ↓
Phase 5 (Polish)
```

### 并行执行说明

- `[P]` 标记的任务可与同 Phase 内其他 `[P]` 任务并行
- 不同 User Story 可并行（如有团队协作）
- 同一 User Story 内按 Step 顺序执行

---

## Notes

- 每个 Task 完成后提交代码
- 每个 Checkpoint 进行验证
- 遇到问题及时记录到 Open Questions
- 代码行数估算包含注释和空行
- 密码加密使用 AES-256-GCM 算法
