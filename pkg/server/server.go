package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.rtnl.ai/uptime/pkg"
	"go.rtnl.ai/uptime/pkg/config"
	"go.rtnl.ai/x/probez"
)

const (
	ServiceName       = "uptime"
	ReadHeaderTimeout = 20 * time.Second
	WriteTimeout      = 20 * time.Second
	IdleTimeout       = 180 * time.Second
	ShutdownTimeout   = 45 * time.Second
)

var (
	tracer = otel.Tracer("go.rtnl.ai/uptime/server")
)

type Server struct {
	sync.RWMutex
	probez.Handler
	conf    config.Config
	srv     *http.Server
	router  *gin.Engine
	url     *url.URL
	started time.Time
	errc    chan error
}

func New(conf *config.Config) (s *Server, err error) {
	if conf == nil {
		if conf, err = config.New(); err != nil {
			return nil, fmt.Errorf("could not load configuration: %w", err)
		}
	}

	s = &Server{
		conf: *conf,
		errc: make(chan error, 1),
	}

	// Configure the gin router
	gin.SetMode(s.conf.Mode)
	s.router = gin.New()
	s.router.RedirectTrailingSlash = true
	s.router.RedirectFixedPath = false
	s.router.HandleMethodNotAllowed = true
	s.router.ForwardedByClientIP = true
	s.router.UseRawPath = false
	s.router.UnescapePathValues = true
	if err = s.setupRoutes(); err != nil {
		return nil, err
	}

	// Create the http server
	s.srv = &http.Server{
		Addr:              s.conf.BindAddr,
		Handler:           s.router,
		ErrorLog:          nil,
		ReadHeaderTimeout: ReadHeaderTimeout,
		WriteTimeout:      WriteTimeout,
		IdleTimeout:       IdleTimeout,
	}

	return s, nil
}

func (s *Server) Serve() (err error) {
	// Catch OS signals for graceful shutdowns
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		s.errc <- s.Shutdown()
	}()

	// Create a socket to listen on and infer the final URL.
	// NOTE: if the bindaddr is 127.0.0.1:0 for testing, a random port will be assigned,
	// manually creating the listener will allow us to determine which port.
	// When we start listening all incoming requests will be buffered until the server
	// actually starts up in its own go routine below.
	var sock net.Listener
	if sock, err = net.Listen("tcp", s.srv.Addr); err != nil {
		return fmt.Errorf("could not listen on bind addr %s: %s", s.srv.Addr, err)
	}

	s.setURL(sock.Addr())
	s.started = time.Now()
	s.Healthy()

	// Listen for HTTP requests and handle them
	go func() {
		// Make sure we don't use external err to avoid data races.
		if serr := s.serve(sock); !errors.Is(serr, http.ErrServerClosed) {
			s.errc <- serr
		}
	}()

	// Mark the server as live and ready
	s.Ready()

	// TODO: setup logging with go.rtnl.ai/x/rlog
	slog.Default().Info(
		"server started",
		slog.String("url", s.URL()),
		slog.String("version", pkg.Version(true)),
		slog.Bool("maintenance", s.conf.Maintenance),
	)
	return <-s.errc
}

func (s *Server) serve(sock net.Listener) (err error) {
	if s.srv.TLSConfig != nil {
		return s.srv.ServeTLS(sock, "", "")
	}
	return s.srv.Serve(sock)
}

func (s *Server) Shutdown() (err error) {
	slog.Default().Info("gracefully shutting down server")
	s.NotReady()

	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()

	// Shutdown the server
	s.Unhealthy()
	s.srv.SetKeepAlivesEnabled(false)
	if serr := s.srv.Shutdown(ctx); serr != nil {
		err = errors.Join(serr, fmt.Errorf("could not shutdown server: %w", serr))
	}

	slog.Default().Debug("server shutdown", slog.Any("error", err))
	return err
}

func (s *Server) URL() string {
	s.RLock()
	defer s.RUnlock()
	return s.url.String()
}

func (s *Server) setURL(addr net.Addr) {
	s.Lock()
	defer s.Unlock()
	s.url = &url.URL{
		Scheme: "http",
		Host:   addr.String(),
	}

	if s.srv.TLSConfig != nil {
		s.url.Scheme = "https"
	}

	if tcp, ok := addr.(*net.TCPAddr); ok && tcp.IP.IsUnspecified() {
		s.url.Host = fmt.Sprintf("localhost:%d", tcp.Port)
	}
}

func Span(c *gin.Context, name string, attributes ...attribute.KeyValue) (context.Context, trace.Span) {
	return tracer.Start(c.Request.Context(), name, trace.WithAttributes(attributes...))
}
