package grpc

import (
	"context" // استيراد حزمة السياق
	"errors"  // استيراد حزمة الأخطاء
	"testing" // استيراد حزمة الاختبار

	"github.com/quangdangfit/gocommon/logger" // استيراد حزمة التسجيل
	"github.com/stretchr/testify/mock"        // استيراد حزمة محاكاة الشهادات
	"github.com/stretchr/testify/suite"       // استيراد حزمة إدارة الاختبارات

	"goshop/internal/user/dto"           // استيراد حزمة بيانات المستخدم
	"goshop/internal/user/model"         // استيراد حزمة نموذج المستخدم
	"goshop/internal/user/service/mocks" // استيراد حزمة محاكاة خدمة المستخدم
	"goshop/pkg/config"                  // استيراد حزمة التكوين
	pb "goshop/proto/gen/go/user"        // استيراد حزمة البروتوكول gRPC للمستخدم
)

// تعريف مجموعة اختبار لوحدة معالجة المستخدمين
type UserHandlerTestSuite struct {
	suite.Suite                     // تضمين بنية مجموعة الاختبارات من حزمة الشهادات
	mockService *mocks.IUserService // خدمة المستخدم الوهمية للاختبار
	handler     *UserHandler        // وحدة معالجة المستخدمين للاختبار
}

// تهيئة الاختبار
func (suite *UserHandlerTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv) // تهيئة المسجل في بيئة الإنتاج

	suite.mockService = mocks.NewIUserService(suite.T()) // إنشاء خدمة المستخدم الوهمية
	suite.handler = NewUserHandler(suite.mockService)    // إنشاء وحدة معالجة المستخدمين باستخدام الخدمة الوهمية
}

// تنفيذ مجموعة الاختبارات لوحدة معالجة المستخدمين
func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}

// اختبار تسجيل الدخول الناجح
func (suite *UserHandlerTestSuite) TestUserAPI_LoginSuccess() {
	req := &pb.LoginReq{
		Email:    "login@test.com",
		Password: "test123456",
	}

	// تحديد الاستجابة المتوقعة للخدمة الوهمية عند استدعاء دالة Login
	suite.mockService.On("Login", mock.Anything, &dto.LoginReq{
		Email:    req.Email,
		Password: req.Password,
	}).Return(
		&model.User{
			Email:    "login@test.com",
			Password: "test123456",
		},
		"access-token",
		"refresh-token",
		nil,
	).Times(1)

	// استدعاء دالة Login بوحدة المعالجة والتحقق من النتائج
	res, err := suite.handler.Login(context.Background(), req)
	suite.Nil(err)
	suite.Equal(req.Email, res.User.Email)
	suite.Equal("access-token", res.AccessToken)
	suite.Equal("refresh-token", res.RefreshToken)
}

// اختبار فشل تسجيل الدخول
func (suite *UserHandlerTestSuite) TestUserAPI_LoginFail() {
	req := &pb.LoginReq{
		Email:    "login@test.com",
		Password: "test123456",
	}

	// تحديد الاستجابة المتوقعة للخدمة الوهمية عند استدعاء دالة Login بفشل
	suite.mockService.On("Login", mock.Anything, &dto.LoginReq{
		Email:    req.Email,
		Password: req.Password,
	}).Return(nil, "", "", errors.New("error")).Times(1)

	// استدعاء دالة Login بوحدة المعالجة والتحقق من النتائج
	res, err := suite.handler.Login(context.Background(), req)
	suite.Nil(res)
	suite.NotNil(err)
}

// اختبار التسجيل الناجح
func (suite *UserHandlerTestSuite) TestUserAPI_RegisterSuccess() {
	req := &pb.RegisterReq{
		Email:    "register@test.com",
		Password: "test123456",
	}

	// تحديد الاستجابة المتوقعة للخدمة الوهمية عند استدعاء دالة Register
	suite.mockService.On("Register", mock.Anything, &dto.RegisterReq{
		Email:    req.Email,
		Password: req.Password,
	}).Return(
		&model.User{
			Email:    "register@test.com",
			Password: "test123456",
		},
		nil,
	).Times(1)

	// استدعاء دالة Register بوحدة المعالجة والتحقق من النتائج
	res, err := suite.handler.Register(context.Background(), req)
	suite.Nil(err)
	suite.Equal(req.Email, res.User.Email)
}

// اختبار فشل التسجيل
func (suite *UserHandlerTestSuite) TestUserAPI_RegisterFail() {
	req := &pb.RegisterReq{
		Email:    "register@test.com",
		Password: "test123456",
	}

	// تحديد الاستجابة المتوقعة للخدمة الوهمية عند استدعاء دالة Register بفشل
	suite.mockService.On("Register", mock.Anything, &dto.RegisterReq{
		Email:    req.Email,
		Password: req.Password,
	}).Return(nil, errors.New("error")).Times(1)

	// استدعاء دالة Register بوحدة المعالجة والتحقق من النتائج
	res, err := suite.handler.Register(context.Background(), req)
	suite.Nil(res)
	suite.NotNil(err)
}

// اختبار الحصول على معلومات المستخدم بنجاح
func (suite *UserHandlerTestSuite) TestUserAPI_GetMeSuccess() {
	userId := "123456"
	ctx := context.WithValue(context.Background(), "userId", userId)

	// تحديد الاستجابة المتوقعة للخدمة الوهمية عند استدعاء دالة GetUserByID
	suite.mockService.On("GetUserByID", mock.Anything, userId).
		Return(
			&model.User{
				ID:       userId,
				Email:    "user@test.com",
				Password: "test123456",
			},
			nil,
		).Times(1)

	// استدعاء دالة GetMe بوحدة المعالجة والتحقق من النتائج
	res, err := suite.handler.GetMe(ctx, &pb.GetMeReq{})
	suite.Nil(err)
	suite.Equal(userId, res.User.Id)
	suite.Equal("user@test.com", res.User.Email)
}

