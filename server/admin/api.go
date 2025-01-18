package admin

import "github.com/gin-gonic/gin"

func (s *Server) RegisterAPIRoutes(g *gin.RouterGroup) {
	g.Use(ApiAuthMiddleware(s))
	s.RegisterTunnelsRoutes(g.Group("/tunnels"))
}
