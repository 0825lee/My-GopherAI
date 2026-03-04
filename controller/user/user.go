package user

import (
	"GopherAI/common/code"
	"GopherAI/controller"
	"GopherAI/service/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ( //这个type的作用是为了定义一些结构体，这些结构体主要是用来接收前端传入的参数以及返回给前端的数据的
	//这里的Username只能是账号登录
	LoginRequest struct {
		Username string `json:"username"`
		Password string `json:password`
	}
	// omitempty当字段为空的时候，不返回这个东西
	LoginResponse struct {
		controller.Response
		Token string `json:"token,omitempty"`
	}
	//验证码由后端生成，存放到redis中，固然需要先发送一次请求CaptchaRequest,然后用返回的验证码
	//邮箱以及密码进行注册，后续再将账号进行返回
	RegisterRequest struct {
		Email    string `json:"email" binding:"required"`
		Captcha  string `json:"captcha"`
		Password string `json:"password"`
	}
	//注册成功之后，直接让其进行登录状态
	RegisterResponse struct {
		controller.Response
		Token string `json:"token,omitempty"`
	}
	//发送验证码的请求结构体
	CaptchaRequest struct {
		Email string `json:"email" binding:"required"`
	}
	//发送验证码的响应结构体
	CaptchaResponse struct {
		controller.Response
	}
)

// 这里的登录接口需要前端传入账号和密码，后端会先验证账号是否存在
// 存在的话就验证密码是否正确
func Login(c *gin.Context) {
	//1:先进行参数校验，为了保证前端传入的参数是合法的，合法的话就进行下一步的处理
	req := new(LoginRequest)
	res := new(LoginResponse)
	if err := c.ShouldBindJSON(req); err != nil { //ShouldBindJSON是gin框架提供的一个方法，可以将前端传入的json数据绑定到结构体上去
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return
	}
	//2:再进行登录的处理，登录成功的话就返回一个Token，登录失败的话就返回对应的错误码
	token, code_ := user.Login(req.Username, req.Password)
	if code_ != code.CodeSuccess {
		c.JSON(http.StatusOK, res.CodeOf(code_))
		return
	}

	res.Success()
	res.Token = token
	c.JSON(http.StatusOK, res) //将登录成功和token返回给前端

}

// 这里的注册接口需要前端传入邮箱，密码以及验证码
// 后端会先验证验证码是否正确，正确的话就生成一个账号
func Register(c *gin.Context) {

	req := new(RegisterRequest)
	res := new(RegisterResponse)
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return
	}

	token, code_ := user.Register(req.Email, req.Password, req.Captcha)
	if code_ != code.CodeSuccess {
		c.JSON(http.StatusOK, res.CodeOf(code_))
		return
	}
	//注册成功之后直接返回一个Token
	res.Success()
	res.Token = token
	c.JSON(http.StatusOK, res)
}

// 发送验证码的接口，前端只需要传入邮箱地址，后端会生成一个验证码
// 存放到redis中，并且将验证码发送到对应的邮箱上去
func HandleCaptcha(c *gin.Context) {
	req := new(CaptchaRequest)
	res := new(CaptchaResponse)
	//解析参数
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return
	}

	//给service层进行处理
	code_ := user.SendCaptcha(req.Email) //发送验证码到对应邮箱上去
	if code_ != code.CodeSuccess {
		c.JSON(http.StatusOK, res.CodeOf(code_))
		return
	}
	//匿名字段，其实本身res.Success()调用就是res.Response.Success()
	//res.Response.Success()
	res.Success()
	c.JSON(http.StatusOK, res)
}
