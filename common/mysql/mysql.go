package mysql

import (
	"GopherAI/config"
	"GopherAI/model"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 全局数据库连接句柄，供整个项目调用
var DB *gorm.DB

// InitMysql 初始化 MySQL 数据库连接
func InitMysql() error {
	// 1. 从配置文件获取数据库连接参数
	host := config.GetConfig().MysqlHost
	port := config.GetConfig().MysqlPort
	dbname := config.GetConfig().MysqlDatabaseName
	username := config.GetConfig().MysqlUser
	password := config.GetConfig().MysqlPassword
	charset := config.GetConfig().MysqlCharset

	// 2. 拼接 DSN (Data Source Name) 数据库连接字符串
	// 包含：用户名:密码@tcp(地址:端口)/数据库名?配置参数
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local", username, password, host, port, dbname, charset)

	// 3. 配置数据库日志模式
	var log logger.Interface
	if gin.Mode() == "debug" {
		// 开发模式下显示所有 SQL 语句及其执行耗时
		log = logger.Default.LogMode(logger.Info)
	} else {
		log = logger.Default
	}

	// 4. 使用 GORM 打开数据库连接
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，兼容旧版 MySQL
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式
		DontSupportRenameColumn:   true,  // 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动初始化
	}), &gorm.Config{
		Logger: log,
	})
	if err != nil {
		return err
	}

	// 5. 获取底层的 sql.DB 对象以配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)           // 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxOpenConns(100)          // 设置打开数据库连接的最大数量
	sqlDB.SetConnMaxLifetime(time.Hour) // 设置连接可复用的最大时间

	// 将初始化好的连接赋值给全局变量 DB
	DB = db

	// 6. 执行数据库迁移（自动建表）
	return migration()
}

// migration 自动迁移函数，根据模型结构体自动创建或更新数据库表
func migration() error {
	return DB.AutoMigrate(
		new(model.User),    // 用户表
		new(model.Session), // 会话表
		new(model.Message), // 消息表
	)
}

// InsertUser 插入一个新用户到数据库
func InsertUser(user *model.User) (*model.User, error) {
	err := DB.Create(&user).Error
	return user, err
}

// GetUserByUsername 根据用户名从数据库获取用户信息
func GetUserByUsername(username string) (*model.User, error) {
	user := new(model.User)
	// 使用 First 查找第一条符合条件的记录
	err := DB.Where("username = ?", username).First(user).Error
	return user, err
}
