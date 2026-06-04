<p align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go" alt="Go">
  <img src="https://img.shields.io/badge/Vue-3.x-4FC08D?logo=vue.js" alt="Vue">
  <img src="https://img.shields.io/badge/MySQL-8.0-4479A1?logo=mysql" alt="MySQL">
  <img src="https://img.shields.io/badge/Redis-7-DC382D?logo=redis" alt="Redis">
  <img src="https://img.shields.io/badge/Docker-✔-2496ED?logo=docker" alt="Docker">
</p>

<h1 align="center">ProjectFlow</h1>
<p align="center">轻量级软件研发协作平台</p>

<p align="center">
  <a href="#特性概览">特性概览</a> •
  <a href="#技术栈">技术栈</a> •
  <a href="#快速开始">快速开始</a> •
  <a href="#项目结构">项目结构</a> •
  <a href="#api-文档">API 文档</a>
</p>

---

## 项目简介

ProjectFlow 是一套轻量级研发协作平台，解决传统微信群、Excel 管理项目的痛点——项目状态不透明、任务分配混乱、进度跟踪困难、文件版本丢失。支持从项目立项到归档的全生命周期管理，内置 Kanban 看板、评论协同、文件附件、操作审计等功能。

**适用场景**：小型研发团队（5-50 人）、敏捷开发 Sprint 管理、外包项目交付跟踪。

---

## 特性概览

🎯 **项目管理**
- 6 状态工作流：Draft → PendingApproval → Approved → Developing → Completed → Archived
- 项目经理审批、成员管理、操作日志审计

📋 **任务管理**
- 5 状态 Kanban 看板：Todo → InProgress → Review → Testing → Done
- 乐观锁并发控制——多人同时编辑时防止覆盖
- 按优先级/负责人筛选，5s 自动刷新

🔐 **用户与权限**
- RBAC 三级角色：Guest / Developer / ProjectManager
- JWT + Refresh Token 双令牌认证
- 图形验证码、登录失败限流

💬 **评论协同**
- 项目级 / 任务级评论，嵌套回复
- `@username` 自动解析并发送站内通知

📎 **附件管理**
- 图片 / 文档 / 压缩包上传，单文件 ≤20MB
- 文件类型校验，UUID 命名防冲突

📊 **看板与缓存**
- 5 列 Kanban 布局，拖拽更新状态
- Redis Cache-Aside 模式，故障自动回源 MySQL

---

## 技术栈

| 领域 | 技术选型 | 说明 |
|---|---|---|
| 后端框架 | Go 1.21 + Gin 1.9 | 高性能 HTTP 框架 |
| ORM | GORM 1.25 | 自动迁移 + 关联预加载 |
| 数据库 | MySQL 8.0 | 8 张表，乐观锁版本控制 |
| 缓存 | Redis 7 | Cache-Aside，10min TTL |
| 对象存储 | MinIO | 文件存储（可切换本地） |
| 认证 | JWT (golang-jwt/v5) | HMAC-SHA256 签名 |
| 前端框架 | Vue 3 + TypeScript | Composition API |
| 构建工具 | Vite 5 | HMR 热更新 |
| 状态管理 | Pinia | 响应式 store |
| UI 风格 | Material Design | 自写 CSS，零依赖组件库 |
| 部署 | Docker Compose | MySQL + Redis + MinIO + Nginx + App |

---

## 快速开始

### 前置要求

- Go 1.21+
- Node.js 18+
- Docker & Docker Compose

### 一键启动

```bash
# 1. 启动基础设施
docker compose up -d mysql redis minio

# 2. 启动后端 (端口 8080)
cd backend
cp config/config.yaml.example config/config.yaml  # 按需修改数据库密码
go run ./cmd/server

# 3. 启动前端 (端口 3000)
cd frontend
npm install
npm run dev
```

启动后访问 `http://localhost:3000`，注册账号即可使用。

### 服务端口

| 服务 | 地址 | 默认凭据 |
|---|---|---|
| 后端 API | http://localhost:8080 | — |
| 前端 | http://localhost:3000 | — |
| MySQL | localhost:3306 | root / root |
| Redis | localhost:6379 | — |
| MinIO Console | http://localhost:9001 | minioadmin / minioadmin123 |

