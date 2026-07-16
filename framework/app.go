// Package framework is the Flux core: a Laravel + Inertia inspired layer on
// top of Gin, GORM and Vue. Application code (controllers, services) depends
// only on this package — never on Gin directly.
package framework

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"flux/config"
)

// App is the application container. It owns every framework-level dependency
// (config, database, logger, mailer, auth) and is the single injection point
// for services and controllers.
type App struct {
	config    *config.Config
	db        *gorm.DB
	logger    *slog.Logger
	mailer    Mailer
	auth      *AuthManager
	validator *Validator
	inertia   *Inertia
	engine    *gin.Engine
	router    *Router
}

// Option customizes App construction.
type Option func(*App)

// WithDB injects an already-opened GORM connection (used by tests and the CLI).
func WithDB(db *gorm.DB) Option {
	return func(a *App) { a.db = db }
}

// New builds the application container and the underlying HTTP engine.
func New(cfg *config.Config, opts ...Option) (*App, error) {
	app := &App{config: cfg}

	for _, opt := range opts {
		opt(app)
	}

	app.logger = NewLogger(cfg)
	slog.SetDefault(app.logger)

	app.mailer = NewMailer(cfg, app.logger)
	app.auth = NewAuthManager(cfg)
	app.validator = NewValidator()

	inertia, err := NewInertia(cfg)
	if err != nil {
		return nil, fmt.Errorf("framework: initialising inertia adapter: %w", err)
	}
	app.inertia = inertia

	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.MaxMultipartMemory = 32 << 20 // 32 MiB

	app.engine = engine
	app.router = newRouter(app, &engine.RouterGroup)

	// Static assets: the compiled Vite bundle and anything in public/.
	engine.Static("/"+cfg.Vite.BuildDir, "./public/"+cfg.Vite.BuildDir)
	engine.StaticFile("/favicon.ico", "./public/favicon.ico")
	engine.StaticFile("/favicon.svg", "./public/favicon.svg")

	engine.Use(app.requestLogger())

	return app, nil
}

// Config returns the loaded application configuration.
func (a *App) Config() *config.Config { return a.config }

// DB returns the shared GORM connection. Services receive this via
// constructor injection; the framework never hides it.
func (a *App) DB() *gorm.DB { return a.db }

// Logger returns the application slog logger.
func (a *App) Logger() *slog.Logger { return a.logger }

// Mailer returns the configured mail transport.
func (a *App) Mailer() Mailer { return a.mailer }

// Auth returns the JWT/session manager.
func (a *App) Auth() *AuthManager { return a.auth }

// Router returns the root router used to register application routes.
func (a *App) Router() *Router { return a.router }

// Handler exposes the app as an http.Handler (useful for tests).
func (a *App) Handler() http.Handler { return a.engine }

// Serve starts the HTTP server and blocks until it exits.
func (a *App) Serve() error {
	addr := fmt.Sprintf("%s:%d", a.config.App.Host, a.config.App.Port)
	a.logger.Info("flux server started",
		slog.String("addr", "http://"+addr),
		slog.String("env", a.config.App.Env),
	)

	server := &http.Server{
		Addr:              addr,
		Handler:           a.engine,
		ReadHeaderTimeout: 10 * time.Second,
	}
	return server.ListenAndServe()
}

// requestLogger logs each request through slog with method, path, status and
// latency — the framework equivalent of Laravel's access log.
func (a *App) requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		attrs := []any{
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", c.Writer.Status()),
			slog.Duration("latency", time.Since(start)),
			slog.String("ip", c.ClientIP()),
		}
		switch {
		case c.Writer.Status() >= 500:
			a.logger.Error("request", attrs...)
		case c.Writer.Status() >= 400:
			a.logger.Warn("request", attrs...)
		default:
			a.logger.Info("request", attrs...)
		}
	}
}
