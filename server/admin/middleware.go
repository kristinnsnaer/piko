package admin

import (
	"github.com/gin-gonic/gin"
)

func ApiAuthMiddleware(s *Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("x-api-token")

		if !s.apiAuthConfig.Enabled() {
			c.Next()
			return
		}

		if token != s.apiAuthConfig.Token {
			c.AbortWithStatus(401)
			return
		}
	}
}
