package model

import (
	"time" // استيراد حزمة الوقت

	"github.com/google/uuid" // استيراد حزمة UUID لتوليد معرفات فريدة
	"gorm.io/gorm"           // استيراد حزمة GORM لإدارة قاعدة البيانات

	"goshop/pkg/utils" // استيراد حزمة utils لتجزئة كلمة المرور
)

// تعريف نوع UserRole كنوع سلسلة (string)
type UserRole string

// تعريف القيم الممكنة لنوع UserRole
const (
	UserRoleAdmin    UserRole = "admin"    // دور المدير
	UserRoleCustomer UserRole = "customer" // دور العميل
)

// تعريف بنية المستخدم
type User struct {
	ID        string     `json:"id" gorm:"unique;not null;index;primary_key"`       // معرف المستخدم
	CreatedAt time.Time  `json:"created_at"`                                        // وقت الإنشاء
	UpdatedAt time.Time  `json:"updated_at"`                                        // وقت التحديث
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`                           // وقت الحذف (يمكن أن يكون فارغاً)
	Email     string     `json:"email" gorm:"unique;not null;index:idx_user_email"` // البريد الإلكتروني للمستخدم
	Password  string     `json:"password"`                                          // كلمة المرور
	Role      UserRole   `json:"role"`                                              // دور المستخدم
}

// دالة يتم استدعاؤها قبل إنشاء سجل المستخدم في قاعدة البيانات
func (user *User) BeforeCreate(tx *gorm.DB) error {
	// توليد معرف فريد للمستخدم
	user.ID = uuid.New().String()
	// تجزئة وتشفير كلمة المرور
	user.Password = utils.HashAndSalt([]byte(user.Password))
	// تعيين الدور الافتراضي كعميل إذا لم يتم تحديد دور
	if user.Role == "" {
		user.Role = UserRoleCustomer
	}
	return nil
}
