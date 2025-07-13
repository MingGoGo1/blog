# Redis Token认证实现

## 概述

已成功将JWT token验证机制与Redis集成，实现了更安全和高效的用户认证系统。

## 主要改进

### 1. 登录时的Redis存储
- **原来**: 只生成JWT token，无Redis存储
- **现在**: 登录成功后将token存储到Redis
  ```go
  // 存储token信息，包含用户ID
  tokenKey := fmt.Sprintf("token:%s", token)
  userTokensKey := fmt.Sprintf("user_tokens:%d", user.ID)
  
  global.Redis.Set(ctx, tokenKey, user.ID, 24*time.Hour*30)
  global.Redis.SAdd(ctx, userTokensKey, token)
  global.Redis.Expire(ctx, userTokensKey, 24*time.Hour*30)
  ```

### 2. Token验证机制
- **原来**: 只解析JWT，然后查询数据库
- **现在**: 先检查Redis中token是否存在，再查询数据库
  ```go
  // 检查Redis中是否存在该token
  tokenKey := fmt.Sprintf("token:%s", token)
  userIDStr, err := global.Redis.Get(ctx, tokenKey).Result()
  if err != nil {
      return nil, errors.New("token已过期或无效")
  }
  ```

### 3. 用户注销功能
- **新增**: 注销时从Redis删除token
  ```go
  // 删除token
  tokenKey := fmt.Sprintf("token:%s", token)
  userTokensKey := fmt.Sprintf("user_tokens:%d", claims.UserID)
  
  global.Redis.Del(ctx, tokenKey)
  global.Redis.SRem(ctx, userTokensKey, token)
  ```

## 新增功能

### 1. 单设备注销
- 路由: `POST /api/v1/logout`
- 功能: 注销当前设备的token

### 2. 多设备登录支持
- 每个设备有独立的token
- 支持同一用户在多个设备同时登录

### 3. 全设备注销
- 方法: `LogoutAllDevices()`
- 功能: 注销用户所有设备的token

## 数据结构

### Redis存储结构
```
token:{token_string} -> {user_id}           # 单个token映射
user_tokens:{user_id} -> Set{token1, token2, ...}  # 用户所有token集合
```

### 过期时间
- Token过期时间: 30天 (24*time.Hour*30)
- 与JWT token过期时间保持一致

## 安全优势

1. **即时注销**: 注销后token立即失效，无需等待JWT过期
2. **集中管理**: 可以统一管理用户的所有token
3. **性能优化**: Redis查询比数据库查询更快
4. **防重放攻击**: 注销的token无法再次使用

## API接口

### 登录
```bash
POST /api/v1/login
{
  "username": "admin",
  "password": "123456"
}
```

### 注销
```bash
POST /api/v1/logout
Authorization: Bearer {token}
```

### 访问受保护资源
```bash
GET /api/v1/profile
Authorization: Bearer {token}
```

## 测试

测试文件位于 `test/` 目录：
- `test_redis_auth.go`: 单元测试
- `test_api_auth.md`: API测试说明

## 配置要求

确保 `config.yaml` 中Redis配置正确：
```yaml
redis:
  host: "localhost"
  port: "6379"
  password: ""
  db: 7
```

## 部署注意事项

1. 确保Redis服务正常运行
2. 生产环境建议设置Redis密码
3. 考虑Redis持久化配置
4. 监控Redis内存使用情况
