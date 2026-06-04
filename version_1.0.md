# cowork v1.0 —— 软件研发协作平台 版本记录

> 生成日期：2026-06-04  
> 技术栈：Go 1.21+ / Gin / GORM / MySQL / Redis / MinIO / Vue 3 / Vite  
> 架构：分层架构 (Handler → Service → Repository)，Cache-Aside 缓存，乐观锁并发控制

---

## 1. 项目概述

cowork (ProjectFlow) 是一套轻量级软件研发协作平台，实现项目全生命周期管理。支持项目经理规划审批项目、分配任务；开发者完成任务开发、提交成果、协同沟通。

### 核心能力

- **项目管理**：6 状态工作流（Draft → PendingApproval → Approved → Developing → Completed → Archived）
- **任务管理**：5 状态 Kanban 看板（Todo → InProgress → Review → Testing → Done），乐观锁防冲突
- **用户体系**：RBAC 三级权限（Guest / Developer / ProjectManager），JWT + Refresh Token 认证
- **协同安全**：乐观锁版本控制，图形验证码，登录限流
- **缓存策略**：Redis Cache-Aside 模式，故障自动回源数据库
- **文件存储**：MinIO 对象存储，单文件 ≤20MB
- **审计追溯**：全操作日志记录（操作人/IP/内容/时间）

---

## 2. 项目结构

