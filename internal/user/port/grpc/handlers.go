package grpc

import (
	"context" // استيراد حزمة السياق
	"errors"  // استيراد حزمة الأخطاء

	"github.com/quangdangfit/gocommon/logger" // استيراد حزمة التسجيل

	"goshop/internal/user/dto"     // استيراد حزمة بيانات المستخدم
	"goshop/internal/user/service" // استيراد حزمة خدمة المستخدم
	"goshop/pkg/utils"             // استيراد حزمة الأدوات المساعدة
	pb "goshop/proto/gen/go/user"  // استيراد حزمة البروتوكول gRPC للمستخدم
)

// تعريف وحدة معالجة المستخدمين
type UserHandler struct {
	pb.UnimplementedUserServiceServer // تضمين السيرفر غير المنفذ من gRPC

	service service.IUserService // خدمة المستخدم
}

// دالة إنشاء وحدة معالجة المستخدمين الجديدة
func NewUserHandler(service service.IUserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// دالة تسجيل الدخول
func (h *UserHandler) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginRes, error) {
	// استدعاء خدمة تسجيل الدخول والحصول على النتائج
	user, accessToken, refreshToken, err := h.service.Login(ctx, &dto.LoginReq{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		logger.Error("Failed to login ", err) // تسجيل الخطأ في حالة الفشل
		return nil, err
	}

	var res pb.LoginRes
	utils.Copy(&res.User, &user)    // نسخ بيانات المستخدم من النتيجة
	res.AccessToken = accessToken   // تعيين الرمز المميز للوصول
	res.RefreshToken = refreshToken // تعيين الرمز المميز للتحديث
	return &res, nil
}

// دالة التسجيل
func (h *UserHandler) Register(ctx context.Context, req *pb.RegisterReq) (*pb.RegisterRes, error) {
	// استدعاء خدمة التسجيل والحصول على النتائج
	user, err := h.service.Register(ctx, &dto.RegisterReq{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		logger.Error("Failed to register ", err) // تسجيل الخطأ في حالة الفشل
		return nil, err
	}

	var res pb.RegisterRes
	utils.Copy(&res.User, &user) // نسخ بيانات المستخدم من النتيجة
	return &res, nil
}

// دالة الحصول على معلومات المستخدم
func (h *UserHandler) GetMe(ctx context.Context, _ *pb.GetMeReq) (*pb.GetMeRes, error) {
	userID, _ := ctx.Value("userId").(string) // الحصول على معرف المستخدم من السياق
	if userID == "" {
		return nil, errors.New("unauthorized") // إرجاع خطأ في حالة عدم التحقق
	}

	// استدعاء خدمة الحصول على المستخدم بواسطة معرف المستخدم
	user, err := h.service.GetUserByID(ctx, userID)
	if err != nil {
		logger.Error("Failed to get user info ", err) // تسجيل الخطأ في حالة الفشل
		return nil, err
	}

	var res pb.GetMeRes
	utils.Copy(&res.User, &user) // نسخ بيانات المستخدم من النتيجة
	return &res, nil
}

// دالة تحديث الرمز المميز
func (h *UserHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenReq) (*pb.RefreshTokenRes, error) {
	userID, _ := ctx.Value("userId").(string) // الحصول على معرف المستخدم من السياق
	if userID == "" {
		return nil, errors.New("unauthorized") // إرجاع خطأ في حالة عدم التحقق
	}

	// استدعاء خدمة تحديث الرمز المميز
	accessToken, err := h.service.RefreshToken(ctx, userID)
	if err != nil {
		logger.Error("Failed to refresh token ", err) // تسجيل الخطأ في حالة الفشل
		return nil, err
	}

	res := pb.RefreshTokenRes{
		AccessToken: accessToken, // تعيين الرمز المميز للوصول
	}
	return &res, nil
}

// دالة تغيير كلمة المرور
func (h *UserHandler) ChangePassword(ctx context.Context, req *pb.ChangePasswordReq) (*pb.ChangePasswordRes, error) {
	userID, _ := ctx.Value("userId").(string) // الحصول على معرف المستخدم من السياق
	if userID == "" {
		return nil, errors.New("unauthorized") // إرجاع خطأ في حالة عدم التحقق
	}

	// استدعاء خدمة تغيير كلمة المرور
	err := h.service.ChangePassword(ctx, userID, &dto.ChangePasswordReq{
		Password:    req.Password,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		logger.Error("Failed to change password ", err) // تسجيل الخطأ في حالة الفشل
		return nil, err
	}

	return &pb.ChangePasswordRes{}, nil
}
