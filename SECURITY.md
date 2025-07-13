# 安全政策

## 🔒 安全配置

### 生产环境部署前必须修改的配置

#### 1. JWT 密钥
```yaml
# config.yaml
jwt:
  secret: "your-strong-jwt-secret-key-at-least-32-characters"
```

#### 2. 数据库密码
```yaml
# config.yaml
database:
  password: "your-strong-mysql-password"
```

#### 3. Redis 密码（推荐）
```yaml
# config.yaml
redis:
  password: "your-redis-password"
```

#### 4. MinIO 访问密钥
```yaml
# config.yaml
minio:
  access_key_id: "your-minio-access-key"
  secret_access_key: "your-strong-minio-secret-key"
```

#### 5. 服务器模式
```yaml
# config.yaml
server:
  mode: "release"  # 生产环境使用 release 模式
```

## 🛡️ 安全最佳实践

### 1. 环境变量
推荐使用环境变量存储敏感信息：
```bash
export JWT_SECRET="your-jwt-secret"
export DB_PASSWORD="your-db-password"
export REDIS_PASSWORD="your-redis-password"
```

### 2. HTTPS
生产环境务必使用 HTTPS：
- 配置反向代理（Nginx/Apache）
- 使用有效的 SSL 证书
- 强制 HTTPS 重定向

### 3. 防火墙
- 只开放必要的端口
- 限制数据库和Redis的访问
- 使用VPC或私有网络

### 4. 定期更新
- 定期更新依赖包
- 关注安全漏洞公告
- 及时应用安全补丁

## 🚨 报告安全漏洞

如果您发现安全漏洞，请：

1. **不要**在公开的 Issue 中报告
2. 发送邮件到项目维护者
3. 提供详细的漏洞描述和复现步骤
4. 我们会在24小时内回复

## 📋 安全检查清单

部署前请确认：

- [ ] 修改了所有默认密码
- [ ] JWT 密钥足够强壮（至少32字符）
- [ ] 数据库不允许外网访问
- [ ] Redis 设置了密码认证
- [ ] 服务器模式设置为 release
- [ ] 配置了 HTTPS
- [ ] 设置了适当的防火墙规则
- [ ] 定期备份数据库
- [ ] 监控系统日志

## 🔍 安全功能

本项目已实现的安全功能：

- ✅ JWT Token 认证
- ✅ Redis Token 验证
- ✅ 密码 bcrypt 加密
- ✅ SQL 注入防护（GORM）
- ✅ CORS 跨域保护
- ✅ 请求日志记录
- ✅ 统一错误处理

## 📚 相关资源

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [Go Security Checklist](https://github.com/Checkmarx/Go-SCP)
- [JWT Best Practices](https://auth0.com/blog/a-look-at-the-latest-draft-for-jwt-bcp/)
