package service

import (
	"github.com/Danche23/Evenstar-Writings/internal/dto"
	"github.com/Danche23/Evenstar-Writings/internal/model"
	"github.com/Danche23/Evenstar-Writings/internal/repository"
	apperrors "github.com/Danche23/Evenstar-Writings/pkg/errors"
	"github.com/Danche23/Evenstar-Writings/pkg/utils"
)

// MessageService 留言业务逻辑
type MessageService struct {
	messageRepo *repository.MessageRepository
	userRepo    *repository.UserRepository
}

// NewMessageService 创建留言服务
func NewMessageService(messageRepo *repository.MessageRepository, userRepo *repository.UserRepository) *MessageService {
	return &MessageService{messageRepo: messageRepo, userRepo: userRepo}
}

// List 分页留言（带留言者简略信息）
func (s *MessageService) List(page, size int) (*dto.MessageListResponse, error) {
	page, size = utils.ClampPage(page, size, 10, 50)
	rows, total, err := s.messageRepo.List(page, size)
	if err != nil {
		return nil, apperrors.ErrInternalError
	}
	list := make([]dto.Message, 0, len(rows))
	for i := range rows {
		list = append(list, s.toMessage(&rows[i]))
	}
	return &dto.MessageListResponse{List: list, Total: total, Page: page, PageSize: size}, nil
}

// Create 留言（需登录）
func (s *MessageService) Create(userID uint, content string) (*dto.Message, error) {
	m := &model.Message{UserID: &userID, Content: content}
	if err := s.messageRepo.Create(m); err != nil {
		return nil, apperrors.ErrInternalError
	}
	r := s.toMessage(m)
	return &r, nil
}

// Delete 删除留言（自己或管理员）
func (s *MessageService) Delete(userID, messageID uint, isAdmin bool) error {
	m, err := s.messageRepo.FindByID(messageID)
	if err != nil {
		return apperrors.ErrResourceNotFound
	}
	if !isAdmin && (m.UserID == nil || *m.UserID != userID) {
		return apperrors.ErrForbidden
	}
	return s.messageRepo.Delete(messageID)
}

// toMessage 组装留言 DTO（用户注销显示 null）
func (s *MessageService) toMessage(m *model.Message) dto.Message {
	msg := dto.Message{
		ID:        m.ID,
		UserID:    m.UserID,
		Content:   m.Content,
		CreatedAt: m.CreatedAt,
	}
	if m.UserID != nil {
		if user, err := s.userRepo.FindByID(*m.UserID); err == nil {
			msg.User = &dto.CommentUser{ID: user.ID, Nickname: user.Nickname, Avatar: user.Avatar}
		}
	}
	return msg
}
