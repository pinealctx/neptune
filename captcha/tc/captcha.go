package tc

import (
	"errors"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/captcha/v20190722"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tErr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"go.uber.org/zap"

	"github.com/pinealctx/neptune/ulog"
)

var ErrCaptchaVerify = errors.New("captcha.verify.error")

type Captcha struct {
	cli       *v20190722.Client
	appID     uint64
	appSecret string
}

func New(cnf *Config) *Captcha {
	var credential = common.NewCredential(cnf.SecretID, cnf.SecretKey)
	var cpf = profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "captcha.tencentcloudapi.com"
	var client, _ = v20190722.NewClient(credential, "", cpf)
	return &Captcha{cli: client, appID: cnf.AppID, appSecret: cnf.AppSecret}
}

type VerifyReq struct {
	Ticket  string
	Randstr string
	UserIP  string
	Mini    bool
}

func (c *Captcha) Verify(req *VerifyReq) error {
	if req.Mini {
		return c.verifyMini(req.Ticket, req.UserIP)
	}
	return c.verify(req.Ticket, req.Randstr, req.UserIP)
}

func (c *Captcha) verify(ticket, randstr, userIP string) error {
	var request = v20190722.NewDescribeCaptchaResultRequest()
	request.CaptchaType = common.Uint64Ptr(9)
	request.Ticket = common.StringPtr(ticket)
	request.UserIp = common.StringPtr(userIP)
	request.Randstr = common.StringPtr(randstr)
	request.CaptchaAppId = common.Uint64Ptr(c.appID)
	request.AppSecretKey = common.StringPtr(c.appSecret)
	var response, err = c.cli.DescribeCaptchaResult(request)
	if err != nil {
		var sErr, ok = err.(*tErr.TencentCloudSDKError)
		if ok {
			ulog.Error("verify.SDK.err", zap.Reflect("sErr", sErr))
			return ErrCaptchaVerify
		}
		ulog.Error("verify.err", zap.Error(err))
		return err
	}
	if *response.Response.CaptchaCode != 1 {
		ulog.Error("verify.fail", zap.Reflect("response", response.Response))
		return ErrCaptchaVerify
	}
	return nil
}

func (c *Captcha) verifyMini(ticket, userIP string) error {
	var request = v20190722.NewDescribeCaptchaMiniResultRequest()
	request.CaptchaType = common.Uint64Ptr(9)
	request.Ticket = common.StringPtr(ticket)
	request.UserIp = common.StringPtr(userIP)
	request.CaptchaAppId = common.Uint64Ptr(c.appID)
	request.AppSecretKey = common.StringPtr(c.appSecret)
	var response, err = c.cli.DescribeCaptchaMiniResult(request)
	if err != nil {
		var sErr, ok = err.(*tErr.TencentCloudSDKError)
		if ok {
			ulog.Error("verifyMini.SDK.err", zap.Reflect("sErr", sErr))
			return ErrCaptchaVerify
		}
		ulog.Error("verifyMini.err", zap.Error(err))
		return err
	}
	if *response.Response.CaptchaCode != 1 {
		ulog.Error("verifyMini.fail", zap.Reflect("response", response.Response))
		return ErrCaptchaVerify
	}
	return nil
}
