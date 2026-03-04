package redis

import (
	"GopherAI/config"
	"fmt"
)

// 这个文件的作用是定义Redis相关的键生成函数，键是用于在Redis中存储和检索数据的标识符。在这个文件中，我们主要关注与验证码相关的键生成逻辑。
// GenerateCaptcha函数的主要作用是根据给定的邮箱地址生成一个唯一的验证码键。这个键通常用于在Redis中存储与特定邮箱相关的验证码数据，以便后续验证用户输入的验证码是否正确。
// 通过使用配置文件中的前缀格式，可以确保生成的键具有一致的命名规范，便于在Redis中存储和检索验证码数据。
// 例如，如果配置文件中的CaptchaPrefix设置为"captcha:%s"，当调用GenerateCaptcha("user@example.com")时，会生成键"captcha:user@example.com"。

// github.com/go-redis/redis/v8
// key:特定邮箱-> 验证码
func GenerateCaptcha(email string) string {
	return fmt.Sprintf(config.DefaultRedisKeyConfig.CaptchaPrefix, email)
}
