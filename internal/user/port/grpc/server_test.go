package grpc

import (
	"testing" // استيراد حزمة الاختبارات

	"github.com/quangdangfit/gocommon/validation" // استيراد حزمة التحقق من البيانات
	goGRPC "google.golang.org/grpc"               // استيراد حزمة gRPC

	"goshop/pkg/dbs/mocks" // استيراد حزمة القواعد البيانات الوهمية
)

// TestRegisterHandlers هو وظيفة الاختبار لتسجيل المعالجين
func TestRegisterHandlers(t *testing.T) {
	// إنشاء قاعدة بيانات وهمية للاختبار
	mockDB := mocks.NewIDatabase(t)
	// تسجيل المعالجين مع خادم gRPC جديد
	RegisterHandlers(goGRPC.NewServer(), mockDB, validation.New())
}
