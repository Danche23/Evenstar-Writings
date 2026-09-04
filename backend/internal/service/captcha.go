package service

import (
	"fmt"

	captcha20230305 "github.com/alibabacloud-go/captcha-20230305/client"
	openapiutil "github.com/alibabacloud-go/darabonba-openapi/v2/utils"
	"github.com/alibabacloud-go/tea/dara"

	"github.com/Danche23/Evenstar-Writings/pkg/config"
	apperrors "github.com/Danche23/Evenstar-Writings/pkg/errors"
	"github.com/Danche23/Evenstar-Writings/pkg/logger"
	"go.uber.org/zap"
)

// CaptchaService 阿里云验证码 2.0（滑块验证）服务端校验。
// 职责：接收前端提交的 captcha_verify_param，调用阿里云 VerifyIntelligentCaptcha 校验，
// 绝不信任前端的"验证成功"状态。
type CaptchaService struct {
	cfg    *config.CaptchaConfig
	client *captcha20230305.Client
}

// NewCaptchaService 创建验证码服务（配置不完整时先告警，调用时返回错误而非伪造成功）
func NewCaptchaService(cfg *config.CaptchaConfig) *CaptchaService {
	s := &CaptchaService{cfg: cfg}
	if err := s.initClient(); err != nil {
		logger.Warn("验证码服务初始化失败（该接口将返回错误，请补全 captcha 配置）", zap.Error(err))
	}
	return s
}

// initClient 初始化阿里云 Captcha 2.0 客户端
func (s *CaptchaService) initClient() error {
	c := s.cfg
	if c == nil || c.SceneID == "" || c.AccessKeyID == "" || c.AccessKeySecret == "" {
		return fmt.Errorf("captcha 配置缺失（scene_id / access_key_id / access_key_secret）")
	}
	endpoint := c.Endpoint
	if endpoint == "" {
		endpoint = "captcha.cn-shanghai.aliyuncs.com"
	}
	client, err := captcha20230305.NewClient(&openapiutil.Config{
		AccessKeyId:     dara.String(c.AccessKeyID),
		AccessKeySecret: dara.String(c.AccessKeySecret),
		Endpoint:        dara.String(endpoint),
	})
	if err != nil {
		return err
	}
	s.client = client
	return nil
}

// Verify 服务端校验滑块凭证。
// 返回约定（统一走 apperrors，供上层 response.BizError 直接映射）：
//   - 配置缺失/调用异常 → 内部错误（不泄露阿里云内部信息）
//   - 凭证为空 → 1007 captcha_required
//   - 校验不通过 → 1008 captcha_verify_failed
func (s *CaptchaService) Verify(captchaVerifyParam string) error {
	if s.client == nil {
		return apperrors.New(apperrors.CodeInternalError, "验证服务未配置，请联系管理员")
	}
	if captchaVerifyParam == "" {
		return apperrors.NewDefault(apperrors.CodeCaptchaRequired)
	}

	resp, err := s.client.VerifyIntelligentCaptchaWithOptions(
		&captcha20230305.VerifyIntelligentCaptchaRequest{
			SceneId:            dara.String(s.cfg.SceneID),
			CaptchaVerifyParam: dara.String(captchaVerifyParam),
		},
		&dara.RuntimeOptions{},
	)
	if err != nil {
		// 详细错误只进日志，不返回给前端
		logger.Error("阿里云验证码服务端校验调用失败", zap.Error(err))
		return apperrors.New(apperrors.CodeInternalError, "验证服务异常，请稍后重试")
	}
	if resp == nil || resp.Body == nil || resp.Body.Result == nil {
		logger.Warn("阿里云验证码返回结构异常")
		return apperrors.New(apperrors.CodeInternalError, "验证服务异常，请稍后重试")
	}
	if !dara.BoolValue(resp.Body.Result.VerifyResult) {
		return apperrors.NewDefault(apperrors.CodeCaptchaVerifyFailed)
	}
	return nil
}
