package admin

import "github.com/gin-gonic/gin"

func (s *Server) RegisterAPIRoutes(g *gin.RouterGroup) {
	s.RegisterTunnelsRoutes(g.Group("/tunnels"))
}
