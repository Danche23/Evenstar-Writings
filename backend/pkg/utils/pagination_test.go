package utils

import "testing"

func TestClampPage(t *testing.T) {
	cases := []struct {
		name               string
		page, size         int
		defSize, maxSize   int
		wantPage, wantSize int
	}{
		{"默认值", 1, 10, 10, 50, 1, 10},
		{"page 为 0 归 1", 0, 10, 10, 50, 1, 10},
		{"page 为负归 1", -5, 10, 10, 50, 1, 10},
		{"size 为 0 取默认", 2, 0, 10, 50, 2, 10},
		{"size 为负取默认", 2, -1, 10, 50, 2, 10},
		{"size 超上限取上限", 1, 999999, 10, 50, 1, 50},
		{"后台上限 100", 1, 500, 10, 100, 1, 100},
		{"正常值不变", 3, 20, 10, 50, 3, 20},
	}
	for _, c := range cases {
		gotPage, gotSize := ClampPage(c.page, c.size, c.defSize, c.maxSize)
		if gotPage != c.wantPage || gotSize != c.wantSize {
			t.Errorf("%s: ClampPage(%d,%d,%d,%d) = (%d,%d)，期望 (%d,%d)",
				c.name, c.page, c.size, c.defSize, c.maxSize, gotPage, gotSize, c.wantPage, c.wantSize)
		}
	}
}

func TestClampPageDefensive(t *testing.T) {
	// 默认值为 0、max 小于 def 时也应产出正数，避免 Limit(-1) 全表扫描
	p, s := ClampPage(1, 0, 0, 0)
	if p != 1 || s < 1 {
		t.Errorf("异常入参未兜底: page=%d size=%d", p, s)
	}
}
