package server

import (
	"github.com/gin-gonic/gin"
	"go.rtnl.ai/gimlet/logger"
	"go.rtnl.ai/uptime/pkg"
	"go.rtnl.ai/uptime/pkg/telemetry"
)

func (s *Server) setupRoutes() (err error) {
	// Create observability middleware
	var observability gin.HandlerFunc
	if observability, err = telemetry.Middleware(); err != nil {
		return err
	}

	// Application Middleware
	middlewares := []gin.HandlerFunc{
		// o11y should be on the outside so we can record the correct latency of requests
		// NOTE: o11y panics will not recover due to middleware ordering.
		observability,

		// Panic recovery middleware
		gin.Recovery(),

		// Optional logging middleware
		logger.Logger(ServiceName, pkg.Version(true)),

		// Return unavailable if the server is in maintenance mode
		s.Maintenance(),
	}

	// Kubernetes liveness probes added before middleware.
	s.router.GET("/healthz", gin.WrapF(s.Healthz))
	s.router.GET("/livez", gin.WrapF(s.Healthz))
	s.router.GET("/readyz", gin.WrapF(s.Readyz))

	// Add the middleware to the router
	for _, middleware := range middlewares {
		if middleware != nil {
			s.router.Use(middleware)
		}
	}

	// Not Found and Not Allowed routes
	s.router.NoRoute(s.NotFound)
	s.router.NoMethod(s.NotAllowed)

	// Error Pages for Front-End Error Redirection
	s.router.GET("/not-found", s.NotFound)
	s.router.GET("/not-allowed", s.NotAllowed)
	s.router.GET("/error", s.InternalError)

	// Unauthenticated API Routes
	v1o := s.router.Group("/v1")
	{
		// Status/Heartbeat endpoint
		v1o.GET("/status", s.Status)
	}

	return nil
}
