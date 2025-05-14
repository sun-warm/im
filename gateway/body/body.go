package body

type LoginRequest struct {
	UserName string `json:"user_name"`
	PassWord string `json:"pass_word"`
}

type LoginResponse struct {
	StatusCode string `json:"status_code"`
}
type SendMessageRequest struct {
	UserName    string `json:"user_name"`
	Receiver    string `json:"receiver"`
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
}

type SendMessageResponse struct {
	StatusCode int `json:"status_code"`
}
