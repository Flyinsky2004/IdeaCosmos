package service

import (
	"back-end/config"
	"back-end/entity/dto"
	"back-end/entity/pojo"
	"back-end/util"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"time"
)

/**
 * @author Flyinsky
 * @email w2084151024@gmail.com
 * @date 2024/12/24 14:44
 */
type UserLoginBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var body UserLoginBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(200, dto.ErrorResponse[string](400, "提交了错误的JSON"))
		return
	}
	var user pojo.User
	if err := config.MysqlDataBase.Where("username = ?", body.Username).First(&user).Error; err != nil {
		c.JSON(200, dto.ErrorResponse[string](401, "用户名不存在"))
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		c.JSON(200, dto.ErrorResponse[string](401, "用户名或密码错误"))
		return
	}
	var token, err = util.GenerateToken(int(user.ID), user.Username)
	if err != nil {
		c.JSON(200, dto.ErrorResponse[string](401, "生成用户token时发生错误，请再尝试一次，或联系管理员"))
		return
	}
	c.JSON(200, dto.SuccessResponse(token))
}

// SaveCodeToRedis 保存验证码到 Redis
func SaveCodeToRedis(key, code string, ttl time.Duration) error {
	return config.RedisClient.Set(context.Background(), key, code, ttl).Err()
}

// GetCodeFromRedis 获取验证码
func GetCodeFromRedis(key string) (string, error) {
	return config.RedisClient.Get(context.Background(), key).Result()
}
func DeleteCodeToRedis(key string) error {
	return config.RedisClient.Del(context.Background(), key).Err()
}

// CheckIfCodeExists 检查验证码是否存在
func CheckIfCodeExists(key string) (bool, error) {
	cmd := config.RedisClient.Exists(context.Background(), key)
	exists, err := cmd.Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

type RegisterRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Code     string `json:"code"`
}
type SendVerifyCodeRequestBody struct {
	Email string `json:"email"`
}

func SendVerifyCode(c *gin.Context) {
	var reqBody SendVerifyCodeRequestBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(200, dto.ErrorResponse[string](400, "请求表单格式有误"))
		return
	}
	isHasExisted, err := CheckIfCodeExists(reqBody.Email)
	if err != nil {
		c.JSON(200, dto.ErrorResponse[string](400, "查询验证码时发生错误"))
		return
	}
	if isHasExisted {
		c.JSON(200, dto.ErrorResponse[string](400, "你已经发送了一条未失效的验证码"))
		return
	}
	code := util.GenerateCode(6)
	if err := SaveCodeToRedis(reqBody.Email, code, time.Minute*3); err != nil {
		c.JSON(200, dto.ErrorResponse[string](500, "缓存验证码时发生错误"))
		return
	}
	if err := util.SendEmail(reqBody.Email, "OneAPIWay验证码", "您的验证码为："+code); err != nil {
		c.JSON(200, dto.ErrorResponse[string](400, "邮件系统发送验证码时发生错误"))
		return
	}
	c.JSON(200, dto.SuccessResponse("验证码已发送至您的邮箱，请前往查看～"))
}
func Register(c *gin.Context) {
	var registerBody RegisterRequestBody
	if err := c.ShouldBindJSON(&registerBody); err != nil {
		c.JSON(200, dto.ErrorResponse[string](400, "解构JSON时发生错误，请检查您的表单格式"))
		return
	}
	isHasExisted, err := CheckIfCodeExists(registerBody.Email)
	if err != nil {
		c.JSON(200, dto.ErrorResponse[string](400, "查询验证码时发生错误"))
		return
	}
	if !isHasExisted {
		c.JSON(200, dto.ErrorResponse[string](400, "你需要先请求您的邮箱验证码"))
		return
	}
	code, err := GetCodeFromRedis(registerBody.Email)
	if err != nil {
		c.JSON(200, dto.ErrorResponse[string](400, "查询验证码时发生错误"))
		return
	}
	if registerBody.Code != code {
		c.JSON(200, dto.ErrorResponse[string](400, "您提交的验证码与邮件验证码不符"))
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerBody.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(200, dto.ErrorResponse[string](400, "加密用户密码时发生错误，请稍后重试。"))
		return
	}

	user := pojo.User{
		Username:   registerBody.Username,
		Password:   string(hashedPassword),
		Email:      registerBody.Email,
		Tokens:     0,
		Permission: 0,
		Group:      0,
	}

	tx := config.MysqlDataBase.Begin()
	if err := tx.Where("username = ?", user.Username).First(&user).Error; err == nil {
		tx.Rollback()
		c.JSON(200, dto.ErrorResponse[string](400, "用户名已存在，更换一个吧～"))
		return
	}
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(200, dto.ErrorResponse[string](400, "创建用户时发生错误，请稍后重试。详细信息："+err.Error()))
		return
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(200, dto.ErrorResponse[string](400, "创建用户时发生错误，请稍后重试。详细信息："+err.Error()))
		return
	}
	c.JSON(200, dto.SuccessResponse("好极了！用户创建成功，欢迎您来到创剧星球!"))
	if err := DeleteCodeToRedis(registerBody.Email); err != nil {
		fmt.Println("注销验证码时发生错误:" + err.Error())
	}
}
