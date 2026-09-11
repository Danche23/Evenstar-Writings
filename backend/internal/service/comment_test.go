package service

import (
	"github.com/Danche23/Evenstar-Writings/internal/model"
	"testing"
)

func TestReplyTargetsMustBelongToArticle(t *testing.T) {
	parent := &model.Comment{ArticleID: 1}
	if err := validateReplyArticle(parent, 2); err == nil {
		t.Fatal("cross-article parent should be rejected")
	}
	if err := validateReplyArticle(parent, 1); err != nil {
		t.Fatalf("same-article parent rejected: %v", err)
	}
}
