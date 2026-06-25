package app

import (
	"log"
	"net/http"
	"os"
	"personal-finance/app/auth"
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

	// Route groups
	API := e.Group("/api")

	// Rate limiter
	rateLimitConfig := getRateLimitConfig()
	API.Use(middleware.RateLimiterWithConfig(rateLimitConfig))

	// Register feature routes
	auth.RegisterRoutes(API.Group("/auth"))
	user.RegisterRoutes(API.Group("/users"))

	a := &App{
		e:      e,
		AppCfg: config.Application,
	}
	return a
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
