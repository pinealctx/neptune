package mock

import (
	"github.com/pinealctx/neptune/ulog"
	"go.uber.org/zap"
)

type Mock struct {
}

func NewMockSMS() *Mock {
	return &Mock{}
}

func (m *Mock) SendCode(areaCode, phone, code string) error {
	ulog.Debug("Mock.SendCode", zap.String("areaCode", areaCode),
		zap.String("phone", phone), zap.String("code", code))
	return nil
}
