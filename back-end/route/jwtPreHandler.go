package route

import (
	"back-end/entity/dto"
	"back-end/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

/**
 * @author Flyinsky
 * @email w2084151024@gmail.com
 * @date 2024/12/24 11:34
 */
func CorsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:5173") // 或者指定具体的域名，如 http://example.com
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true") // 是否允许携带 cookie

		// 对于预检请求，直接返回状态码 204
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func preHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 前置处理：比如检查Token，记录日志等
		if c.Request.RequestURI == "/api/auth/login" || c.Request.RequestURI == "/api/auth/register" {
			c.Next()
			return
		}
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(200, dto.ErrorResponse[string](401, "未提供令牌"))
			c.Abort()
			return
		}
		claims, err := util.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "token不合法或已过期，请尝试重新登陆获取。"))
			c.Abort()
			return
		}
		c.Set("userId", claims.UserID)

		// 允许请求继续进行
		c.Next()

		//statusCode := c.Writer.Status()
		//c.Writer.Header().Set("X-Response-Time", "123ms")
		//if statusCode != 200 {
		//	fmt.Printf("Request failed with status: %d\n", statusCode)
		//}
	}
}
