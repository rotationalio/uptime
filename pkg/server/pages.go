package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.rtnl.ai/uptime/pkg/web/scene"
)

func (s *Server) Index(c *gin.Context) {
	Span(c, "ui.index")

	data := scene.New(c).ForPage("index")
	c.HTML(http.StatusOK, "pages/index.html", data)
}
