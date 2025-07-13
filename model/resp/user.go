package resp

import "blog/model/entity"

type UserLoginResponse struct {
	Token string       `json:"token"`
	User  *entity.User `json:"user"`
}

type UserProfileResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio"`
}

// ToUserProfileResponse 转换为用户资料响应格式
func ToUserProfileResponse(u *entity.User) *UserProfileResponse {
	return &UserProfileResponse{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Nickname: u.Nickname,
		Avatar:   u.Avatar,
		Bio:      u.Bio,
	}
}
