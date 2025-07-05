package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"onepenny-server/util"
	"strings"
)

// AuthMiddleware 从 Authorization: Bearer <token> 验证 JWT 并写入 ctx
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		hdr := c.GetHeader("Authorization")
		if hdr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			return
		}

		parts := strings.SplitN(hdr, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			return
		}

		token := parts[1]
		userID, err := util.ValidateJWT(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		// 写入 userID，到后续 handler 可通过 c.Get("userID") 拿到
		c.Set("userID", userID)
		c.Next()
	}
}
