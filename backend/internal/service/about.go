package service

import (
	"encoding/json"
	"strings"

	"github.com/Danche23/Evenstar-Writings/internal/dto"
	"github.com/Danche23/Evenstar-Writings/internal/repository"
	apperrors "github.com/Danche23/Evenstar-Writings/pkg/errors"
	"gorm.io/gorm"
)

// aboutKey site_settings 中「关于」页内容的 KV key
const aboutKey = "about"

// AboutService 前台「关于」页内容管理
type AboutService struct {
	aboutRepo *repository.AboutRepository
}

// NewAboutService 创建关于页服务
func NewAboutService(aboutRepo *repository.AboutRepository) *AboutService {
	return &AboutService{aboutRepo: aboutRepo}
}

// Get 返回关于页内容（未保存过时返回默认文案，不落库）
func (s *AboutService) Get() (*dto.AboutContent, error) {
	setting, err := s.aboutRepo.Get(aboutKey)
	if err == nil {
		var doc dto.AboutContent
		if jerr := json.Unmarshal([]byte(setting.Value), &doc); jerr == nil {
			return &doc, nil
		}
		// 历史数据解析失败：回退默认，不阻塞页面
	} else if err != gorm.ErrRecordNotFound {
		return nil, apperrors.ErrInternalError
	}
	return defaultAbout(), nil
}

// Save 保存关于页内容（幂等 upsert）
func (s *AboutService) Save(content *dto.AboutContent) error {
	// 规整：去首尾空白、去空标签项，避免脏数据
	stack := make([]string, 0, len(content.Stack))
	for _, t := range content.Stack {
		t = strings.TrimSpace(t)
		if t != "" {
			stack = append(stack, t)
		}
	}
	content.Stack = stack
	content.Intro = strings.TrimSpace(content.Intro)
	content.Footnote = strings.TrimSpace(content.Footnote)

	b, err := json.Marshal(content)
	if err != nil {
		return apperrors.ErrInternalError
	}
	if err := s.aboutRepo.Save(aboutKey, string(b)); err != nil {
		return apperrors.ErrInternalError
	}
	return nil
}

// defaultAbout 默认「关于」页内容（与前端首次展示一致）
func defaultAbout() *dto.AboutContent {
	return &dto.AboutContent{
		Intro:    "「暮星随笔」是一个以技术为主、内容不限的个人博客。整站由我独立实现：后端 Go + Gin + GORM，MySQL + Redis，JWT 鉴权，接入阿里云 OSS 与验证码；前端 Vue3 + Element Plus，Markdown 写作与渲染。\n\n不追热点，只写实践过、想明白的东西。偶尔也写点生活，放在角落，不影响主线。",
		Stack:    []string{"Go", "Gin", "GORM", "MySQL", "Redis", "JWT", "Vue3", "Element Plus", "Markdown"},
		Footnote: "站名 Evenstar（暮星）取自萨福笔下那颗把白昼散落之物带回家的星——愿每篇文章，都把你带回来。",
	}
}
