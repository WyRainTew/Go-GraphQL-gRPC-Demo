package database

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// UserInfo 存储用户信息的结构体
type UserInfo struct {
	ID   string `gorm:"primaryKey"`
	Name string
	Age  int
	Sex  string
}

// Database 数据库访问层
type Database struct {
	db *gorm.DB
}

// NewDatabase 创建数据库连接并初始化表
func NewDatabase(dbPath string) (*Database, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 自动迁移 schema
	err = db.AutoMigrate(&UserInfo{})
	if err != nil {
		return nil, err
	}

	// 检查是否已有测试数据
	var count int64
	db.Model(&UserInfo{}).Count(&count)

	// 如果没有数据，插入测试数据
	if count == 0 {
		testUser := UserInfo{
			ID:   "aaa",
			Name: "测试用户",
			Age:  30,
			Sex:  "男",
		}
		result := db.Create(&testUser)
		if result.Error != nil {
			log.Printf("插入测试数据失败: %v", result.Error)
		}
	}

	return &Database{db: db}, nil
}

// GetUserByID 根据ID获取用户信息
func (d *Database) GetUserByID(id string) (*UserInfo, error) {
	// 定义一个UserInfo类型的变量user
	var user UserInfo
	// 在数据库中查找id等于传入参数id的用户，并将结果赋值给user
	result := d.db.First(&user, "id = ?", id)
	// 如果查找过程中出现错误，则返回nil和错误信息
	if result.Error != nil {
		return nil, result.Error
	}
	// 否则返回user和nil
	return &user, nil
}

// Close 关闭数据库连接
func (d *Database) Close() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