```
cowork/
├── backend/                        # Go 后端
│   ├── cmd/server/main.go          # 入口：配置加载 → DB/Redis 初始化 → 自动迁移 → Gin 启动
│   ├── internal/
│   │   ├── config/config.go        # Viper 配置加载 (Server/DB/Redis/MinIO/JWT)
│   │   ├── model/                  # GORM 数据模型 (7 个实体)
│   │   │   ├── user.go             # 用户 (ID, Username, Email, PasswordHash, Role)
│   │   │   ├── project.go          # 项目 + ProjectMember (Status, Owner, Version)
│   │   │   ├── task.go             # 任务 (Priority, Status, Assignee, Version)
│   │   │   ├── comment.go          # 评论 (自引用 ParentID 支持嵌套回复)
│   │   │   ├── attachment.go       # 附件 (FileURL, FileSize, ContentType)
│   │   │   ├── notification.go     # 站内通知 (Type, IsRead, RelatedID)
│   │   │   └── operation_log.go    # 操作日志 (UserID, Action, IP)
│   │   ├── handler/                # HTTP 处理层 (参数校验 → 调用 Service)
│   │   │   ├── auth_handler.go     # 注册/登录/刷新/个人资料/验证码
│   │   │   ├── project_handler.go  # 项目 CRUD / 审批 / 成员管理 (12 端点)
│   │   │   ├── task_handler.go     # 任务 CRUD / 指派 / 状态流转 (7 端点)
│   │   │   ├── kanban_handler.go   # Kanban 看板 (5 列分组 + 筛选)
│   │   │   ├── comment_handler.go  # 评论 / 回复 / 删除 (6 端点)
│   │   │   ├── attachment_handler.go # 上传 / 列表 / 删除 (4 端点)
│   │   │   ├── notification_handler.go # 通知列表 / 未读数 / 已读
│   │   │   └── log_handler.go      # 操作日志查询 (PM only)
│   │   ├── service/                # 业务逻辑层
│   │   │   ├── auth_service.go     # 注册 (密码复杂度/唯一性) / 登录 (验证码 + bcrypt)
│   │   │   ├── project_service.go  # 状态流转校验 + Cache-Aside
│   │   │   ├── task_service.go     # 乐观锁 + 状态流转 + Cache-Aside
│   │   │   ├── comment_service.go  # @mention 解析 + 自动通知
│   │   │   ├── attachment_service.go # 文件大小/类型校验 + 本地存储
│   │   │   ├── notification_service.go
│   │   │   └── log_service.go      # 异步记录
│   │   ├── repository/             # 数据访问层
│   │   │   ├── user_repo.go / project_repo.go / task_repo.go
│   │   │   ├── comment_repo.go / notification_repo.go / log_repo.go
│   │   ├── middleware/             # 中间件
│   │   │   ├── auth.go             # JWT Bearer Token 验证
│   │   │   ├── rbac.go             # RequireRole("ProjectManager")
│   │   │   ├── cors.go             # 跨域
│   │   │   ├── ratelimit.go        # IP 限流 (Redis)
│   │   │   └── logger.go           # 请求日志
│   │   ├── router/router.go        # 路由注册 (40+ API 端点)
│   │   ├── cache/redis.go          # 统一缓存接口 (Get/Set/Del/DeletePattern)
│   │   ├── db/db.go                # DB + Redis 全局连接
│   │   └── dto/
│   │       ├── request/            # 请求 DTO (Gin binding tags)
│   │       └── response/           # 统一响应 {code, message, data}
│   ├── pkg/
│   │   ├── jwt/jwt.go              # Access Token + Refresh Token (HMAC-SHA256)
│   │   ├── captcha/captcha.go      # 数学验证码 (PNG, base64, Redis 5min TTL)
│   │   └── errcode/errcode.go      # 错误码常量 (7 个范围: 9xxxx ~ 6xxxx)
│   ├── config/config.yaml          # 环境配置
│   ├── Dockerfile                  # 多阶段构建 (golang:1.21-alpine → alpine:3.19)
│   └── Makefile
├── frontend/                       # Vue 3 前端
│   ├── src/
│   │   ├── views/                  # 6 个页面
│   │   │   ├── Login.vue           # 登录 (验证码图片)
│   │   │   ├── Register.vue        # 注册
│   │   │   ├── Dashboard.vue       # 主页 (Header + Sidebar + Kanban)
│   │   │   ├── ProjectList.vue     # 项目卡片列表
│   │   │   ├── ProjectDetail.vue   # 项目详情 + 成员管理
│   │   │   └── Settings.vue        # 个人设置
│   │   ├── components/             # 8 个通用组件
│   │   │   ├── AppHeader.vue       # 顶部导航栏 (Logo + 搜索 + 通知铃铛 + 头像)
│   │   │   ├── AppSidebar.vue      # 左侧项目列表
│   │   │   ├── KanbanBoard.vue     # 5 列看板 (5s 轮询刷新)
│   │   │   ├── KanbanColumn.vue    # 单列 + 任务计数
│   │   │   ├── KanbanCard.vue      # 任务卡片 (优先级标签 + 截止日期)
│   │   │   ├── TaskDetailModal.vue # 任务详情弹窗 (编辑/状态流转/评论)
│   │   │   ├── CommentSection.vue  # 评论区 (嵌套回复)
│   │   │   └── NotificationBell.vue # 通知铃铛 + 未读 badge
│   │   ├── api/                    # 6 个 API 模块 (Axios + 拦截器)
│   │   │   ├── index.ts            # 请求/响应拦截器 (Token 附加 + 信封解包 + 错误码异常)
│   │   │   ├── auth.ts / project.ts / task.ts / comment.ts / attachment.ts / notification.ts
│   │   ├── stores/                 # 4 个 Pinia Store
│   │   │   ├── auth.ts             # 登录态 (ref, 非 computed 避免缓存陷阱)
│   │   │   ├── project.ts / task.ts / notification.ts
│   │   ├── router/index.ts         # 路由守卫 (requiresAuth + guest)
│   │   └── style.css               # Material Design 全局样式
│   ├── vite.config.ts              # Vite 配置 (API 代理到 :8080)
│   ├── package.json
│   └── tsconfig.json
├── nginx/nginx.conf                # Nginx 反向代理
├── docker-compose.yml              # MySQL 8.0 + Redis 7 + MinIO + App + Nginx
├── .gitignore
└── README.md
```

---

## 3. 数据库设计

### 表结构 (8 张表, GORM AutoMigrate)

| 表名 | 关键字段 | 索引 |
|---|---|---|
| `users` | id, username(U), email(U), password_hash, role | username, email |
| `projects` | id, name, description, status, owner_id(FK), version | owner_id |
| `project_members` | id, project_id(FK), user_id(FK), role | project_id, user_id |
| `tasks` | id, title, priority, status, project_id(FK), assignee_id(FK), creator_id(FK), deadline, version | project_id, assignee_id, creator_id |
| `comments` | id, content, user_id(FK), project_id, task_id, parent_id(自引用) | user_id, project_id, task_id, parent_id |
| `attachments` | id, file_name, file_url, file_size, content_type, user_id(FK), task_id, project_id | user_id, task_id, project_id |
| `notifications` | id, user_id, content, type, is_read, related_id | user_id |
| `operation_logs` | id, user_id(FK), action, ip, detail, created_at | user_id |

### 乐观锁字段

- `projects.version` — 每次更新 +1，WHERE 条件携带旧版本号
- `tasks.version` — 同上

---

## 4. API 端点清单 (40+)

