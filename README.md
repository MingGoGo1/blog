# 📝 Blog API

基于 Go + Gin + GORM 的现代化博客系统 API，集成 Redis 认证、Swagger 文档和 Apifox 一键导入功能。

[![Go Version](https://img.shields.io/badge/Go-1.19+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-supported-blue.svg)](docker-compose.yml)

## ✨ 功能特性

- 🔐 用户注册、登录、JWT + Redis 双重认证
- 📄 文章的增删改查、分页查询
- 📁 文件上传功能（MinIO 对象存储）
- 🧪 测试管理功能
- 📖 自动生成 Swagger 文档
- 🎯 Apifox 一键导入
- 🐳 Docker 一键部署
- 🚀 多设备登录支持
- 🛡️ 安全的认证机制

## 🛠️ 技术栈

- **后端框架**: Gin
- **数据库**: MySQL 8.0+
- **缓存**: Redis 6.0+
- **ORM**: GORM v2
- **认证**: JWT + Redis
- **对象存储**: MinIO
- **配置管理**: Viper
- **API文档**: Swagger/OpenAPI 3.0
- **容器化**: Docker & Docker Compose

## 🚀 快速开始

### 1. 环境要求

- Go 1.19+
- MySQL 8.0+
- Redis 6.0+
- MinIO (可选，用于文件存储)
- Docker & Docker Compose (可选)

### 2. 克隆项目

```bash
git clone <repository-url>
cd blog
```

### 3. 安装依赖

```bash
make deps
# 或者
go mod download && go mod tidy
```

### 4. 配置环境

复制配置文件并修改：
```bash
cp config.example.yaml config.yaml
```

修改 `config.yaml` 中的配置信息：
- 数据库连接信息
- Redis 连接信息
- MinIO 配置（如需要）
- JWT 密钥（生产环境必须修改）

⚠️ **安全提醒**: 生产环境部署前请参考 [SECURITY.md](SECURITY.md) 修改所有默认密码和密钥。

### 5. 启动项目

```bash
# 开发模式 (推荐)
make dev

# 或者直接运行
make run
```

### 6. 访问服务

- **API 地址**: http://localhost:8868
- **Swagger 文档**: http://localhost:8868/swagger/index.html
- **Apifox 导入**: http://localhost:8868/apifox

## 📖 API 文档

### 🎯 Apifox 一键导入

1. 访问 http://localhost:8868/apifox
2. 点击 **"🎯 一键导入 Apifox"** 按钮
3. 或手动导入 URL: `http://localhost:8868/swagger/doc.json`

### 📊 API 分组 (14个接口)

| 分组 | 接口数 | 描述 |
|------|--------|------|
| 👤 用户管理 | 3个 | 注册、登录、资料管理 |
| 📄 文章管理 | 5个 | 文章增删改查、列表 |
| 📁 文件管理 | 1个 | 文件上传 |
| 🧪 测试管理 | 5个 | 测试相关接口 |

### 🔑 认证说明

需要认证的接口需要在 Header 中添加：
```
Authorization: Bearer <your-jwt-token>
```

## 🛠️ 开发命令

```bash
# 查看所有可用命令
make help

# 常用命令
make dev          # 开发模式 (格式化+检查+生成文档+运行)
make run          # 直接运行
make test         # 运行测试
make build        # 编译项目
make clean        # 清理文件

# Swagger 文档
make swagger-gen     # 生成文档
make swagger-serve   # 启动服务查看文档
make swagger-clean   # 清理文档

# 快速命令
make quick-start     # 快速启动 (生成文档+运行)
make info           # 显示项目信息
```

## 📁 项目结构

```
blog/
├── docs/                  # Swagger 生成的文档
├── internal/              # 内部包
│   ├── config/           # 配置管理
│   ├── global/           # 全局变量
│   ├── handler/          # HTTP 处理器 (含 Swagger 注释)
│   ├── init/             # 初始化逻辑
│   ├── middleware/       # 中间件
│   ├── router/           # 路由配置
│   ├── service/          # 业务逻辑
│   └── utils/            # 工具函数
├── model/                 # 数据模型
│   ├── entity/           # 数据库实体
│   ├── req/              # 请求模型
│   └── resp/             # 响应模型
├── config.yaml           # 配置文件
├── docker-compose.yml    # Docker 配置
├── main.go               # 程序入口 (含 Swagger 基本信息)
├── Makefile              # 构建脚本
├── swagger.bat           # Windows 批处理脚本
└── README.md             # 项目说明
```

## 🔄 开发流程

### 添加新接口

1. **定义模型** - 在 `model/req` 和 `model/resp` 中定义
2. **实现业务逻辑** - 在 `internal/service` 中实现
3. **创建处理器** - 在 `internal/handler` 中实现并添加 Swagger 注释
4. **注册路由** - 在 `internal/router` 中注册
5. **生成文档** - 运行 `make swagger-gen`
6. **测试接口** - 在 Swagger UI 中测试

### Swagger 注释示例

```go
// CreateArticle 创建文章
// @Summary 创建文章
// @Description 创建新文章
// @Tags 文章管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body req.ArticleCreateRequest true "文章信息"
// @Success 200 {object} entity.Article "创建成功"
// @Failure 400 {object} utils.Response "参数错误"
// @Router /api/v1/articles [post]
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
    // 实现代码
}
```

## 🐳 Docker 部署

```bash
# 构建并启动
make docker-build
make docker-run

# 查看日志
make docker-logs

# 停止服务
make docker-stop
```

## 📝 更新文档

当修改 API 接口后：

```bash
# 重新生成文档
make swagger-gen

# 重启服务
make run

# 重新导入到 Apifox
# 访问 http://localhost:8868/apifox 重新导入
```

## 🎉 特色功能

- ✅ **自动文档生成** - 基于代码注释自动生成
- ✅ **一键导入 Apifox** - 专门的导入页面和协议链接
- ✅ **完美分组** - API 按功能模块清晰分组
- ✅ **交互测试** - Swagger UI 支持直接测试
- ✅ **Windows 友好** - 提供批处理脚本和兼容命令

## 🤝 贡献

我们欢迎任何形式的贡献！请查看 [CONTRIBUTING.md](CONTRIBUTING.md) 了解如何参与项目。

### 贡献者

感谢所有为这个项目做出贡献的开发者！

## 🔒 安全

如果您发现安全漏洞，请查看 [SECURITY.md](SECURITY.md) 了解如何报告。

## 📄 许可证

本项目基于 [MIT License](LICENSE) 开源。

## ⭐ Star History

如果这个项目对您有帮助，请给我们一个 Star！

## 📞 联系我们

- 创建 [Issue](../../issues) 报告问题或建议
- 提交 [Pull Request](../../pulls) 贡献代码
- 查看 [Wiki](../../wiki) 获取更多文档

---

**Happy Coding! 🎉**
