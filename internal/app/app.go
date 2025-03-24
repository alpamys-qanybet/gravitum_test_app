package app

import (
	"context"
	"fmt"
	"gravitum-test-app/config"
	"gravitum-test-app/internal/handler"
	"gravitum-test-app/internal/repository/postgres"
	"gravitum-test-app/internal/service"
	"gravitum-test-app/pkg/logger"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

type App struct {
	cfg    *config.Config
	log    *logger.Logger
	Db     *pgxpool.Pool
	Server *http.Server
}

func New(
	cfg *config.Config,
	log *logger.Logger,
) *App {
	return &App{
		cfg: cfg,
		log: log,
	}
}

func (app *App) Run(ctx context.Context) error {

	err := app.ConnectDB(ctx, app.cfg.GetDbConfig().GetDsn())
	if err != nil {
		app.log.Error(fmt.Sprintf("couldn't instantiate db: %s", err))
		return err
	}
	defer app.Db.Close()

	repo := postgres.NewRepository(app.cfg, app.Db)

	service := service.NewService(
		app.cfg,
		repo,
	)

	handler := handler.NewHandler(app.cfg, service, app.log)

	r := gin.New()
	r.Use(gin.Recovery()) // recovery middleware
	r.Use(secure.New(secure.Config{
		BrowserXssFilter:   true,
		ContentTypeNosniff: true,
		ReferrerPolicy:     "no-referrer",
	}))
	app.setupRouter(r, handler)

	app.Server = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", app.cfg.App.Host, app.cfg.App.Port),
		Handler: r,
	}

	return app.Server.ListenAndServe()
}

func (app *App) ConnectDB(ctx context.Context, dsn string) error {
	dbpool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return err
	}

	dbpool.Ping(ctx)
	app.Db = dbpool
	return nil
}

func (app *App) setupRouter(r *gin.Engine, h *handler.Handler) {

	api := r.Group("/api")

	if app.cfg.Security.CorsEnabled {
		if len(app.cfg.Security.CorsAllowOrigins) > 0 {
			allowOrigins := strings.Split(app.cfg.Security.CorsAllowOrigins, ",")
			if len(allowOrigins) > 0 {
				corsMiddleware := cors.New(cors.Config{
					AllowOrigins:     allowOrigins,
					AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
					AllowHeaders:     []string{"Origin", "Content-Type", "Content-Language", "Accept", "Authorization", "X-API-SECRET-KEY"},
					ExposeHeaders:    []string{"Content-Length", "Authorization"},
					AllowCredentials: true,
					MaxAge:           12 * time.Hour,
				})

				api.Use(corsMiddleware)

				// Add OPTIONS handler to users group for preflight requests
				api.OPTIONS("/*path", func(c *gin.Context) {
					c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
					c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
					c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-API-SECRET-KEY, Content-Language")
					c.Header("Access-Control-Allow-Credentials", "true")
					c.Status(204) // No Content
				})
			}
		}
	}

	users := api.Group("/users")

	// user routes
	users.GET("/", h.User.GetList)   // api - get user list
	users.GET("/:id", h.User.Get)    // api - get user
	users.POST("/", h.User.Create)   // api - create user
	users.PUT("/:id", h.User.Update) // api - update user method

	// Public API root
	r.GET("/api", func(c *gin.Context) { // api - root(it works)
		c.String(http.StatusOK, "it works")
	})
}
