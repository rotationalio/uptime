package web_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/confire/contest"
	"go.rtnl.ai/uptime/pkg/web"
)

var env = contest.Env{
	"UPTIME_STATIC_SERVE": "true",
	"UPTIME_STATIC_ROOT":  "./static",
}

func TestRender(t *testing.T) {
	t.Cleanup(env.Set())
	renderer, err := web.NewRender(web.Templates())
	require.NoError(t, err)
	require.NotNil(t, renderer)
}

func TestFuncMap(t *testing.T) {
	t.Cleanup(env.Set())
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
