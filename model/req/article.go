package req

type ArticleCreateRequest struct {
	Title      string `json:"title" binding:"required,max=200"`
	Content    string `json:"content" binding:"required"`
	Summary    string `json:"summary" binding:"max=500"`
	CoverImage string `json:"cover_image" binding:"max=255"`
	Status     int    `json:"status" binding:"oneof=0 1"`
	Tags       string `json:"tags" binding:"max=255"`
}

type ArticleUpdateRequest struct {
	Title      string `json:"title" binding:"max=200"`
	Content    string `json:"content"`
	Summary    string `json:"summary" binding:"max=500"`
	CoverImage string `json:"cover_image" binding:"max=255"`
	Status     int    `json:"status" binding:"oneof=0 1"`
	Tags       string `json:"tags" binding:"max=255"`
}

type ArticleListRequest struct {
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
	Status   int    `form:"status" binding:"oneof=0 1"`
	AuthorID uint   `form:"author_id"`
	Keyword  string `form:"keyword"`
}
