package web_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/confire/contest"
	"go.rtnl.ai/uptime/pkg/config"
	"go.rtnl.ai/uptime/pkg/web"
)

var env = contest.Env{
	"UPTIME_STATIC_SERVE": "true",
	"UPTIME_STATIC_ROOT":  "./static",
}

func TestRender(t *testing.T) {
	t.Cleanup(env.Set())
	t.Cleanup(config.Reset)

	renderer, err := web.NewRender(web.Templates())
	require.NoError(t, err)
	require.NotNil(t, renderer)
}

func TestStaticServe(t *testing.T) {
	staticEnv := contest.Env{
		"UPTIME_STATIC_SERVE": "true",
		"UPTIME_STATIC_ROOT":  "./static",
	}
	t.Cleanup(staticEnv.Set())
	t.Cleanup(config.Reset)

	render := &web.Render{}
	funcMap := render.FuncMap()
	f := funcMap["static"].(func(string) string)
	require.Equal(t, "/static/favicon.ico", f("favicon.ico"))
	require.Equal(t, "/static/robots.txt", f("robots.txt"))
	require.Equal(t, "/static/css/style.css", f("css/style.css"))
	require.Equal(t, "/static/js/script.js", f("js/script.js"))
	require.Equal(t, "/static/images/logo.png", f("images/logo.png"))
	require.Equal(t, "/static/fonts/font.woff", f("fonts/font.woff"))
}

func TestStaticCDN(t *testing.T) {
	cdnEnv := contest.Env{
		"UPTIME_STATIC_SERVE": "false",
		"UPTIME_STATIC_URL":   "https://cdn.example.com/",
	}
	t.Cleanup(cdnEnv.Set())
	t.Cleanup(config.Reset)

	render := &web.Render{}
	funcMap := render.FuncMap()
	f := funcMap["static"].(func(string) string)
	require.Equal(t, "https://cdn.example.com/favicon.ico", f("favicon.ico"))
	require.Equal(t, "https://cdn.example.com/robots.txt", f("robots.txt"))
	require.Equal(t, "https://cdn.example.com/css/style.css", f("css/style.css"))
	require.Equal(t, "https://cdn.example.com/js/script.js", f("js/script.js"))
	require.Equal(t, "https://cdn.example.com/images/logo.png", f("images/logo.png"))
	require.Equal(t, "https://cdn.example.com/fonts/font.woff", f("fonts/font.woff"))
}

func TestFuncMap(t *testing.T) {
	t.Cleanup(env.Set())
	t.Cleanup(config.Reset)

	renderer := &web.Render{}
	funcMap := renderer.FuncMap()

	t.Run("titlecase", func(t *testing.T) {
		testCases := []struct {
			input    string
			expected string
		}{
			{input: "hello world", expected: "Hello World"},
			{input: "Hello World", expected: "Hello World"},
			{input: "HELLO WORLD", expected: "Hello World"},
			{input: "hello world", expected: "Hello World"},
		}

		f := funcMap["titlecase"].(func(string) string)
		for _, testCase := range testCases {
			require.Equal(t, testCase.expected, f(testCase.input))
		}
	})

	t.Run("lowercase", func(t *testing.T) {
		testCases := []struct {
			input    string
			expected string
		}{
			{input: "Hello World", expected: "hello world"},
			{input: "HELLO WORLD", expected: "hello world"},
			{input: "hello world", expected: "hello world"},
		}

		f := funcMap["lowercase"].(func(string) string)
		for _, testCase := range testCases {
			require.Equal(t, testCase.expected, f(testCase.input))
		}
	})

	t.Run("uppercase", func(t *testing.T) {
		testCases := []struct {
			input    string
			expected string
		}{
			{input: "hello world", expected: "HELLO WORLD"},
			{input: "Hello World", expected: "HELLO WORLD"},
			{input: "HELLO WORLD", expected: "HELLO WORLD"},
			{input: "hello world", expected: "HELLO WORLD"},
		}

		f := funcMap["uppercase"].(func(string) string)
		for _, testCase := range testCases {
			require.Equal(t, testCase.expected, f(testCase.input))
		}
	})

	t.Run("datetime", func(t *testing.T) {
		testCases := []struct {
			input    time.Time
			expected string
		}{
			{input: time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC), expected: "January 1, 2026 12:00 PM"},
			{input: time.Date(2026, 4, 7, 9, 12, 38, 3123, time.UTC), expected: "April 7, 2026 9:12 AM"},
		}

		f := funcMap["datetime"].(func(time.Time) string)
		for _, testCase := range testCases {
			require.Equal(t, testCase.expected, f(testCase.input))
		}
	})

	t.Run("currentYear", func(t *testing.T) {
		f := funcMap["currentYear"].(func() int)
		require.Equal(t, time.Now().Year(), f())
	})

	t.Run("static", func(t *testing.T) {
		testCases := []struct {
			input    string
			expected string
		}{
			{input: "favicon.ico", expected: "/static/favicon.ico"},
			{input: "robots.txt", expected: "/static/robots.txt"},
			{input: "css/style.css", expected: "/static/css/style.css"},
			{input: "js/script.js", expected: "/static/js/script.js"},
			{input: "images/logo.png", expected: "/static/images/logo.png"},
			{input: "fonts/font.woff", expected: "/static/fonts/font.woff"},
			{input: "fonts/font.woff2", expected: "/static/fonts/font.woff2"},
			{input: "fonts/font.ttf", expected: "/static/fonts/font.ttf"},
			{input: "fonts/font.eot", expected: "/static/fonts/font.eot"},
		}

		f := funcMap["static"].(func(string) string)
		for _, testCase := range testCases {
			require.Equal(t, testCase.expected, f(testCase.input))
		}
	})
}
