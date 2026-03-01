# CloudWeGo-Book

基于 [CloudWeGo](https://github.com/cloudwego/biz-demo) 官方示例修改和扩展的微服务实战项目。

本项目在 CloudWeGo 官方 biz-demo 的基础上进行了大量改进和扩展，包含完整的电商系统示例和教程，适合学习 CloudWeGo 技术栈。

## 项目列表

### 1. Gomall - 电商微服务平台

#### 简介
Gomall 是一个教学性质的微服务电商项目，展示了如何使用 CloudWeGo 的 Kitex (RPC) 和 Hertz (HTTP) 框架构建生产级别的微服务架构。项目实现了完整的电商系统，包含 8 个独立服务。

#### 服务架构

| 服务 | 端口 | 类型 | 描述 |
|---------|------|------|-------------|
| **frontend** | 8080 | HTTP | Web 前端服务 |
| **user** | 8881 | RPC | 用户认证与管理 |
| **product** | 8882 | RPC | 商品目录与库存 |
| **cart** | 8883 | RPC | 购物车管理 |
| **payment** | 8884 | RPC | 支付处理 |
| **checkout** | 8885 | RPC | 结算编排 |
| **order** | 8886 | RPC | 订单管理 |
| **email** | 8887 | RPC | 邮件通知 |

#### 基础设施组件

| 组件 | 端口 | 用途 |
|-----------|------|---------|
| MySQL | 3306 | 主数据库 |
| Redis | 6379 | 缓存与会话存储 |
| Consul | 8500 | 服务发现与注册 |
| Jaeger | 16686 | 分布式链路追踪 |
| Prometheus | 9090 | 指标采集 |
| Grafana | 3000 | 监控面板 |
| NATS | 4222 | 消息队列 |
| Loki | 3100 | 日志聚合 |

#### 技术栈

| 技术 | 用途 | 文档 |
|------------|---------|---------------|
| [Kitex](https://github.com/cloudwego/kitex) | 高性能 RPC 框架 | [文档](https://www.cloudwego.io/docs/kitex/) |
| [Hertz](https://github.com/cloudwego/hertz) | 高性能 HTTP 框架 | [文档](https://www.cloudwego.io/docs/hertz/) |
| [cwgo](https://github.com/cloudwego/cwgo) | CloudWeGo 代码生成工具 | [文档](https://www.cloudwego.io/docs/cwgo/) |
| [Bootstrap](https://getbootstrap.com/) | 前端 UI 框架 | [文档](https://getbootstrap.com/docs/) |

#### 已实现功能

- [x] 用户认证与授权
- [x] 用户注册与登录/登出
- [x] 商品分类与目录
- [x] 商品搜索与筛选
- [x] 购物车管理
- [x] 实时购物车徽章更新
- [x] 结算流程
- [x] 支付处理
- [x] 订单管理与历史记录
- [x] 邮件通知 (通过 NATS 异步处理)
- [x] 会话管理
- [x] 分布式链路追踪
- [x] 指标监控

#### 教程章节

项目包含循序渐进的教程，帮助学习者逐步掌握微服务开发：

- **ch01-ch06**: 基础入门
- **ch07-ch10**: 服务拆分与 RPC 通信
- **ch11-ch14**: 中间件集成
- **ch15-ch17**: 生产级特性

#### 详细文档
- [Gomall README](./gomall/README.md)
- [教程说明](./gomall/tutorial/README.md)

---

### 2. Book Shop - 书店系统

#### 简介
展示如何在 Kitex 项目中集成中间件 (如 ElasticSearch、Redis 等)，以及如何在不同复杂度的项目中使用 Hertz 和 Kitex 进行代码分层。

#### 业务场景
电商系统，包含商家管理商品、消费者管理个人账户并下单购买商品。

#### 服务架构
- **facade**: HTTP 服务，处理 HTTP 请求并通过 RPC 调用其他服务
- **user**: RPC 服务，处理用户管理
- **item**: RPC 服务，处理商品管理
- **order**: RPC 服务，处理订单管理

#### 核心技术
- [x] Hertz 作为网关
- [x] Kitex 作为 RPC 框架构建微服务
- [x] Hertz swagger、jwt、pprof、gzip 中间件
- [x] ETCD 服务注册
- [x] MySQL 数据库
- [x] Redis 缓存
- [x] ElasticSearch 搜索引擎

#### 详细文档
[Book Shop](./book-shop/README.md)

---

### 3. Bookinfo

#### 简介
展示如何在 Istio 中使用 Kitex proxyless，以及如何使用 CloudWeGo 实现全流程流量泳道。

#### 业务场景
重写 [Bookinfo](https://istio.io/latest/docs/examples/bookinfo/) 项目，应用分为四个微服务：
- **productpage**: 产品页面微服务，调用 details 和 reviews 微服务
- **details**: 详情微服务，包含书籍信息
- **reviews**: 评论微服务，包含书籍评论，调用 ratings 微服务
- **ratings**: 评分微服务，包含书籍评分信息

#### 核心技术
- [x] 使用 istiod 作为 xDS 服务器
- [x] 使用 wire 依赖注入
- [x] 使用 opentelemetry 链路追踪
- [x] 使用 Kitex-xds 和 opentelemetry baggage 实现 proxyless 流量泳道
- [x] 使用 arco-design react 实现前端界面

#### 详细文档
[Bookinfo](./bookinfo/README.md)

---

### 4. Easy Note - 笔记服务

#### 简介
展示 Hertz 和 Kitex 协作入门，以及项目结构设计。

#### 业务场景
笔记服务，允许用户创建、删除、更新和查询笔记。

#### 服务架构
- **demoapi**: HTTP 服务，处理 HTTP 请求并通过 RPC 调用其他服务
- **demouser**: RPC 服务，处理用户相关操作
- **demonote**: RPC 服务，处理笔记相关操作

#### 核心技术
- [x] 使用 hz 和 kitex 生成代码
- [x] Hertz requestid、jwt、pprof、gzip 中间件
- [x] go-tagexpr 和 thrift-gen-validator 验证请求
- [x] obs-opentelemetry 链路追踪
- [x] etcd 服务注册
- [x] GORM 数据库操作
- [x] MySQL 数据库

#### 详细文档
[easy_note](./easy_note/README.md)

---

### 5. Open Payment Platform - 开放支付平台

#### 简介
展示如何使用 Kitex 泛化调用作为 HTTP 网关，以及如何使用 Kitex 实现 Go 的整洁架构。

#### 核心技术
- [x] Hertz 作为网关
- [x] Kitex 泛化调用客户端路由请求
- [x] Kitex 作为 RPC 框架构建微服务
- [x] 整洁架构设计
- [x] ent 实体框架
- [x] wire 依赖注入
- [x] Nacos 服务注册
- [x] MySQL 数据库

#### 详细文档
[Open Payment Platform](./open-payment-platform/README.md)

---

## 快速开始

### 环境要求

- Go 1.23+
- Docker & Docker Compose
- IDE (VS Code、GoLand 等)

### 启动 Gomall

```bash
# 克隆仓库
git clone https://github.com/LANSGANBS/CloudWeGo-Book.git
cd CloudWeGo-Book/gomall

# 启动基础设施服务
docker-compose up -d

# 初始化环境变量
make init

# 下载依赖
make tidy

# 启动服务
make run svc=user
make run svc=product
make run svc=frontend
# ... 其他服务
```

### 访问应用

- **前端网站**: http://localhost:8080
- **Consul 控制台**: http://localhost:8500
- **Jaeger 链路追踪**: http://localhost:16686
- **Grafana 监控**: http://localhost:3000
- **Prometheus 指标**: http://localhost:9090

## 致谢

本项目基于 [CloudWeGo biz-demo](https://github.com/cloudwego/biz-demo) 进行修改和扩展，感谢 CloudWeGo 团队提供的优秀示例项目。

## 许可证

[Apache License 2.0](LICENSE)
