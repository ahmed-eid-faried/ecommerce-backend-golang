package grpc

import (
	"github.com/quangdangfit/gocommon/validation" // استيراد حزمة التحقق من البيانات
	"google.golang.org/grpc"                      // استيراد حزمة gRPC

	"goshop/internal/user/repository" // استيراد حزمة مستودع المستخدم
	"goshop/internal/user/service"    // استيراد حزمة خدمة المستخدم
	"goshop/pkg/dbs"                  // استيراد حزمة القواعد البيانات
	pb "goshop/proto/gen/go/user"     // استيراد حزمة البروتوكول gRPC للمستخدم
)

// RegisterHandlers يسجل المعالجين مع خادم gRPC المحدد
func RegisterHandlers(svr *grpc.Server, db dbs.IDatabase, validator validation.Validation) {
	// إنشاء مستودع المستخدم باستخدام قاعدة البيانات المعطاة
	userRepo := repository.NewUserRepository(db)
	// إنشاء خدمة المستخدم باستخدام المعالج المحدد ومستودع المستخدم
	userSvc := service.NewUserService(validator, userRepo)
	// إنشاء وحدة معالجة المستخدمين باستخدام خدمة المستخدم المحددة
	userHandler := NewUserHandler(userSvc)

	// تسجيل وحدة معالجة المستخدمين مع خادم gRPC
	pb.RegisterUserServiceServer(svr, userHandler)
}