// اختبار فشل الحصول على معلومات المستخدم بسبب عدم التحقق
func (suite *UserHandlerTestSuite) TestUserAPI_GetMeUnauthorized() {
	// استدعاء دالة GetMe بوحدة المعالجة والتحقق من النتائج عند عدم وجود المستخدم في السياق
	res, err := suite.handler.GetMe(context.Background(), &pb.GetMeReq{})
	suite.Nil(res)
	suite.NotNil(err)
}

// اختبار فشل الحصول على معلومات المستخدم
func (suite *UserHandlerTestSuite) TestUserAPI_GetMeFail() {
	userId := "123456"
	ctx := context.WithValue(context.Background(), "userId", userId)

	// تحديد الاستجابة المتوقعة للخدمة الوهمية عند استدعاء دالة GetUserByID بفشل
	suite.mockService.On("GetUserByID", mock.Anything, "123456").
		Return(nil, errors.New("error")).Times(1)

	// استدعاء دالة GetMe بوحدة المعالجة والتحقق من النتائج
	res, err := suite.handler.GetMe(ctx, &pb.GetMeReq{})
	suite.Nil(res)
	suite.NotNil(err)
}

// اختبار نجاح تحديث الرمز
func (suite *UserHandlerTestSuite) TestUserAPI_RefreshTokenSuccess() {
	userId := "123456"
	ctx := context.WithValue(context.Background(), "userId", userId)

	// تحديد الاستجابة المتوقعة للخدمة الوهمية عند استدعاء دالة RefreshToken
	suite.mockService.On("RefreshToken", mock.Anything, userId).
		Return("access-token", nil).Times(1)

	// استدعاء دالة RefreshToken بوحدة المعالجة والتحقق من النتائج
	res, err := suite.handler.RefreshToken(ctx, &pb.RefreshTokenReq{})
	suite.Nil(err)
	suite.Equal("access-token", res.AccessToken)
}

// اختبار فشل تحديث الرمز بسبب عدم التحقق
func (suite *UserHandlerTestSuite) TestUserAPI_RefreshTokenUnauthorized() {
	// استدعاء دالة RefreshToken بوحدة المعالجة والتحقق من النتائج عند عدم وجود المستخدم في السياق
	res, err := suite.handler.RefreshToken(context.Background(), &pb.RefreshTokenReq{})
	suite.Nil(res)
	suite.NotNil(err)
}

// اختبار فشل تحديث الرمز
func (suite *UserHandlerTestSuite) TestUserAPI_RefreshTokenFail() {
	userId := "123456"
	ctx := context.WithValue(context.Background(), "userId", userId)

	// تحديد الاستجابة المتوقعة للخدمة الوهمية عند استدعاء دالة RefreshToken بفشل
	suite.mockService.On("RefreshToken", mock.Anything, "123456").
		Return("", errors.New("error")).Times(1)

	// استدعاء دالة RefreshToken بوحدة المعالجة والتحقق من النتائج
	res, err := suite.handler.RefreshToken(ctx, &pb.RefreshTokenReq{})
	suite.Nil(res)
	suite.NotNil(err)
}

// اختبار نجاح تغيير كلمة المرور
func (suite *UserHandlerTestSuite) TestUserAPI_ChangePasswordSuccess() {
	req := &pb.ChangePasswordReq{
		Password:    "test123456",
		NewPassword: "new-test123456",
	}

	userId := "123456"
	ctx := context.WithValue(context.Background(), "userId", userId)

	// تحديد الاستجابة المتوقعة للخدمة الوهمية عند استدعاء دالة ChangePassword
	suite.mockService.On("ChangePassword", mock.Anything, userId, &dto.ChangePasswordReq{
		Password:    req.Password,
		NewPassword: req.NewPassword,
	}).Return(nil).Times(1)

	// استدعاء دالة ChangePassword بوحدة المعالجة والتحقق من النتائج
	res, err := suite.handler.ChangePassword(ctx, req)
	suite.NotNil(res)
	suite.Nil(err)
}

// اختبار فشل تغيير كلمة المرور بسبب عدم التحقق
func (suite *UserHandlerTestSuite) TestUserAPI_ChangePasswordUnauthorized() {
	// استدعاء دالة ChangePassword بوحدة المعالجة والتحقق من النتائج عند عدم وجود المستخدم في السياق
	res, err := suite.handler.ChangePassword(context.Background(), &pb.ChangePasswordReq{})
	suite.Nil(res)
	suite.NotNil(err)
}

// اختبار فشل تغيير كلمة المرور
func (suite *UserHandlerTestSuite) TestUserAPI_ChangePasswordFail() {
	req := &pb.ChangePasswordReq{
		Password:    "test123456",
		NewPassword: "new-test123456",
	}

	userId := "123456"
	ctx := context.WithValue(context.Background(), "userId", userId)

	// تحديد الاستجابة المتوقعة للخدمة الوهمية عند استدعاء دالة ChangePassword بفشل
	suite.mockService.On("ChangePassword", mock.Anything, "123456", &dto.ChangePasswordReq{
		Password:    req.Password,
		NewPassword: req.NewPassword,
	}).Return(errors.New("error")).Times(1)

	// استدعاء دالة ChangePassword بوحدة المعالجة والتحقق من النتائج
	res, err := suite.handler.ChangePassword(ctx, req)
	suite.Nil(res)
	suite.NotNil(err)
}
