// Copyright 2024 Andrew Dunstall. All rights reserved.
//
// Use of this source code is governed by a MIT style license that can be
// found in the LICENSE file.

package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/andydunstall/pico/pkg/conn"
	"github.com/andydunstall/pico/pkg/log"
	"github.com/andydunstall/pico/pkg/rpc"
	"github.com/andydunstall/pico/server/config"
	"github.com/andydunstall/pico/server/middleware"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// Server is the HTTP server used for both upstream listeners and downstream
// clients.
//
// /pico is reserved for upstream listeners and management, then all other
// routes will be proxied.
type Server struct {
	httpServer        *http.Server
	rpcServer         *rpcServer
	router            *gin.Engine
	websocketUpgrader *websocket.Upgrader

	shutdownCtx    context.Context
	shutdownCancel func()

	addr string

	registry *prometheus.Registry

	conf   *config.Config
	logger *log.Logger
}

func NewServer(
	addr string,
	registry *prometheus.Registry,
	conf *config.Config,
	logger *log.Logger,
) *Server {
	router := gin.New()
	// Recover from panics.
	router.Use(gin.Recovery())
	router.Use(middleware.NewLogger(logger))

	if registry != nil {
		router.Use(middleware.NewMetrics(registry))
	}

	shutdownCtx, shutdownCancel := context.WithCancel(context.Background())

	s := &Server{
		router: router,
		httpServer: &http.Server{
			Addr:    addr,
			Handler: router,
		},
		rpcServer:         newRPCServer(),
		websocketUpgrader: &websocket.Upgrader{},

		shutdownCtx:    shutdownCtx,
		shutdownCancel: shutdownCancel,

		addr:     addr,
		registry: registry,
		conf:     conf,
		logger:   logger.WithSubsystem("server.http"),
	}
	s.registerRoutes()
	return s
}

func (s *Server) Serve() error {
	s.logger.Info("starting http server", zap.String("addr", s.addr))

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("http serve: %w", err)
	}
	return nil
}

// Shutdown attempts to gracefully shutdown the server by closing open
// WebSockets and waiting for pending requests to complete.
func (s *Server) Shutdown(ctx context.Context) error {
	// Shutdown listeners.
	s.shutdownCancel()

	return s.httpServer.Shutdown(ctx)
}

func (s *Server) registerRoutes() {
	pico := s.router.Group("/pico/v1")
	pico.GET("/listener/:endpointID", s.listener)
	pico.GET("/health", s.health)

	if s.registry != nil {
		pico.GET("/metrics", s.metricsHandler())
	}

	// Handle not found routes, which includes all proxied endpoints.
	s.router.NoRoute(s.notFound)
}

func (s *Server) notFound(c *gin.Context) {
	// All /pico endpoints are reserved. All others are proxied.
	if strings.HasPrefix(c.Request.URL.Path, "/pico") {
		c.Status(http.StatusNotFound)
		return
	}
	s.proxy(c)
}

// proxy handles proxied requests from downstream clients.
func (s *Server) proxy(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}

// listener handles WebSocket connections from upstream listeners.
func (s *Server) listener(c *gin.Context) {
	endpointID := c.Param("endpointID")

	wsConn, err := s.websocketUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		// Upgrade replies to the client so nothing else to do.
		s.logger.Warn("failed to upgrade websocket", zap.Error(err))
		return
	}
	stream := rpc.NewStream(
		conn.NewWebsocketConn(wsConn),
		s.rpcServer.Handler(),
		s.logger,
	)
	defer stream.Close()

	s.logger.Debug(
		"upstream listener connected",
		zap.String("endpoint-id", endpointID),
	)

	listener := newListener(endpointID, stream, s.conf.Upstream, s.logger)
	if err := listener.Monitor(s.shutdownCtx); err != nil {
		s.logger.Warn("listener unexpectly disconnected", zap.Error(err))
	}
}

func (s *Server) health(c *gin.Context) {
	c.Status(http.StatusOK)
}

func (s *Server) metricsHandler() gin.HandlerFunc {
	h := promhttp.HandlerFor(
		s.registry,
		promhttp.HandlerOpts{Registry: s.registry},
	)
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func init() {
	// Disable Gin debugging.
	gin.SetMode(gin.ReleaseMode)
}