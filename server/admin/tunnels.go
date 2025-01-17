package admin

import (
	"errors"
	"strings"

	"github.com/andydunstall/piko/server/dbmanager"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (s *Server) RegisterTunnelsRoutes(g *gin.RouterGroup) {
	g.POST("/", s.CreateTunnel)
	g.GET("/:id", s.GetTunnel)
}

type CreateTunnelBody struct {
	Name       string `json:"name"`
	EndpointID string `json:"endpoint_id"`
}

func (s *Server) CreateTunnel(c *gin.Context) {
	var body CreateTunnelBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	tunnel, err := s.dbmanager.TunnelManager.CreateTunnel(body.Name, body.EndpointID)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			c.JSON(400, gin.H{
				"message": "Failed to create tunnel",
				"error":   "Endpoint already has a tunnel",
			})
			return
		}

		c.JSON(500, gin.H{
			"message": "Failed to create tunnel",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Tunnel created",
		"tunnel":  CreateTunnelResponse(tunnel),
	})
}

func (s *Server) GetTunnel(c *gin.Context) {
	param := c.Param("id")
	if strings.TrimSpace(param) == "" {
		c.JSON(400, gin.H{
			"message": "Missing id",
		})
		return
	}

	tunnel, err := s.dbmanager.TunnelManager.GetTunnel(param)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{
				"message": "Failed to retrieve tunnel",
				"error":   "Tunnel not found",
			})
			return
		}
		c.JSON(500, gin.H{
			"message": "Failed to retrieve tunnel",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Tunnel retrieved",
		"tunnel":  TunnelResponse(tunnel),
	})
}

func CreateTunnelResponse(tunnel *dbmanager.Tunnel) gin.H {
	return gin.H{
		"id":             tunnel.ID.String(),
		"name":           tunnel.Name,
		"endpoint_id":    tunnel.EndpointID,
		"created_at":     tunnel.CreatedAt.String(),
		"upstream_token": tunnel.UpstreamToken,
		"proxy_token":    tunnel.ProxyToken,
	}
}

func TunnelResponse(tunnel *dbmanager.Tunnel) gin.H {
	return gin.H{
		"id":          tunnel.ID.String(),
		"name":        tunnel.Name,
		"endpoint_id": tunnel.EndpointID,
		"created_at":  tunnel.CreatedAt.String(),
	}
}
