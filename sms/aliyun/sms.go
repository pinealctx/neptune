package aliyun

import (
	"errors"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/pinealctx/neptune/cryptx"
	"github.com/pinealctx/neptune/ulog"
	"go.uber.org/zap"
	"strings"
)

type Config struct {
	Access      string `json:"access" toml:"access"`
	Secret      string `json:"secret" toml:"secret"`
	CNTmpl      string `json:"cn_tmpl" toml:"cn_tmpl"`
	OverseaTmpl string `json:"oversea_tmpl" toml:"oversea_tmpl"`
	SignName    string `json:"sign_name" toml:"sign_name"`
}

func (c *Config) Decrypt() error {
	var err error
	c.Secret, err = cryptx.DecryptSenInfo(c.Secret)
	return err
}

type SMS struct {
	client      *dysmsapi.Client
	cnTmpl      string
	overseaTmpl string
	signName    string
}

func NewAliSMS(cnf *Config) (*SMS, error) {
	var client, err = dysmsapi.NewClientWithAccessKey("cn-hangzhou", cnf.Access, cnf.Secret)
	if err != nil {
		return nil, err
	}
	return &SMS{
		client:      client,
		cnTmpl:      cnf.CNTmpl,
		overseaTmpl: cnf.OverseaTmpl,
		signName:    cnf.SignName,
	}, nil
}

func (s *SMS) SendCode(areaCode, phone, code string) error {
	if areaCode == "+86" {
		return s.sendCNSMS(phone, code)
	}
	areaCode = strings.ReplaceAll(areaCode, "+", "")
	return s.sendOverseaSMS(areaCode+phone, code)
}

func (s *SMS) sendCNSMS(phoneNum, code string) error {
	return s.sendSMS(phoneNum, code, s.cnTmpl)
}

func (s *SMS) sendOverseaSMS(phoneNum, code string) error {
	return s.sendSMS(phoneNum, code, s.overseaTmpl)
}

func (s *SMS) sendSMS(phoneNum, code string, tmpl string) error {
	var request = dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = phoneNum
	request.SignName = s.signName
	request.TemplateCode = tmpl
	request.TemplateParam = fmt.Sprintf(`{"code":"%s"}`, code)

	var rsp, err = s.client.SendSms(request)
	if err != nil {
		ulog.Info("send.sms.code.error", zap.Error(err))
		return err
	}
	if rsp.Code != "OK" {
		ulog.Info("send.sms.code.not.ok",
			zap.String("msg", rsp.Message), zap.String("code", rsp.Code))
		return errors.New(rsp.Message)
	}
	return nil
}
