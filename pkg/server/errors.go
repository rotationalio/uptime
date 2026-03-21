package server

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.rtnl.ai/uptime/pkg/api/v1"
	"go.rtnl.ai/uptime/pkg/web/scene"
)

var (
	ErrNotFound   = errors.New("resource not found")
	ErrNotAllowed = errors.New("method not allowed")
)

// Log the error and return an appropriate response to the client. If HTML is requested
// in the Accept header then a 500 error page is returned. If JSON is requested then
// the error is rendered as a JSON response.
func (s *Server) Error(c *gin.Context, err error) {
	if err != nil {
		c.Error(err)
	}

	c.Negotiate(http.StatusInternalServerError, gin.Negotiate{
		Offered:  []string{binding.MIMEJSON, binding.MIMEHTML},
		HTMLName: "errors/status/500.html",
		HTMLData: scene.New(c).Error(err),
		JSONData: api.Error(err),
	})
}

// Render the not found error page
func (s *Server) NotFound(c *gin.Context) {
	c.Negotiate(http.StatusNotFound, gin.Negotiate{
		Offered:  []string{binding.MIMEJSON, binding.MIMEHTML},
		HTMLName: "errors/status/404.html",
		HTMLData: scene.New(c).Error(ErrNotFound),
		JSONData: api.NotFound,
	})
}

// Render the not allowed error page
func (s *Server) NotAllowed(c *gin.Context) {
	c.Negotiate(http.StatusMethodNotAllowed, gin.Negotiate{
		Offered:  []string{binding.MIMEJSON, binding.MIMEHTML},
		HTMLName: "errors/status/405.html",
		HTMLData: scene.New(c).Error(ErrNotAllowed),
		JSONData: api.NotAllowed,
	})
}

// Render the internal server error page
func (s *Server) InternalError(c *gin.Context) {
	c.Negotiate(http.StatusInternalServerError, gin.Negotiate{
		Offered:  []string{binding.MIMEJSON, binding.MIMEHTML},
		HTMLName: "errors/status/500.html",
		HTMLData: scene.New(c).Error(nil),
		JSONData: api.InternalError,
	})
}
