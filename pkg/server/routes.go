package server

import (
	"github.com/gin-gonic/gin"
	"go.rtnl.ai/gimlet/logger"
	"go.rtnl.ai/uptime/pkg"
)

func (s *Server) setupRoutes() (err error) {

	// Application Middleware
	middlewares := []gin.HandlerFunc{

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