### 公开 (4)
| 方法 | 路径 | 说明 |
|---|---|---|
| GET | /api/health | 健康检查 |
| POST | /api/auth/register | 注册 |
| POST | /api/auth/login | 登录 |
| POST | /api/auth/refresh | 刷新 Token |
| GET | /api/auth/captcha | 获取验证码 |

### 认证保护 (5)
| 方法 | 路径 | 说明 |
|---|---|---|
| GET | /api/auth/profile | 个人资料 |
| POST | /api/projects | 创建项目 (PM) |
| GET | /api/projects | 我的项目列表 |
| GET | /api/projects/:id | 项目详情 |
| PUT | /api/projects/:id | 编辑项目 |

### 项目工作流 (6)
| 方法 | 路径 | 说明 |
|---|---|---|
| POST | /api/projects/:id/submit | 提交审批 |
| POST | /api/projects/:id/approve | 审批通过 (PM) |
| POST | /api/projects/:id/start | 开始开发 |
| POST | /api/projects/:id/complete | 完成 |
| POST | /api/projects/:id/archive | 归档 |
| POST | /api/projects/:id/members | 添加成员 |
| DELETE | /api/projects/:id/members/:uid | 移除成员 |
| GET | /api/projects/:id/members | 成员列表 |

### 任务 (7)
| 方法 | 路径 | 说明 |
|---|---|---|
| POST | /api/projects/:id/tasks | 创建任务 |
| GET | /api/projects/:id/tasks | 任务列表 (status/priority/assignee 筛选) |
| GET | /api/tasks/:id | 任务详情 |
| PUT | /api/tasks/:id | 编辑任务 (需 version) |
| DELETE | /api/tasks/:id | 删除任务 |
| PUT | /api/tasks/:id/assign | 指派 |
| PUT | /api/tasks/:id/status | 状态流转 |

### 看板/评论/附件/通知/日志 (各模块)
| 方法 | 路径 | 说明 |
|---|---|---|
| GET | /api/projects/:id/kanban | 看板 (5 列) |
| GET/POST | /api/projects/:id/comments | 项目评论 |
| GET/POST | /api/tasks/:id/comments | 任务评论 |
| POST | /api/comments/:id/reply | 回复评论 |
| DELETE | /api/comments/:id | 删除评论 |
| POST | /api/attachments/upload | 上传文件 |
| GET | /api/tasks/:id/attachments | 任务附件 |
| GET | /api/projects/:id/attachments | 项目附件 |
| DELETE | /api/attachments/:id | 删除附件 |
| GET | /api/notifications | 通知列表 |
| GET | /api/notifications/unread-count | 未读数 |
| PUT | /api/notifications/:id/read | 标记已读 |
| PUT | /api/notifications/read-all | 全部已读 |
| GET | /api/logs | 操作日志 (PM) |

---

## 5. 关键技术决策

### 5.1 乐观锁 (Optimistic Locking)
每条 Task 和 Project 维护 `version` 字段。更新时 `WHERE id=? AND version=?`，若 RowsAffected=0 则返回冲突提示。

### 5.2 Cache-Aside 模式
读取：Redis → 未命中 → MySQL → 写 Redis (10min TTL)  
更新：MySQL → 删除 Redis key  
Redis 故障时自动回源 MySQL（`db.Redis == nil` 检查）

### 5.3 前端响应拦截器
axios 拦截器自动解包 `{code, message, data}` 信封，非零 code 抛出异常，统一错误处理。

### 5.4 Pinia 登录态
`isLoggedIn` 使用 `ref` 而非 `computed`——避免 `localStorage.getItem` 作为非响应式依赖导致缓存永不更新。

### 5.5 GORM 乐观锁 Bug
GORM `Updates(map/struct)` 存在字段过滤行为，改用原生 SQL `db.Exec("UPDATE ... WHERE ...")` 确保 `status` 字段可靠写入。

---

## 6. 构建与部署

### 开发环境
```bash
# 后端
cd backend && go run ./cmd/server   # → :8080

# 前端
cd frontend && npm install && npm run dev   # → :3000

# 基础设施
docker compose up -d mysql redis minio
```

### 生产构建
```bash
# 后端
cd backend && go build -o bin/server ./cmd/server

# 前端
cd frontend && npm run build   # → dist/

# Docker
docker compose up -d
```

---

## 7. 已知限制与后续迭代

- [ ] WebSocket 实时推送（当前看板使用 5s 轮询）
- [ ] MinIO 完整对接（当前上传存储在本地 `./uploads/`）
- [ ] 单元测试覆盖率 ≥80%
- [ ] 邮件通知
- [ ] Kubernetes 部署配置
- [ ] GitHub OAuth 登录
