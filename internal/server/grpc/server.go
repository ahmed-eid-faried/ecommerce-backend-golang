package grpc

import (
	"fmt" // استيراد حزمة تنسيق النص
	"net" // استيراد حزمة الشبكات

	"github.com/quangdangfit/gocommon/logger"     // استيراد حزمة التسجيل
	"github.com/quangdangfit/gocommon/validation" // استيراد حزمة التحقق
	"google.golang.org/grpc"                      // استيراد حزمة gRPC
	"google.golang.org/grpc/reflection"           // استيراد حزمة انعكاس gRPC

	cartGRPC "goshop/internal/cart/port/grpc" // استيراد حزمة gRPC لعربة التسوق
	userGRPC "goshop/internal/user/port/grpc" // استيراد حزمة gRPC للمستخدم
	"goshop/pkg/config"                       // استيراد حزمة التكوين
	"goshop/pkg/dbs"                          // استيراد حزمة قاعدة البيانات
	"goshop/pkg/middleware"                   // استيراد حزمة البرامج الوسيطة
	"goshop/pkg/redis"                        // استيراد حزمة Redis
)

// تعريف هيكل الخادم
type Server struct {
	engine    *grpc.Server          // محرك gRPC
	cfg       *config.Schema        // التكوين
	validator validation.Validation // المدقق
	db        dbs.IDatabase         // قاعدة البيانات
	cache     redis.IRedis          // ذاكرة التخزين المؤقت
}

// دالة إنشاء خادم جديد
func NewServer(validator validation.Validation, db dbs.IDatabase, cache redis.IRedis) *Server {
	// إنشاء معترض التحقق من الصلاحيات
	interceptor := middleware.NewAuthInterceptor(config.AuthIgnoreMethods)

	// إنشاء محرك gRPC مع معترض التحقق من الصلاحيات
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.Unary(),
		),
	)

	// إرجاع خادم جديد
	return &Server{
		engine:    grpcServer,
		cfg:       config.GetConfig(),
		validator: validator,
		db:        db,
		cache:     cache,
	}
}

// دالة تشغيل الخادم
func (s Server) Run() error {

	// تسجيل معالجات gRPC للمستخدم وعربة التسوق
	userGRPC.RegisterHandlers(s.engine, s.db, s.validator)
	cartGRPC.RegisterHandlers(s.engine, s.db, s.validator)

	// تسجيل انعكاس gRPC
	reflection.Register(s.engine)

	// الاستماع على المنفذ المحدد في التكوين
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.GrpcPort))
	logger.Info("GRPC server is listening on PORT: ", s.cfg.GrpcPort)
	if err != nil {
		logger.Error("Failed to listen: ", err)
		return err
	}

	// بدء تشغيل خادم gRPC
	err = s.engine.Serve(lis)
	if err != nil {
		logger.Fatal("Failed to serve grpc: ", err)
		return err
	}

	return nil
}
