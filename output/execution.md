# cy-402 律师案件管理系统 执行记录

- 项目编号/名称：cy-402 律师案件管理系统（cylawcase）
- 日期：2026-08-16
- 短名：cylawcase
- 端口：前端 28031 / 后端 29069 / PostgreSQL 38102
- 技术栈：React 18 + TypeScript + Ant Design + Vite + Zustand；Go 1.22 + Gin + GORM；PostgreSQL 16；JWT + RBAC

## Docker Compose 结果

| 容器 | 状态 | 端口 |
| --- | --- | --- |
| cylawcase-db | Up (healthy) | 38102:5432 |
| cylawcase-backend | Up (healthy) | 29069:8080 |
| cylawcase-frontend | Up | 28031:80 |

`docker compose config --quiet` 通过；`docker compose up -d --build` 一键启动成功。

## 关键 API 冒烟结果（38/38 通过）

| 接口 | 方法 | 状态码 | 结果摘要 |
| --- | --- | --- | --- |
| /healthz | GET | 200 | ok |
| /api/healthz（Nginx 反代） | GET | 200 | ok |
| /api/v1/healthz（Nginx 反代） | GET | 200 | ok |
| /api/v1/auth/register | POST | 200 | 注册成功 |
| /api/v1/auth/login | POST | 200 | 返回 JWT + 用户 |
| /api/v1/users/me | GET | 200 | 当前用户 |
| /api/v1/users/me（未授权） | GET | 401 | 未登录 |
| /api/v1/users（非管理员） | GET | 403 | 越权拦截 |
| /api/v1/clients?page_size=5 | GET | 200 | 客户列表 |
| /api/v1/clients | POST | 200 | 新建客户 |
| /api/v1/clients/:id | GET | 200 | 客户详情+历史案件 |
| /api/v1/clients/:id | PUT | 200 | 编辑客户 |
| /api/v1/clients?keyword= | GET | 200 | 检索客户 |
| /api/v1/cases?page_size=5 | GET | 200 | 案件列表 |
| /api/v1/cases | POST | 200 | 创建案件 |
| /api/v1/cases/:id | GET | 200 | 案件详情 |
| /api/v1/cases/:id/status（filed→investigating） | POST | 200 | 状态流转 |
| /api/v1/cases/:id/status（跳跃） | POST | 409 | 状态机拦截 |
| /api/v1/cases/:id/status（investigating→hearing→closed） | POST | 200 | 状态流转 |
| /api/v1/cases/:id/assign | POST | 200 | 分配律师 |
| /api/v1/cases?lead_lawyer_id=2 | GET | 200 | 按律师筛选 |
| /api/v1/documents | POST | 200 | 上传文档记录 |
| /api/v1/documents/by-case/:id | GET | 200 | 按案件查文档 |
| /api/v1/documents?page_size=5 | GET | 200 | 文档中心 |
| /api/v1/documents/:id | DELETE | 200 | 删除文档 |
| /api/v1/billings | POST | 200 | 创建账单 |
| /api/v1/billings?page_size=5 | GET | 200 | 账单列表 |
| /api/v1/billings/:id/paid | POST | 200 | 标记支付 |
| /api/v1/billings/:id/paid（重复） | POST | 409 | 状态冲突 |
| /api/v1/billings/:id/invoiced | POST | 200 | 开票 |
| /api/v1/billings/summary | GET | 200 | 本月应收/已收/待收 |
| /api/v1/billings/by-case/:id | GET | 200 | 按案件查账单 |
| /api/v1/billings/:id/void | POST | 200 | 作废账单 |
| /api/v1/auth/login（admin） | POST | 200 | 管理员登录 |
| /api/v1/audit-logs（admin） | GET | 200 | 审计日志 |
| /api/v1/audit-logs（律师） | GET | 403 | 越权拦截 |
| /api/v1/upload/file | POST | 200 | 文件上传返回 url |

异常路径覆盖：401 未授权、403 越权（用户列表/审计日志）、409 状态机冲突（案件/账单）、404 无效资源。

## 浏览器验证结论（内置 Playwright，无外部 Chrome）

- /login 登录页打开正常，输入 lawyer/User@123 登录后进入案件列表，头部显示「张律师」与全部菜单（案件列表/客户管理/费用中心/文档中心/审计日志/个人中心）。
- /cases 渲染真实案件数据：CY20260001 华信科技买卖合同纠纷、CY20260002 陈晓明民间借贷纠纷、CY20260003，主办律师为张律师。
- /billing 渲染本月应收/已收/待收统计卡片与账单列表（BILL2026080001 已支付、0002 待支付、0003 已开票，律师费/诉讼费）。
- /clients 渲染真实客户（深圳华信科技有限公司、陈晓明、冒烟客户公司(改)）。
- 截图：output/cylawcase_cases.png、output/cylawcase_cases2.png。

## README 检查项

- Docker Compose 一键启动命令在最前；本地开发命令；技术栈表格（后端 Go 1.22 + Gin + GORM）；目录结构；环境变量；部署说明；License。
- 枚举出现位置清单：CaseStatus、BillingType、BillingStatus 前后端出现位置已列出。

## 其他质量项

- 后端 `go build ./...` 通过；`go vet` 通过。
- 单元测试：internal/service 与 internal/util 表驱动测试通过（go test ./... ok）。
- 前端 `npm run build` 零错误（tsc + vite build 通过）。
- database/init.sql 含建表（与 GORM AutoMigrate 兼容的命名唯一约束 + setval 重置序列）与种子数据（3 用户 + 2 客户 + 3 案件 + 3 文档 + 3 账单 + 审计日志）。

- 提交记录：init commit 44d7996；本文件独立提交 docs commit。
