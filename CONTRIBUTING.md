# 贡献指南

感谢您对 Blog API 项目的关注！我们欢迎任何形式的贡献。

## 🤝 如何贡献

### 报告问题
- 在提交问题前，请先搜索现有的 Issues
- 使用清晰的标题和详细的描述
- 包含复现步骤和环境信息

### 提交代码
1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 📝 开发规范

### 代码风格
- 使用 `go fmt` 格式化代码
- 遵循 Go 官方编码规范
- 添加必要的注释，特别是公共函数

### 提交信息
使用清晰的提交信息格式：
```
type(scope): description

[optional body]

[optional footer]
```

类型包括：
- `feat`: 新功能
- `fix`: 修复bug
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 重构
- `test`: 测试相关
- `chore`: 构建过程或辅助工具的变动

### API 文档
- 为新的 API 接口添加 Swagger 注释
- 更新相关的文档文件
- 确保 API 文档与实际实现一致

## 🧪 测试

### 运行测试
```bash
make test
```

### 添加测试
- 为新功能添加单元测试
- 确保测试覆盖率不降低
- 测试文件命名为 `*_test.go`

## 📋 Pull Request 检查清单

- [ ] 代码已格式化 (`make fmt`)
- [ ] 代码通过检查 (`make vet`)
- [ ] 所有测试通过 (`make test`)
- [ ] 添加了必要的测试
- [ ] 更新了相关文档
- [ ] Swagger 文档已更新

## 🚀 开发环境设置

1. 安装 Go 1.19+
2. 克隆仓库
3. 安装依赖：`make deps`
4. 配置数据库和Redis
5. 启动开发模式：`make dev`

## 📞 联系方式

如有任何问题，请通过以下方式联系：
- 创建 Issue
- 发送邮件到项目维护者

再次感谢您的贡献！🎉
