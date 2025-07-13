package resp

type FileUploadResponse struct {
	ID       uint   `json:"id"`
	FileName string `json:"file_name"`
	FileURL  string `json:"file_url"`
	FileSize int64  `json:"file_size"`
	FileType string `json:"file_type"`
}
