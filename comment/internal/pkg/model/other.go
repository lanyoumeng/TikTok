package model

// 用于序列化redis
type RComment struct {
	Id      int64  `json:"id,omitempty"`
	UserId  int64  `json:"user_id"`
	Content string `json:"content,omitempty"`
}
