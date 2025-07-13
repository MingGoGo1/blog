package service

import (
	"context"
	"errors"
	"fmt"

	"blog/internal/global"
	"blog/internal/utils"
	"blog/model/entity"
	"blog/model/req"
	"blog/model/resp"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

// getDB 获取数据库连接，支持事务
func (s *UserService) getDB(ctx context.Context) *gorm.DB {
	return global.GetDB(ctx)
}

// Register 用户注册
func (s *UserService) Register(ctx context.Context, req *req.UserRegisterRequest) (*entity.User, error) {
	db := s.getDB(ctx)

	// 检查用户名是否存在
	var existUser entity.User
	if err := db.Where("username = ? OR email = ?", req.Username, req.Email).First(&existUser).Error; err == nil {
		return nil, errors.New("用户名或邮箱已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建用户
	user := &entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Nickname: req.Nickname,
		Status:   1,
	}

	if user.Nickname == "" {
		user.Nickname = user.Username
	}

	if err := db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// Login 用户登录
func (s *UserService) Login(ctx context.Context, req *req.UserLoginRequest) (*resp.UserLoginResponse, error) {
	db := s.getDB(ctx)

	var user entity.User
	if err := db.Where("username = ? OR email = ?", req.Username, req.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("密码错误")
	}

	// 生成JWT token
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	// 将token存储到Redis，支持多设备登录
	tokenKey := fmt.Sprintf("token:%s", token)
	userTokensKey := fmt.Sprintf("user_tokens:%d", user.ID)

	// 获取JWT过期时间
	expireDuration := utils.GetJWTExpireDuration()

	// 存储token信息，包含用户ID
	global.Redis.Set(ctx, tokenKey, user.ID, expireDuration)
	// 将token添加到用户的token集合中
	global.Redis.SAdd(ctx, userTokensKey, token)
	// 设置用户token集合的过期时间
	global.Redis.Expire(ctx, userTokensKey, expireDuration)

	return &resp.UserLoginResponse{
		Token: token,
		User:  &user,
	}, nil
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(ctx context.Context, id uint) (*entity.User, error) {
	db := s.getDB(ctx)

	var user entity.User
	if err := db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// ValidateToken 验证token
func (s *UserService) ValidateToken(token string) (*entity.User, error) {
	// 1. 解析JWT token获取用户ID
	claims, err := utils.ParseToken(token)
	if err != nil {
		return nil, err
	}

	// 2. 检查Redis中是否存在该token
	ctx := context.Background()
	tokenKey := fmt.Sprintf("token:%s", token)
	userIDStr, err := global.Redis.Get(ctx, tokenKey).Result()
	if err != nil {
		return nil, errors.New("token已过期或无效")
	}

	// 3. 验证token中的用户ID与Redis中存储的是否一致
	if fmt.Sprintf("%d", claims.UserID) != userIDStr {
		return nil, errors.New("token无效")
	}

	// 4. 从数据库获取用户信息
	return s.GetUserByID(context.Background(), claims.UserID)
}

// Logout 用户注销
func (s *UserService) Logout(ctx context.Context, token string) error {
	// 解析token获取用户ID
	claims, err := utils.ParseToken(token)
	if err != nil {
		return err
	}

	// 从Redis中删除token
	tokenKey := fmt.Sprintf("token:%s", token)
	userTokensKey := fmt.Sprintf("user_tokens:%d", claims.UserID)

	// 删除token
	global.Redis.Del(ctx, tokenKey)
	// 从用户token集合中移除该token
	global.Redis.SRem(ctx, userTokensKey, token)

	return nil
}

// LogoutAllDevices 注销用户所有设备
func (s *UserService) LogoutAllDevices(ctx context.Context, userID uint) error {
	userTokensKey := fmt.Sprintf("user_tokens:%d", userID)

	// 获取用户所有token
	tokens, err := global.Redis.SMembers(ctx, userTokensKey).Result()
	if err != nil {
		return err
	}

	// 删除所有token
	for _, token := range tokens {
		tokenKey := fmt.Sprintf("token:%s", token)
		global.Redis.Del(ctx, tokenKey)
	}

	// 删除用户token集合
	global.Redis.Del(ctx, userTokensKey)

	return nil
}

// UpdateProfile 更新用户资料
func (s *UserService) UpdateProfile(ctx context.Context, userID uint, updates map[string]interface{}) error {
	db := s.getDB(ctx)
	return db.Model(&entity.User{}).Where("id = ?", userID).Updates(updates).Error
}
