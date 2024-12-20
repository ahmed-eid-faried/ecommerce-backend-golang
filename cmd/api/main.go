package main

import (
	"log"
	"os"
	"time"

	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"

	orderModel "goshop/internal/order/model"
	productModel "goshop/internal/product/model"
	grpcServer "goshop/internal/server/grpc"
	httpServer "goshop/internal/server/http"
	userModel "goshop/internal/user/model"
	"goshop/pkg/config"
	"goshop/pkg/dbs"
	"goshop/pkg/redis"
)

//	@title			GoShop Swagger API
//	@version		1.0
//	@description	Swagger API for GoShop.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Quang Dang
//	@contact.email	quangdangfit@gmail.com

//	@license.name	MIT
//	@license.url	https://github.com/MartinHeinz/go-project-blueprint/blob/master/LICENSE

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization

// @BasePath	/api/v1
func main() {
	cfg := config.LoadConfig()
	logger.Initialize(cfg.Environment)
	//*********************************************

	// db, err := dbs.NewDatabase(cfg.DatabaseURI)
	// if err != nil {
	// 	logger.Fatal("Cannot connect to database", err)
	// }

	// *********************************************
	var db *dbs.Database
	var err error
	for i := 0; i < 10; i++ { // retry up to 10 times
		db, err = dbs.NewDatabase(cfg.DatabaseURI)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to database. Retrying in 5 seconds... (%d/10)\n", i+1)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		logger.Fatal("Cannot connect to database", err)
		os.Exit(1)
	}
	//*********************************************

	err = db.AutoMigrate(&userModel.User{}, &productModel.Product{}, orderModel.Order{}, orderModel.OrderLine{})
	if err != nil {
		logger.Fatal("Database migration fail", err)
	}

	validator := validation.New()

	cache := redis.New(redis.Config{
		Address:  cfg.RedisURI,
		Password: cfg.RedisPassword,
		Database: cfg.RedisDB,
	})

	go func() {
		httpSvr := httpServer.NewServer(validator, db, cache)
		if err = httpSvr.Run(); err != nil {
			logger.Fatal(err)
		}
	}()

	grpcSvr := grpcServer.NewServer(validator, db, cache)
	if err = grpcSvr.Run(); err != nil {
		logger.Fatal(err)
	}

}
