package service

import (
	"fmt"
	"net/smtp"

	"github.com/Danche23/Evenstar-Writings/pkg/config"
	"github.com/Danche23/Evenstar-Writings/pkg/logger"
)

// 邮件服务
type MailService struct {
	cfg *config.MailConfig
}

// 创建邮件服务
func NewMailService(cfg *config.MailConfig) *MailService {
	return &MailService{cfg: cfg}
}

// Send 发送邮件：mock 时打印日志，否则走 SMTP
func (s *MailService) Send(to, subject, body string) error {
	if s.cfg.Mock {
		logger.Infof("【邮件 mock】收件人=%s 主题=%s 内容=%s", to, subject, body)
		return nil
	}
	return s.sendSMTP(to, subject, body)
}

// sendSMTP 通过 SMTP 真实发送
func (s *MailService) sendSMTP(to, subject, body string) error {
	msg := fmt.Sprintf(
		"From: %s <%s>\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		s.cfg.FromName, s.cfg.FromAddr, to, subject, body,
	)

	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	auth := smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Host)
	if err := smtp.SendMail(addr, auth, s.cfg.FromAddr, []string{to}, []byte(msg)); err != nil {
		return fmt.Errorf("发送邮件失败: %w", err)
	}
	return nil
}
