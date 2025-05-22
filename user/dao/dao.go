package dao

import (
	"fmt"
	"user/dao/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	var err error
	//DB, err = gorm.Open(mysql.Open("root:sw19990807@tcp(192.168.145.133:3306)/mysql?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	// 检查是否有错误
	DB, err = gorm.Open(mysql.Open("root:sw19990807@tcp(localhost:3306)/mysql?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	// 检查是否有错误
	if err != nil {
		fmt.Println("连接数据库失败：", err)
		return err
	}
	fmt.Println("连接mysql成功")
	// 设置表选项
	err = DB.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4").AutoMigrate(&model.UserRelation{})
	if err != nil {
		return err
	}
	//modifyData()
	//创建表
	//result := DB.Migrator().CreateTable(&model.User{})
	//result = DB.Migrator().CreateTable(&model.WorkTime{})
	return nil
}

func modifyData() {
	//User := model.User{UserName: "y-zhou.weihua", Name: "周卫华", Role: 2, Department: "质量控制部"}
	//DB.Create(&User)
	//User = model.User{UserName: "y-sun.wen", Name: "孙文", Role: 1, Department: "质量控制部"}
	//DB.Create(&User)
	//User = model.User{UserName: "y-xu.lei", Name: "许雷", Role: 2, Department: "质量控制部"}
	//DB.Create(&User)

	//WorkTime := model.WorkTime{UserName: "y-zhou.weihua", WorkTime: 114, Department: "质量控制部", YearNMonth: "2024-02"}
	//DB.Create(&WorkTime)
	//WorkTime = model.WorkTime{UserName: "y-sun.wen", WorkTime: 113, Department: "质量控制部", YearNMonth: "2024-02"}
	//DB.Create(&WorkTime)
	//WorkTime = model.WorkTime{UserName: "y-xu.lei", WorkTime: 115, Department: "质量控制部", YearNMonth: "2024-02"}
	//DB.Create(&WorkTime)
	//
	//WorkTime = model.WorkTime{UserName: "y-zhou.weihua", WorkTime: 113, Department: "质量控制部", YearNMonth: "2025-02"}
	//DB.Create(&WorkTime)
	//WorkTime = model.WorkTime{UserName: "y-sun.wen", WorkTime: 112, Department: "质量控制部", YearNMonth: "2025-02"}
	//DB.Create(&WorkTime)
	//WorkTime = model.WorkTime{UserName: "y-xu.lei", WorkTime: 111, Department: "质量控制部", YearNMonth: "2025-02"}
	//DB.Create(&WorkTime)
}
