# OnePenny Server

[![version](https://img.shields.io/badge/version-v0.1.0-blue.svg)](https://github.com/iammm0/geekreward-server/releases/tag/v0.1.0)

OnePenny Server 是一个以 Go 为核心、Gin 框架驱动的后端 API 服务，提供悬赏（Bounty）平台的完整功能，包括用户认证、发布赏金、提交申请、团队邀请、通知、评论、点赞、团队管理等。内置 PostgreSQL（GORM）和 Redis 支持，并自动生成 Swagger API 文档。

---

## 特性

- **用户认证**：注册 / 登录 / JWT 鉴权
- **赏金管理**：发布／查询／更新／删除赏金任务
- **申请管理**：提交／查询／删除赏金申请
- **组队邀请**：发送／响应／取消团队邀请
- **通知系统**：发送／查询／未读统计／标记已读
- **评论与回复**：多级评论、编辑、删除
- **点赞机制**：对赏金、评论、用户等多态点赞
- **团队管理**：创建团队、添加／移除成员、成员列表
- **Redis**：缓存、会话或其它扩展
- **Swagger 文档**：自动生成、托管在 `/swagger` 下

---

## 技术栈

- **语言 & 框架**：Go + Gin
- **ORM**：GORM + PostgreSQL
- **缓存**：Redis
- **配置**：Viper（YAML）
- **文档**：swaggo/swag + gin-swagger
- **依赖管理**：Go Modules

---

## 快速上手

### 前置条件

- Go ≥ 1.22
- PostgreSQL
- Redis

### 克隆项目

```bash
git clone https://github.com/yourusername/onepenny-server.git
cd onepenny-server
```

## 配置
### 在项目根目录`config`,示例：

```yaml
server:
  port: 8080

postgres:
  host: localhost
  port: 5432
  user: postgres
  password: onepenny2507031357
  name: onepenny

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0
```

### 安装依赖 & 生成 Swagger 文档

```bash
# 安装 swag CLI（一次性）
go install github.com/swaggo/swag/cmd/swag@latest

# 下载项目依赖
go mod tidy

# 在项目根目录运行，生成 docs/ 目录
swag init --parseDependency --parseInternal
```

### 部署数据库

```bash
docker-compose up -d  # 根目录下
```

### 启动服务

```bash
go run main.go  # 根目录下
```

或者构建二进制：

```bash
go build -o onepenny-server .  
./onepenny-server
```



访问：

- API 根地址：`http://localhost:8080/api/...`
- Swagger UI：`http://localhost:8080/swagger/index.html`

------

## 项目结构

```csharp
onepenny-server/
├── config.yaml
├── go.mod
├── main.go
├── docs/                          # swagger init 生成
├── database/
│   ├── postgres-connection.go     # GORM + Postgres
│   └── redis-connection.go        # Redis 连接
├── controller/
│   ├── router.go                  # 路由集中注册
│   ├── user/
│   │   ├── auth.go
│   │   └── profile.go
│   ├── bounty/
│   │   └── bounty.go
│   ├── application/
│   ├── invitation/
│   ├── notification/
│   ├── comment/
│   ├── like/
│   └── team/
├── internal/
│   ├── repository/                # 数据库操作
│   └── service/                   # 业务逻辑
├── util/
│   ├── jwt.go                     # JWT 生成与验证
│   └── password.go                # bcrypt 密码哈希
└── uploads/                       # 静态文件目录
```

------

## API 文档

项目启动后访问：

```bash
http://localhost:8080/swagger/index.html
```

即可查看并交互式测试所有接口。
