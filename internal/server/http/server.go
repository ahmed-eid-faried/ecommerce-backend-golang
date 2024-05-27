package http

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"

	// swaggerFiles "github.com/swaggo/files"
	// ginSwagger "github.com/swaggo/gin-swagger"

	_ "goshop/docs"
	orderHttp "goshop/internal/order/port/http"
	productHttp "goshop/internal/product/port/http"
	userHttp "goshop/internal/user/port/http"
	"goshop/pkg/config"
	"goshop/pkg/dbs"
	"goshop/pkg/redis"
	"goshop/pkg/response"
)

// Server represents the HTTP server
type Server struct {
	engine    *gin.Engine           // The underlying Gin engine
	cfg       *config.Schema        // The configuration for the server
	validator validation.Validation // The validator for request data
	db        dbs.IDatabase         // The database instance
	cache     redis.IRedis          // The Redis instance
}

// NewServer creates a new HTTP server with the provided dependencies
func NewServer(validator validation.Validation, db dbs.IDatabase, cache redis.IRedis) *Server {
	return &Server{
		engine:    gin.Default(),
		cfg:       config.GetConfig(),
		validator: validator,
		db:        db,
		cache:     cache,
	}
}

// Run starts the HTTP server
func (s Server) Run() error {

	// Set the trusted proxies for the server
	_ = s.engine.SetTrustedProxies(nil)

	// Set the mode for the server based on the environment
	if s.cfg.Environment == config.ProductionEnv {
		gin.SetMode(gin.ReleaseMode)
	}

	// Map the routes for the server
	if err := s.MapRoutes(); err != nil {
		log.Fatalf("MapRoutes Error: %v", err)
	}

	// Define the health check endpoint
	s.engine.GET("/health", func(c *gin.Context) {
		response.JSON(c, http.StatusOK, nil)
		return
	})

	// Start the HTTP server
	logger.Info("HTTP server is listening on PORT: ", s.cfg.HttpPort)
	if err := s.engine.Run(fmt.Sprintf(":%d", s.cfg.HttpPort)); err != nil {
		log.Fatalf("Running HTTP server: %v", err)
	}

	return nil
}

// GetEngine returns the underlying Gin engine
func (s Server) GetEngine() *gin.Engine {
	return s.engine
}

// MapRoutes maps the routes for the server
func (s Server) MapRoutes() error {
	// Define the routes for the server
	// s.engine.GET("/", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// s.engine.NoRoute(func(c *gin.Context) {
	// 	c.File("./templates/404.html")
	// })
	v1 := s.engine.Group("/api/v1")
	userHttp.Routes(v1, s.db, s.validator)
	productHttp.Routes(v1, s.db, s.validator, s.cache)
	orderHttp.Routes(v1, s.db, s.validator)
	return nil
}
