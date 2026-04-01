/*
Scene provides well structured template contexts and functionality for HTML template
rendering. We chose the word "scene" to represent the context since "context" is an
overloaded term and milieu was too hard to spell.
*/
package scene

import (
	"maps"

	"github.com/gin-gonic/gin"
	"go.rtnl.ai/uptime/pkg"
)

var (
	// Compute the version of the package at runtime so it is static for all contexts.
	version      = pkg.Version(false)
	shortVersion = pkg.Version(true)
	revision     = pkg.GitVersion
	buildDate    = pkg.BuildDate
)

// Keys for default Scene context items
const (
	Version      = "Version"
	ShortVersion = "ShortVersion"
	Revision     = "Revision"
	BuildDate    = "BuildDate"
	Page         = "Page"
	APIData      = "APIData"
	Path         = "Path"
)

type Scene map[string]any

func New(c *gin.Context) Scene {
	if c == nil {
		return Scene{
			Version:      version,
			ShortVersion: shortVersion,
			Revision:     revision,
			BuildDate:    buildDate,
		}
	}

	// Create the basic context
	context := Scene{
		Version:      version,
		ShortVersion: shortVersion,
		Revision:     revision,
		BuildDate:    buildDate,
	}

	// Set path from the gin context
	context[Path] = c.Request.URL.Path
	return context
}

func (s Scene) Update(o Scene) Scene {
	maps.Copy(s, o)
	return s
}

func (s Scene) WithAPIData(data any) Scene {
	s[APIData] = data
	return s
}

func (s Scene) ForPage(page string) Scene {
	s[Page] = page
	return s
}
