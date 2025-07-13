package resp

import "blog/model/entity"

type TestListResponse struct {
	Tests    []entity.Test `json:"tests"`
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
}
