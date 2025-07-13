package resp

import "blog/model/entity"

type ArticleListResponse struct {
	Articles []entity.Article `json:"articles"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"page_size"`
}

// ToArticleResponse 转换为响应格式，隐藏敏感信息
func ToArticleResponse(a *entity.Article) *entity.Article {
	// 隐藏敏感信息
	if a.Author != nil {
		a.Author.Password = ""
		a.Author.Email = ""
	}
	return a
}
