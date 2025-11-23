package main

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger" // 引入 GORM 的 logger 包
)

func main() {
	// 配置 GORM 日志
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io.Writer，输出到控制台
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别：Silent、Error、Warn、Info
			IgnoreRecordNotFoundError: true,        // 是否忽略 ErrRecordNotFound 错误
			Colorful:                  true,        // 是否开启彩色打印
		},
	)

	dsn := "root:123456@tcp(127.0.0.1:3306)/go_demo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	// 使用db进行后续操作

	db.AutoMigrate(&Product{})

	db.AutoMigrate(&User{})

	// 3. 创建记录时，无需指定 ID
	// 数据库会自动为其分配一个递增的唯一 ID
	newUser := User{Name: "张三", Email: "zhangsan@example.com", Age: 30}
	result := db.Create(&newUser) // 注意：这里要传指针

	if result.Error != nil {
		log.Fatalf("failed to create user: %v", result.Error)
		log.Println("创建用户失败")
	}

	// 4. 创建成功后，GORM 会将数据库自动生成的 ID 回填到 newUser 结构体中
	log.Printf("创建用户成功，自动生成的 ID 为: %d", newUser.ID) // 输出类似: 创建用户成功，自动生成的 ID 为: 1

	user := FindFirst(db, newUser) // 输出类似: 查询到的用户: {ID:1 Name:张三 Email:

	Find(db)

	// 更新用户的年龄
	db.Model(&user).Update("Age", 31) // 更新单个字段
}

func FindFirst(db *gorm.DB, newUser User) User {
	log.Printf("db 的类型是%T", db)
	var user User
	s := db.First(&user, newUser.ID)
	if s.Error != nil {
		log.Fatalf("failed to retrieve user: %v", s.Error)
	}
	log.Printf("查询到的用户:%+v ", user)
	return user
}

func Find(db *gorm.DB) {
	var users []User
	db.Find(&users, "name = ?", "张三")
	log.Printf("查询到的所有用户:%+v ", users)
}

// 定义一个模型
type Product struct {
	gorm.Model
	Code  string
	Price uint
}

type User struct {
	ID    uint `gorm:"primaryKey"` // 明确指定为主键，这是可选的，但推荐写上，增加代码可读性
	Name  string
	Email string
	Age   int
}
