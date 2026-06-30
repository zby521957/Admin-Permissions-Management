# RBAC Admin

基于 **Go + Gin + GORM + MySQL + Redis** 构建的 RBAC（基于角色的访问控制）权限管理系统，提供用户、角色、权限的完整 CRUD 操作与细粒度权限校验。

## 功能特性

- **用户管理** — 注册、登录、信息更新、软删除、角色分配
- **角色管理** — 角色创建、更新、删除、权限分配
- **权限管理** — 细粒度资源级权限（如 \user:list\、\ole:create\），支持增删改查
- **JWT 认证** — 基于 HS256 签名，24 小时有效期
- **权限缓存** — Redis 缓存用户权限，减少数据库查询；角色/权限变更时自动失效
- **Docker 一键部署** — 配套 Docker Compose，MySQL + Redis 开箱即用
- **分环境配置** — 通过 \ENV\ 环境变量切换本地 / Docker 配置

## 技术栈

| 技术 | 用途 |
|---|---|
| Go 1.25 | 后端语言 |
| Gin | HTTP 框架 |
| GORM | ORM，支持 AutoMigrate |
| MySQL 8.0 | 关系型数据库 |
| Redis 7 | 权限缓存 |
| JWT | Token 认证 |
| bcrypt | 密码哈希 |
| Viper | 配置管理 |
| Docker | 容器化部署 |

## 架构

\\\mermaid
graph TD
    Client[客户端] -->|HTTP JSON| Router[Gin Router]
    Router -->|公开路由| AuthCtrl[controller/auth.go]
    Router -->|受保护路由| JWTMiddle[中间件: JWT 认证]
    JWTMiddle -->|注入用户信息| PermMiddle[中间件: 权限校验]
    PermMiddle -->|校验权限| Ctrl[controller/*.go]
    Ctrl -->|调用| Svc[service/*.go]
    Svc -->|查询| DB[(MySQL)]
    Svc -->|缓存读写| Cache[(Redis)]
    DB -->|自动迁移| Model[model/*.go]
\\\

## 快速开始

### 前置条件

- Go 1.25+
- MySQL 8.0（运行在 127.0.0.1:3306）
- Redis 7（运行在 127.0.0.1:6379）

### 本地运行

\\\ash
# 1. 克隆项目
git clone <your-repo-url>
cd rbac-admin

# 2. 编辑配置文件 config.yaml，确认数据库连接信息

# 3. 启动服务
go run .

# 4. 验证 - 注册用户
curl -X POST http://localhost:8080/api/v1/register \\
  -H \"Content-Type: application/json\" \\
  -d '{\"username\":\"admin\",\"password\":\"admin123\",\"email\":\"admin@test.com\"}'

# 5. 验证 - 登录获取 Token
curl -X POST http://localhost:8080/api/v1/login \\
  -H \"Content-Type: application/json\" \\
  -d '{\"username\":\"admin\",\"password\":\"admin123\"}'
\\\

### Docker 部署

\\\ash
docker compose up -d
\\\

服务运行在 http://localhost:8080。

> Docker 环境使用独立的 MySQL 与 Redis 容器，数据与本地隔离，需重新注册用户。

## API 接口

| 方法 | 路径 | 权限 | 说明 |
|---|---|---|---|
| POST | /api/v1/register | 公开 | 用户注册 |
| POST | /api/v1/login | 公开 | 用户登录 |
| GET | /api/v1/users | user:list | 用户列表 |
| GET | /api/v1/users/:id | user:list | 用户详情 |
| PUT | /api/v1/users/:id | user:update | 更新用户 |
| DELETE | /api/v1/users/:id | user:delete | 删除用户 |
| POST | /api/v1/users/:id/roles | user:update | 分配角色 |
| GET | /api/v1/roles | role:list | 角色列表 |
| GET | /api/v1/roles/:id | role:list | 角色详情 |
| POST | /api/v1/roles | role:create | 创建角色 |
| PUT | /api/v1/roles/:id | role:update | 更新角色 |
| DELETE | /api/v1/roles/:id | role:delete | 删除角色 |
| POST | /api/v1/roles/:id/permissions | role:update | 分配权限 |
| GET | /api/v1/permissions | permission:list | 权限列表 |
| GET | /api/v1/permissions/:id | permission:list | 权限详情 |
| POST | /api/v1/permissions | permission:create | 创建权限 |
| PUT | /api/v1/permissions/:id | permission:update | 更新权限 |
| DELETE | /api/v1/permissions/:id | permission:delete | 删除权限 |

## 项目结构

\\\
rbac-admin/
├── cache/            # Redis 缓存（权限缓存、缓存失效）
├── config/           # 配置加载（Viper + 多环境支持）
├── controller/       # HTTP 处理器（参数校验 + 调用 service）
├── middleware/       # Gin 中间件（JWT 认证 + 权限校验）
├── model/            # GORM 模型（User / Role / Permission，含钩子）
├── router/           # 路由注册（公开路由 + 受保护路由）
├── service/          # 业务逻辑层（组合模型操作 + 缓存策略）
├── utils/            # 工具库（JWT 生成/解析、密码哈希、统一响应）
├── config.yaml       # 本地开发配置
├── config.docker.yaml # Docker 部署配置
├── docker-compose.yaml # Docker Compose 编排
├── Dockerfile        # 多阶段构建
├── go.mod
├── main.go           # 程序入口
└── README.md
\\\

## 开发命令

| 命令 | 说明 |
|---|---|
| \go run .\ | 启动服务 |
| \go build -o server .\ | 编译 |
| \docker compose up -d\ | Docker 部署 |
| \docker compose down\ | 停止容器（保留数据） |
| \docker compose down -v\ | 停止容器并删除数据 |
