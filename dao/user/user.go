package user

import (
	"GopherAI/common/mysql"
	"GopherAI/model"
	"GopherAI/utils"
	"context"

	"gorm.io/gorm"
)

// 这些常量主要是用来存放一些固定的字符串，这些字符串主要是用来发送邮件的内容的
const (
	CodeMsg     = "GopherAI验证码如下(验证码仅限于2分钟有效): "
	UserNameMsg = "GopherAI的账号如下，请保留好，后续可以用账号进行登录 "
)

var ctx = context.Background()

// 这边只能通过账号进行登录
func IsExistUser(username string) (bool, *model.User) {

	user, err := mysql.GetUserByUsername(username)

	if err == gorm.ErrRecordNotFound || user == nil {
		return false, nil
	}

	return true, user
}

// 这个函数主要是用来注册用户的，注册成功之后会将用户的信息返回给前端，前端可以根据这些信息进行后续的处理
func Register(username, email, password string) (*model.User, bool) {
	if user, err := mysql.InsertUser(&model.User{
		Email:    email,
		Name:     username,
		Username: username,
		Password: utils.MD5(password), //使用MD5加密密码
	}); err != nil {
		return nil, false
	} else {
		return user, true
	}
}
