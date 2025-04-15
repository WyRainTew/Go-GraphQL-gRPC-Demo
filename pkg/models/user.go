package models

// UserInfo 存储用户信息的结构体
type UserInfo struct {
	ID   string `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
	Sex  string `json:"sex"`
} 