### Docker 全量部署

```bash
docker compose up -d          # 启动全部服务
docker compose logs -f app    # 查看后端日志
docker compose down -v        # 停止并清理数据卷
```

---

## 项目结构

```
cowork/
├── backend/                     # Go 后端
│   ├── cmd/server/main.go       # 入口
│   ├── internal/
│   │   ├── config/              # Viper 配置加载
│   │   ├── model/               # GORM 数据模型 (7 实体)
│   │   ├── handler/             # HTTP 处理层 (40+ 端点)
│   │   ├── service/             # 业务逻辑层
│   │   ├── repository/          # 数据访问层
│   │   ├── middleware/           # JWT/RBAC/CORS/限流/日志
│   │   ├── router/              # 路由注册
│   │   ├── cache/               # Redis 缓存层
│   │   ├── db/                  # DB + Redis 连接
│   │   └── dto/                 # 请求/响应 DTO
│   ├── pkg/
│   │   ├── jwt/                 # JWT 工具
│   │   ├── captcha/             # 图形验证码
│   │   └── errcode/             # 错误码
│   ├── config/config.yaml       # 配置文件
│   ├── Dockerfile
│   └── Makefile
├── frontend/                    # Vue 3 前端
│   ├── src/
│   │   ├── views/               # 6 页面
│   │   ├── components/          # 8 通用组件
│   │   ├── api/                 # Axios 封装
│   │   ├── stores/              # Pinia 状态管理
│   │   └── router/              # 路由守卫
│   ├── vite.config.ts
│   └── package.json
├── nginx/nginx.conf
├── docker-compose.yml
└── README.md
```

---

## API 文档

### 用户认证

```bash
# 注册
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"demo","email":"demo@test.com","password":"Demo1234"}'

# 获取验证码
curl http://localhost:8080/api/auth/captcha
# → {"data":{"captcha_id":"xxx","captcha_image":"iVBORw0KGgo..."}}

# 登录 (从 Redis 获取验证码答案后)
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"demo","password":"Demo1234","captcha_id":"xxx","captcha_answer":"42"}'
# → {"data":{"access_token":"eyJ...","refresh_token":"eyJ...","user":{...}}}
```

### 项目与任务

```bash
TOKEN="<access_token>"

# 创建项目 (需要 ProjectManager 角色)
curl -X POST http://localhost:8080/api/projects \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"协作平台开发","description":"v1.0"}'

# 创建任务
curl -X POST http://localhost:8080/api/projects/1/tasks \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"用户登录模块","priority":"High","project_id":1}'

# 获取看板
curl http://localhost:8080/api/projects/1/kanban \
  -H "Authorization: Bearer $TOKEN"

# 流转任务状态
curl -X PUT http://localhost:8080/api/tasks/1/status \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status":"InProgress"}'
```

完整 API 端点清单见 [version_1.0.md](./version_1.0.md)。

---

## 设计亮点

**乐观锁并发控制** — Task 和 Project 维护 `version` 字段，`UPDATE WHERE id=? AND version=?` 防止多人编辑覆盖，冲突时提示刷新重试。

**Cache-Aside 容错** — Redis 不可用时自动回源 MySQL，零宕机降级；`cache.Get/Set/Del` 统一封装，Redis nil 时静默跳过。

**前端响应拦截器** — axios 拦截器统一解包 `{code, message, data}` 后端信封，非零 code 自动转为异常，无需每个调用手动判断。

**Pinia 登录态管理** — `isLoggedIn` 使用 `ref` 显式更新而非 `computed`，避免 `localStorage` 作为非响应式依赖导致缓存陷阱。

**GORM 乐观锁修复** — GORM `Updates(map/struct)` 存在字段过滤，项目改用原生 SQL `Exec("UPDATE ... WHERE ...")` 确保关键字段可靠写入。

---

## 贡献与许可

欢迎提交 PR 或 Issue。

本项目仅用于学习和团队内部协作场景。
