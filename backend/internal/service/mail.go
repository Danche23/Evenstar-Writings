package service

import (
	"crypto/tls"
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

// sendSMTP 通过 SMTP 真实发送。
// 端口 465 = SSL 隐式端口（先 TLS 握手再发命令）；80/25 = 明文端口（smtp.SendMail）。
func (s *MailService) sendSMTP(to, subject, body string) error {
	msg := fmt.Sprintf(
		"From: %s <%s>\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		s.cfg.FromName, s.cfg.FromAddr, to, subject, body,
	)

	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	auth := smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Host)

	// 465 是隐式 SSL：必须先建 TLS 连接，再用 SMTP 命令发信
	if s.cfg.Port == 465 {
		conn, err := tls.Dial("tcp", addr, &tls.Config{ServerName: s.cfg.Host})
		if err != nil {
			return fmt.Errorf("TLS 连接失败: %w", err)
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, s.cfg.Host)
		if err != nil {
			return fmt.Errorf("SMTP 握手失败: %w", err)
		}
		defer client.Close()

		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP 认证失败: %w", err)
		}
		if err := sendData(client, s.cfg.FromAddr, to, msg); err != nil {
			return err
		}
		return nil
	}

	if err := smtp.SendMail(addr, auth, s.cfg.FromAddr, []string{to}, []byte(msg)); err != nil {
		return fmt.Errorf("发送邮件失败: %w", err)
	}
	return nil
}

// sendData 用已连接的 SMTP client 依次发送 MAIL FROM / RCPT TO / DATA
func sendData(client *smtp.Client, from, to, msg string) error {
	if err := client.Mail(from); err != nil {
		return fmt.Errorf("MAIL FROM 失败: %w", err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("RCPT TO 失败: %w", err)
	}
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("DATA 失败: %w", err)
	}
	if _, err := w.Write([]byte(msg)); err != nil {
		w.Close()
		return fmt.Errorf("写入邮件内容失败: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("结束 DATA 失败: %w", err)
	}
	return nil
}
