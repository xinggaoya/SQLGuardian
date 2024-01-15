package middleware

import "github.com/gin-gonic/gin"

// Cors 跨域
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置允许跨域的域名，这里设置为允许所有来源
		c.Header("Access-Control-Allow-Origin", "*")
		// 设置允许的请求方法
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// 设置允许的请求头
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length")
		// 设置是否允许发送凭证信息（例如，使用 cookies）
		c.Header("Access-Control-Allow-Credentials", "true")
		// 设置预检请求的有效期，单位为秒
		c.Header("Access-Control-Max-Age", "12")

		// 继续处理请求
		c.Next()
	}
}
