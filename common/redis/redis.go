package redis

import (
	"GopherAI/config"
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

var Rdb *redis.Client //Rdb是一个全局的Redis客户端连接句柄，供整个项目调用

var ctx = context.Background() //Background()的作用是返回一个空的上下文

func Init() { //Init函数用于初始化Redis连接，读取配置文件中的Redis连接参数，并创建一个新的Redis客户端实例
	conf := config.GetConfig()
	host := conf.RedisConfig.RedisHost
	port := conf.RedisConfig.RedisPort
	password := conf.RedisConfig.RedisPassword
	db := conf.RedisDb
	addr := host + ":" + strconv.Itoa(port)

	Rdb = redis.NewClient(&redis.Options{ //Options是一个结构体，包含了Redis连接的各种配置参数，如地址、密码、数据库编号等
		Addr:     addr,
		Password: password,
		DB:       db,
	})

}

func SetCaptchaForEmail(email, captcha string) error {
	//SetCaptchaForEmail函数用于将生成的验证码存储在Redis中，关联到特定的邮箱地址。
	// 它接受两个参数：email表示用户的邮箱地址，captcha表示生成的验证码字符串。
	key := GenerateCaptcha(email)
	expire := 2 * time.Minute //验证码的过期时间设置为2分钟
	return Rdb.Set(ctx, key, captcha, expire).Err()
}

// CheckCaptchaForEmail函数用于验证用户输入的验证码是否与存储在Redis中的验证码匹配。
func CheckCaptchaForEmail(email, userInput string) (bool, error) {
	key := GenerateCaptcha(email)

	storedCaptcha, err := Rdb.Get(ctx, key).Result() // Get方法用于从Redis中获取与指定键关联的值，如果键不存在或发生错误，会返回一个错误对象。
	if err != nil {
		if err == redis.Nil {

			return false, nil
		}

		return false, err
	}
	// 验证用户输入的验证码是否与存储在Redis中的验证码匹配，使用strings.EqualFold函数进行比较，忽略大小写。
	if strings.EqualFold(storedCaptcha, userInput) {

		// 验证成功后删除 key
		if err := Rdb.Del(ctx, key).Err(); err != nil {

		} else {

		}
		return true, nil
	}

	return false, nil
}
