# DataSource Specification

> **Branch**: `feature/datasource`
> **Spec Path**: `specs/datasource/`
> **Created**: 2026-01-04

> **Status**: Draft

---

## Overview

实现数据源管理功能，支持多种数据源的统一管理，包括数据源列表查询和新增创建功能。用户可以通过可视化的方式管理不同类型的数据源（如 MySQL、PostgreSQL、Redis、MongoDB 等），为后续的数据集成和分析提供基础。

---

## User Stories

### Story 1: 数据源列表查询 (P1)

AS a 数据管理员
I WANT 查看系统中所有数据源的列表信息
SO THAT 可以快速了解已配置的数据源状态和基本信息

**独立测试**: 提供分页、搜索、排序功能，返回符合条件的数据源列表

### Story 2: 数据源新增创建 (P1)

AS a 数据管理员
I WANT 新增一个数据源配置
SO THAT 可以连接并使用新的数据源进行数据操作

**独立测试**: 提交有效的数据源配置信息，创建成功后返回 201 和数据源详情

---

## Acceptance Criteria (EARS)

### 正常流程

| ID | Scenario | Trigger | Expected Behavior |
|----|----------|---------|-------------------|
| AC-01 | 创建数据源成功 | WHEN 用户提交有效的数据源配置 | THE SYSTEM SHALL 验证配置、生成 UUID v7 主键、保存到数据库并返回 201 |
| AC-02 | 列表查询成功 | WHEN 用户访问数据源列表接口 | THE SYSTEM SHALL 返回分页的数据源列表（包含总数、每页数据） |
| AC-03 | 按关键字搜索 | WHEN 用户输入关键字进行搜索 | THE SYSTEM SHALL 模糊匹配名称或描述，返回匹配结果 |
| AC-04 | 按状态筛选 | WHEN 用户按启用/禁用状态筛选 | THE SYSTEM SHALL 只返回指定状态的数据源 |

### 异常处理

| ID | Scenario | Trigger | Expected Behavior |
|----|----------|---------|-------------------|
| AC-10 | 必填参数为空 | WHEN 必填参数（名称、类型、连接信息）为空 | THE SYSTEM SHALL 返回 400，提示"必填参数不能为空" |
| AC-11 | 数据源名称重复 | WHEN 创建时名称与已有数据源重复 | THE SYSTEM SHALL 返回 409，提示"数据源名称已存在" |
| AC-12 | 无效的数据源类型 | WHEN 类型字段不在支持范围内 | THE SYSTEM SHALL 返回 400，提示"不支持的数据源类型" |
| AC-13 | 连接测试失败 | WHEN 新增/编辑时连接测试失败 | THE SYSTEM SHALL 返回 400，提示"数据源连接失败" |
| AC-14 | 偏移量超范围 | WHEN 分页请求的 offset 超出实际范围 | THE SYSTEM SHALL 返回空列表和总数 |

---

## Edge Cases

| ID | Case | Expected Behavior |
|----|------|-------------------|
| EC-01 | 大量数据分页查询 | 支持分页，最大 limit=2000，返回查询时间 <200ms (P99) |
| EC-02 | 同时创建同名数据源 | 仅第一个成功，其他返回 409，避免重复创建 |
| EC-03 | 特殊字符处理 | 名称和描述支持中文、英文、数字、特殊符号，长度限制 1-100 字符 |
| EC-04 | 空列表查询 | 返回空数组和 total_count=0 |

---

## Business Rules

| ID | Rule | Description |
|----|------|-------------|
| BR-01 | 名称唯一性 | 系统内数据源名称不能重复（Single-tenant 模式） |
| BR-02 | 连接信息加密 | 数据库密码等敏感信息需要使用 AES-256-GCM 加密存储 |
| BR-03 | 支持类型 | 必须支持 MySQL、PostgreSQL、Redis、MongoDB、SQLServer |
| BR-04 | 状态管理 | 数据源有启用/禁用两种状态，禁用状态下不可用于数据操作 |
| BR-05 | 主键规范 | 所有数据源记录使用 UUID v7 作为主键 |
| BR-06 | 软删除 | 删除数据源采用软删除，保留历史记录 |
| BR-07 | 租户模式 | 采用 Single-tenant 模式，所有用户共享数据源 |
| BR-08 | 权限控制 | 采用单一管理员权限模型，适合初期简化管理 |

---

## Data Considerations

需要持久化的数据字段：

| Field | Description | Constraints |
|-------|-------------|-------------|
| id | 数据源唯一标识 | CHAR(36), UUID v7, 主键 |
| name | 数据源名称 | 必填，1-100 字符，唯一 |
| type | 数据源类型 | 必填，枚举值：mysql/postgresql/redis/mongodb/sqlserver |
| host | 连接地址 | 必填，IP 或域名，1-200 字符 |
| port | 连接端口 | 必填，1-65535 整数 |
| database | 数据库名/库名 | 可选，0-100 字符（部分类型必填） |
| username | 连接用户名 | 必填，1-100 字符 |
| password | 连接密码 | 必填，加密存储，1-200 字符 |
| description | 描述信息 | 可选，0-500 字符 |
| status | 状态 | 枚举：enabled/disabled，默认 enabled |
| sort_order | 排序权重 | 整数，默认 0，数值越大排序越靠前 |
| created_at | 创建时间 | 时间戳，自动生成 |
| updated_at | 更新时间 | 时间戳，自动更新 |
| deleted_at | 删除时间 | 时间戳，NULL 表示未删除 |

---

## Success Metrics

| ID | Metric | Target |
|----|--------|--------|
| SC-01 | 列表接口响应时间 | < 200ms (P99) |
| SC-02 | 创建接口响应时间 | < 500ms (P99) |
| SC-03 | 测试覆盖率 | > 80% |
| SC-04 | 连接测试成功率 | > 95% |

---

## Clarifications

### Session 2026-01-04

- Q: 密码等敏感信息应该使用什么加密算法和方式？ → A: AES-256-GCM 加密存储
- Q: 是否需要支持多租户隔离？ → A: Single-tenant（单租户模式）
- Q: 数据源权限控制模型如何设计？ → A: 单一管理员权限模型
- Q: 连接测试的具体触发时机是什么？ → A: 仅新增/编辑时测试
- Q: 用户角色如何定义和划分？ → A: 仅有管理员角色

---

## Open Questions

*（无待决问题）*

---

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2026-01-04 | - | 初始版本，定义数据源列表和新增功能 |
