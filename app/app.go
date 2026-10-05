package app

import (
	"log"
	"net/http"
	"os"
	"personal-finance/app/account"
	"personal-finance/app/auth"
	"personal-finance/app/budget"
	"personal-finance/app/category"
	"personal-finance/app/report"
	"personal-finance/app/transaction"
	"personal-finance/app/user"
	"personal-finance/config"
	"personal-finance/utils"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

// App holds the Echo instance and application config
type App struct {
	e      *echo.Echo
	AppCfg *config.App
}

// New initializes a new App: loads config, registers routes
func New() *App {
	e := echo.New()

	config.Application = &config.App{}
	if err := config.Application.InitConfig(); err != nil {
		panic(err)
	}

	// Middleware
	e.Use(middleware.CORSWithConfig(utils.GetMiddleWareConfig()))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Handle OPTIONS pre-flight
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Method == http.MethodOptions {
				return c.NoContent(http.StatusNoContent)
			}
			return next(c)
		}
	})

	registerRoutes(e)

	a := &App{
		e:      e,
		AppCfg: config.Application,
	}
	return a
}

// registerRoutes builds the whole route table. It is separate from New so the
// public/protected split can be tested without a database or a live config.
func registerRoutes(e *echo.Echo) {
	API := e.Group("/api")

	// Rate limiter
	API.Use(middleware.RateLimiterWithConfig(getRateLimitConfig()))

	// Public routes
	auth.RegisterRoutes(API.Group("/auth"))

	// Health check, useful for container orchestration
	API.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, utils.ResponseOK(map[string]string{
			"status":  "ok",
			"version": config.Application.Version,
		}))
	})

	// Protected routes. Every handler below resolves the caller from the JWT
	// via utils.GetUserID, so data is always scoped to the signed-in user.
	protected := API.Group("", utils.JWTMiddleware())
	user.RegisterRoutes(protected.Group("/users"))
	account.RegisterRoutes(protected.Group("/accounts"))
	category.RegisterRoutes(protected.Group("/categories"))
	transaction.RegisterRoutes(protected.Group("/transactions"))
	budget.RegisterRoutes(protected.Group("/budgets"))
	report.RegisterRoutes(protected.Group("/reports"))
}

// Start starts the HTTP server and cron scheduler
func (a *App) Start(addr string) error {
	a.startCron()
	return a.e.Start(addr)
}

// startCron starts background cron jobs (if enabled)
func (a *App) startCron() {
	if !a.AppCfg.CronEnabled {
		log.Println("⏰ Cron scheduler is disabled")
		return
	}

	log.Println("⏰ Starting cron scheduler...")
	go func() {
		ticker := time.NewTicker(time.Second * 30)
		defer ticker.Stop()

		for {
			<-ticker.C
			// TODO: add your cron jobs here
			// a.SomeCronJob()
			log.Println("⏰ Cron tick")
		}
	}()
}

// getRateLimitConfig configures rate limiting per IP
func getRateLimitConfig() middleware.RateLimiterConfig {
	rateVal := 20.0
	if os.Getenv("ENVIRONMENT") == "local" {
		rateVal = 100.0
	}

	return middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{
				Rate:      rate.Limit(rateVal),
				Burst:     int(rateVal),
				ExpiresIn: 5 * time.Minute,
			},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			return ctx.RealIP(), nil
		},
		ErrorHandler: func(ctx echo.Context, err error) error {
			return ctx.JSON(http.StatusTooManyRequests, map[string]interface{}{
				"status":  "error",
				"message": "Too many requests. Please slow down.",
			})
		},
	}
}
