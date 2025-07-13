# Blog API Makefile (Windows Compatible)
# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
BINARY_NAME=blog
BINARY_WIN=$(BINARY_NAME).exe

# Help
help:
	@echo "Blog API 可用命令:"
	@echo ""
	@echo "基础命令:"
	@echo "  build          编译项目"
	@echo "  clean          清理构建文件"
	@echo "  test           运行测试"
	@echo "  deps           下载依赖"
	@echo "  run            运行项目"
	@echo "  dev            开发模式 (格式化+检查+编译+运行)"
	@echo ""
	@echo "代码质量:"
	@echo "  fmt            格式化代码"
	@echo "  vet            代码检查"
	@echo ""
	@echo "Swagger文档:"
	@echo "  swagger-gen    生成Swagger文档"
	@echo "  swagger-serve  启动服务查看文档"
	@echo "  swagger-clean  清理文档文件"
	@echo "  swagger-install 安装Swagger工具"
	@echo ""
	@echo "Docker:"
	@echo "  docker-build   构建Docker镜像"
	@echo "  docker-run     启动Docker容器"
	@echo "  docker-stop    停止Docker容器"

# Build the project
build:
	@echo "🔨 编译项目..."
	$(GOBUILD) -o $(BINARY_WIN) -v .
	@echo "✅ 编译完成: $(BINARY_WIN)"

# Clean build files (Windows compatible)
clean:
	@echo "🧹 清理构建文件..."
	$(GOCLEAN)
	@if exist $(BINARY_WIN) del $(BINARY_WIN)
	@if exist docs\docs.go del docs\docs.go
	@if exist docs\swagger.json del docs\swagger.json
	@if exist docs\swagger.yaml del docs\swagger.yaml
	@echo "✅ 清理完成"

# Run tests
test:
	@echo "🧪 运行测试..."
	$(GOTEST) -v ./...
	@echo "✅ 测试完成"

# Download dependencies
deps:
	@echo "📦 下载依赖..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "✅ 依赖更新完成"

# Run the application
run:
	@echo "🚀 启动项目..."
	$(GOCMD) run main.go

# Format code
fmt:
	@echo "🎨 格式化代码..."
	$(GOCMD) fmt ./...
	@echo "✅ 代码格式化完成"

# Vet code
vet:
	@echo "🔍 代码检查..."
	$(GOCMD) vet ./...
	@echo "✅ 代码检查完成"

# Development workflow
dev:
	@echo "🛠️  开发模式启动..."
	@make deps
	@make fmt
	@make vet
	@make swagger-gen
	@make run

# Docker operations
docker-build:
	@echo "🐳 构建Docker镜像..."
	docker build -t blog-api .
	@echo "✅ Docker镜像构建完成"

docker-run:
	@echo "🐳 启动Docker容器..."
	docker-compose up -d
	@echo "✅ Docker容器启动完成"

docker-stop:
	@echo "🐳 停止Docker容器..."
	docker-compose down
	@echo "✅ Docker容器已停止"

docker-logs:
	@echo "📋 查看Docker日志..."
	docker-compose logs -f blog-api

# Swagger Documentation Operations
swagger-gen:
	@echo "📋 生成Swagger文档..."
	swag init -g main.go -o ./docs --parseDependency --parseInternal
	@echo "✅ Swagger文档生成完成"
	@echo "📖 访问地址: http://localhost:8868/swagger/index.html"
	@echo "🎯 Apifox导入: http://localhost:8868/apifox"

swagger-serve:
	@echo "🚀 启动服务查看Swagger文档..."
	@echo "📖 Swagger文档: http://localhost:8868/swagger/index.html"
	@echo "🎯 Apifox导入页面: http://localhost:8868/apifox"
	@echo "🔗 导入URL: http://localhost:8868/swagger/doc.json"
	$(GOCMD) run main.go

swagger-clean:
	@echo "🧹 清理Swagger文档..."
	@if exist docs\docs.go del docs\docs.go
	@if exist docs\swagger.json del docs\swagger.json
	@if exist docs\swagger.yaml del docs\swagger.yaml
	@echo "✅ Swagger文档清理完成"

swagger-install:
	@echo "📦 安装Swagger工具..."
	$(GOCMD) install github.com/swaggo/swag/cmd/swag@latest
	@echo "✅ Swagger工具安装完成"

# Quick commands for common workflows
quick-start:
	@echo "🚀 快速启动 (生成文档 + 启动服务)..."
	@make swagger-gen
	@make run

quick-dev:
	@echo "🛠️  快速开发 (完整开发流程)..."
	@make dev

# Show project info
info:
	@echo "📊 项目信息:"
	@echo "  项目名称: Blog API"
	@echo "  Go版本: $(shell go version)"
	@echo "  项目路径: $(shell pwd)"
	@echo "  二进制文件: $(BINARY_WIN)"
	@echo ""
	@echo "📖 文档地址:"
	@echo "  Swagger UI: http://localhost:8868/swagger/index.html"
	@echo "  Apifox导入: http://localhost:8868/apifox"
	@echo "  导入URL: http://localhost:8868/swagger/doc.json"

.PHONY: help build clean test deps run fmt vet dev docker-build docker-run docker-stop docker-logs swagger-gen swagger-serve swagger-clean swagger-install quick-start quick-dev info
