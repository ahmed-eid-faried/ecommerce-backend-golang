package main

import (
	"github.com/quangdangfit/gocommon/logger"     // استيراد حزمة تسجيل الدخول
	"github.com/quangdangfit/gocommon/validation" // استيراد حزمة التحقق

	_ "goshop/docs"                              // This is for loading the swag docs 	// هذا لتحميل وثائق Swagger
	orderModel "goshop/internal/order/model"     // استيراد نموذج الطلبات
	productModel "goshop/internal/product/model" // استيراد نموذج المنتجات
	grpcServer "goshop/internal/server/grpc"     // استيراد خادم gRPC
	httpServer "goshop/internal/server/http"     // استيراد خادم HTTP
	userModel "goshop/internal/user/model"       // استيراد نموذج المستخدمين
	"goshop/pkg/config"                          // استيراد حزمة التكوين
	"goshop/pkg/dbs"                             // استيراد حزمة قاعدة البيانات
	"goshop/pkg/redis"                           // استيراد حزمة Redis
)

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @termsOfService	http://swagger.io/terms/
// @contact.name	Quang Dang
// @contact.email	quangdangfit@gmail.com
// @license.name	MIT
// @license.url	https://github.com/MartinHeinz/go-project-blueprint/blob/master/LICENSE

// @Title			GoShop Swagger API
// @version		1.0
// @Description	Swagger API for GoShop.
// @host localhost:8888
// @BasePath /
func main() {

	// تحميل التكوين
	cfg := config.LoadConfig()
	// تهيئة نظام التسجيل
	logger.Initialize(cfg.Environment)

	// إنشاء اتصال قاعدة البيانات
	db, err := dbs.NewDatabase(cfg.DatabaseURI)
	if err != nil {
		// تسجيل خطأ إذا لم يتم الاتصال بقاعدة البيانات
		logger.Fatal("Cannot connect to database", err)
	}

	// ترحيل قاعدة البيانات
	err = db.AutoMigrate(&userModel.User{}, &productModel.Product{}, orderModel.Order{}, orderModel.OrderLine{})
	if err != nil {
		// تسجيل خطأ إذا فشل الترحيل
		logger.Fatal("Database migration fail", err)
	}

	// إنشاء المدقق
	validator := validation.New()

	// إعداد ذاكرة التخزين المؤقت Redis
	cache := redis.New(redis.Config{
		Address:  cfg.RedisURI,      // عنوان Redis
		Password: cfg.RedisPassword, // كلمة مرور Redis
		Database: cfg.RedisDB,       // قاعدة بيانات Redis
	})

	// تشغيل خادم HTTP في روتين مستقل
	go func() {
		httpSvr := httpServer.NewServer(validator, db, cache)
		if err = httpSvr.Run(); err != nil {
			// تسجيل خطأ إذا فشل تشغيل خادم HTTP
			logger.Fatal(err)
		}
	}()

	// إعداد وتشغيل خادم gRPC
	grpcSvr := grpcServer.NewServer(validator, db, cache)
	if err = grpcSvr.Run(); err != nil {
		// تسجيل خطأ إذا فشل تشغيل خادم gRPC
		logger.Fatal(err)
	}
}
