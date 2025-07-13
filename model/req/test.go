package req

type TestCreateRequest struct {
	Id   uint   `json:"id"`
	Test string `json:"test" binding:"max=500"`
}

type TestUpdateRequest struct {
	Id   uint   `json:"id"`
	Test string `json:"test" binding:"max=500"`
}

type TestListRequest struct {
	Test     string `form:"test" binding:"max=500"`
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
}
