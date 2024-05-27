package dto

import (
	"time" // استيراد حزمة الوقت

	_ "goshop/docs" // هذا لتحميل وثائق Swagger
)

// تعريف بنية المستخدم
type User struct {
	ID        string    `json:"id"`         // معرف المستخدم
	Email     string    `json:"email"`      // البريد الإلكتروني للمستخدم
	CreatedAt time.Time `json:"created_at"` // وقت الإنشاء
	UpdatedAt time.Time `json:"updated_at"` // وقت التحديث
}

// تعريف بنية طلب التسجيل
type RegisterReq struct {
	Email    string `json:"email" validate:"required,email"`       // البريد الإلكتروني مع التحقق من البريد الإلكتروني المطلوب
	Password string `json:"password" validate:"required,password"` // كلمة المرور مع التحقق من كلمة المرور المطلوبة
}

// تعريف بنية استجابة التسجيل
type RegisterRes struct {
	User User `json:"user"` // المستخدم المسجل
}

// تعريف بنية طلب تسجيل الدخول
type LoginReq struct {
	Email    string `json:"email" validate:"required,email"`       // البريد الإلكتروني مع التحقق من البريد الإلكتروني المطلوب
	Password string `json:"password" validate:"required,password"` // كلمة المرور مع التحقق من كلمة المرور المطلوبة
}

// تعريف بنية استجابة تسجيل الدخول
type LoginRes struct {
	User         User   `json:"user"`          // المستخدم الذي قام بتسجيل الدخول
	AccessToken  string `json:"access_token"`  // رمز الوصول
	RefreshToken string `json:"refresh_token"` // رمز التحديث
}

// تعريف بنية طلب تحديث الرمز
type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token" validate:"required"` // رمز التحديث مع التحقق المطلوب
}

// تعريف بنية استجابة تحديث الرمز
type RefreshTokenRes struct {
	AccessToken string `json:"access_token"` // رمز الوصول الجديد
}

// تعريف بنية طلب تغيير كلمة المرور
type ChangePasswordReq struct {
	Password    string `json:"password" validate:"required,password"`     // كلمة المرور الحالية مع التحقق من كلمة المرور المطلوبة
	NewPassword string `json:"new_password" validate:"required,password"` // كلمة المرور الجديدة مع التحقق من كلمة المرور المطلوبة
}
